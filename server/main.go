package main

import (

	"net/http"
	"encoding/json"
	"log"
	"context"

	"github.com/jackc/pgx/v5"
)

var users = make(map[string]string)
var conn *pgx.Conn

type UserAuth struct{
	Username string `json:"name"`
	Password string `json:"password"`
}

func main(){
	
}

func handleRegister(){

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
	query := "SELECT password FROM public.users WHERE name = $1"
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
		w.Header().Set("Access-Control-Allow-Origin","Content-Type")
		w.Header().Set("Access-Control-Allow-Origin","POST,GET,OPTIONS")

		if r.Method == "OPTIONS"{
		w.WriteHeader(http.StatusOK)
		return
		}
		h(w,r)
	}
}
