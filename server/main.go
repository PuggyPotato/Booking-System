package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5"
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
	log.Println("Password",dbpass)
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

	//Start http server
	log.Fatal(http.ListenAndServe(":8080",nil))
}

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

	//Saving To PostGresDB
	_,err = conn.Exec(
		context.Background(),
		"INSERT INTO public.users(username,password) VALUES ($1,$2)",
		user.Username,user.Password,
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

	//Check If Password Match
	if user.Password != storedPassword{
		http.Error(w,"Invalid username or password",http.StatusUnauthorized)
		return
	}

	//Success
	response := map[string]string{"message":"Login Succesful"}
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(response)
}

func withCORS(h http.HandlerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter,r * http.Request){
		//Allow all origin for cors
		w.Header().Set("Access-Control-Allow-Origin","*")
		w.Header().Set("Access-Control-Allow-Headers","Content-Type")
		w.Header().Set("Access-Control-Allow-Method","POST, GET, OPTIONS")

		if r.Method == "OPTIONS"{
		w.WriteHeader(http.StatusOK)
		return
		}
		h(w,r)
	}
}
