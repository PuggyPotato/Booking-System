package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var users = make(map[string]string)
var conn *pgx.Conn

type UserAuth struct{
	Username string `json:"username"`
	Password string `json:"password"`
}

func main(){

	_ = godotenv.Load()
	var dbpass = os.Getenv("PASS")
	var err error
	var connectionStr string = fmt.Sprintf("postgres://postgres:%v@localhost:5432/booking",dbpass)
	conn,err = pgx.Connect(context.Background(),connectionStr)

	if err !=nil{
		log.Fatal("Unable To Connect To The Database:",err)
	}
	defer conn.Close(context.Background())
	fmt.Println("Connected To PostgresSQL.")
	fmt.Println("Server is running on localhost:8080")

	//Routes
	http.HandleFunc("/register",withCORS(handleRegister))
	http.HandleFunc("/login",withCORS(handleLogin))
	http.HandleFunc("/appointment",withCORS(bookAppointment))
	http.HandleFunc("/adminAppointments",withCORS(admin))
	http.HandleFunc("/adminLogin",withCORS(adminLogin))

	//Start http server
	log.Fatal(http.ListenAndServe(":8080",nil))
}

//Handle Register
func handleRegister(w http.ResponseWriter,r *http.Request){
	log.Println("HandleRegister called.")

	if r.Method != http.MethodPost{
		http.Error(w,"Bad JSON",http.StatusMethodNotAllowed)
		return
	}

	var user UserAuth
	if err := json.NewDecoder(r.Body).Decode(&user);err != nil{
		http.Error(w,"Bad JSON",http.StatusBadRequest)
		return
	}

	log.Println("Received user:",user)

	var exists bool 
	query := "SELECT EXISTS (SELECT 1 FROM public.users WHERE username = $1)"
	err := conn.QueryRow(context.Background(),query,user.Username).Scan(&exists)

	if err != nil{
		http.Error(w,"User already Exist",http.StatusConflict)
		return
	}

	//Encrypt password
	hashedPassword,err := bcrypt.GenerateFromPassword([]byte(user.Password),bcrypt.DefaultCost )
	if err !=nil{
		log.Println("Password hash Error:",err)
		http.Error(w,"Failed To Hash Password:",http.StatusInternalServerError)
		return
	}
	//Saving To PostGresDB
	_,err = conn.Exec(
		context.Background(),
		"INSERT INTO public.users(username,password) VALUES ($1,$2)",
		user.Username,string(hashedPassword),
	)
	if err !=nil{
		log.Println("Error saving to DB:",err)
		http.Error(w,"Failed To Save User",http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message":"User Received!",}
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(response)

}


//Handle Login
func handleLogin(w http.ResponseWriter,r *http.Request){
	log.Println("handleLogin called.")


	if r.Method != http.MethodPost{
		http.Error(w,"BAD JSON",http.StatusMethodNotAllowed)
		return
	}

	var user UserAuth
	if err :=json.NewDecoder(r.Body).Decode(&user);err !=nil{
		http.Error(w,"Bad JSON",http.StatusBadRequest)
	}

	log.Println("Login Attempt:",user)

	var storedPassword string
	query := "SELECT password FROM public.users WHERE username = $1"
	err:= conn.QueryRow(context.Background(),query,user.Username).Scan(&storedPassword)

	log.Printf("Raw Input: Name='%v',Password='%v'",user.Username,user.Password)

	if err !=nil{
		log.Println("Login DB Error:",err)
		http.Error(w,"Invalid Username Or Password",http.StatusUnauthorized)
	}

	//Turn password To Hashed and check if password match
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword),[]byte(user.Password))
	if err !=nil{
		http.Error(w,"Invalid Username Or Password",http.StatusUnauthorized)
		return
	}

	token,err := generateJWT(user.Username)
	if err !=nil{
		log.Println("Token generation error:",err)
		http.Error(w,"Internal server error",http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":"Login Succesful",
		"token":token,
	})
}

func withCORS(h http.HandlerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter,r * http.Request){
		//Allow all origin for cors
		w.Header().Set("Access-Control-Allow-Origin","*")
		w.Header().Set("Access-Control-Allow-Headers","Content-Type,Authorization")
		w.Header().Set("Access-Control-Allow-Method","POST, GET, OPTIONS")

		if r.Method == "OPTIONS"{
		w.WriteHeader(http.StatusOK)
		return
		}
		h(w,r)
	}
}

func generateJWT(username string) (string,error){
	var secretKey = os.Getenv("SECRET")
	var jwtSecret = []byte(secretKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"username":username,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString(jwtSecret)
}

func bookAppointment(w http.ResponseWriter,r * http.Request){
	if r.Method != http.MethodPost{
		http.Error(w,"Only Post Allowed",http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == ""{
		http.Error(w,"Missing Authorization Header",http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader,"Bearer ")
	username,err := validateJWT(tokenString)
	if err!=nil{
		http.Error(w,"Invalid Token",http.StatusUnauthorized)
		return
	}

	//Parse Request Body
	var payload struct{
		Date string `json:"date"`
		Time string `json:"time"`
		Reason string `json:"reason"`
	}

	if err:= json.NewDecoder(r.Body).Decode(&payload); err!=nil{
		http.Error(w,"Invalid JSON",http.StatusBadRequest)
		return
	}

	//Validate Input
	if payload.Date == "" || payload.Time == "" || payload.Reason == ""{
		http.Error(w,"Missing Fields",http.StatusBadRequest)
		return
	}

	//Save Appointment
	_,err = conn.Exec(
		context.Background(),
		"INSERT INTO appointments (username,date,time,reason) VALUES ($1,$2,$3,$4)",
		username, payload.Date,payload.Time,payload.Reason,
	)

	if err !=nil{
		log.Println("DB Insert Error:",err)
		http.Error(w,"Database Error",http.StatusInternalServerError)
		return
	}

	//Respond With Success
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":"Appointment booked.",
	})
}

func validateJWT(tokenString string) (string,error){
	secret := os.Getenv("SECRET")
	token,err := jwt.Parse(tokenString,func (token * jwt.Token)(interface{},error){
		return []byte(secret),nil
	})

	if claims,ok := token.Claims.(jwt.MapClaims); ok && token.Valid{
		username := claims["username"].(string)
		return username,nil
	}

	return "",err
}

func admin(w http.ResponseWriter,r *http.Request){
	if r.Method != http.MethodGet{
		http.Error(w,"Only Get Allowed",http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == ""{
		http.Error(w,"Missing Authorization Header",http.StatusUnauthorized)
		return
	}

	tokenString:=strings.TrimPrefix(authHeader, "Bearer ")
	username,err := validateJWT(tokenString)
	if err !=nil{
		http.Error(w,"Invalid Token",http.StatusUnauthorized)
		return
	}

	//Restrict To Admin Only
	if username != "potato"{
		http.Error(w,"Unauthorized",http.StatusForbidden)
		return
	}

	//Query all appointment
	rows,err := conn.Query(context.Background(),"SELECT username,date,time,reason FROM appointments")
	if err !=nil{
		log.Println("DB Query Error:",err)
		http.Error(w,"Database Error",http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var appointments []map[string]string

	for rows.Next(){
		var user, timeStr,reason string
		var date time.Time

		if err :=rows.Scan(&user,&date,&timeStr,&reason);err !=nil{
			log.Println("Row scan Error",err)
			continue
		}
		appointments = append(appointments,map[string]string{
			"username":user,
			"date":date.Format("2006-01-02"),
			"time":timeStr,
			"reason":reason,
		})
	}
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(appointments)
}

func adminLogin(w http.ResponseWriter,r * http.Request){
	log.Println("adminLogin called.")


	if r.Method != http.MethodPost{
		http.Error(w,"BAD JSON",http.StatusMethodNotAllowed)
		return
	}

	var creds UserAuth
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Bad JSON", http.StatusBadRequest)
		return
	}

	//Getting Admin Password
	adminUsername := os.Getenv("ADMIN")
	adminPassword := os.Getenv("ADMINPASS")

	//Check Credentials
	if creds.Username != adminUsername || creds.Password !=adminPassword{
		log.Printf("Invalid admin login attempt: username=%s",creds.Username)
		http.Error(w,"Unauthorised",http.StatusUnauthorized)
		return
	}

	token,err := generateJWT(creds.Username)
	if err !=nil{
		log.Println("Token generation error:",err)
		http.Error(w,"Internal server error",http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":"Login Succesful",
		"token":token,
	})
}