package main

import (
	"log"
	"os"

	"net/http"

	"oserver/database"
	"oserver/handlers"
	
	"github.com/rs/cors"
)

const uri = "mongodb://localhost:27017"
const dbName = "avaara"


func main(){
	PORT := os.Getenv("PORT")
	infoLog := log.New(os.Stdout , "INFO\t" , log.Ldate | log.Ltime)
	errorLog := log.New(os.Stderr , "ERROR\t", log.Ldate | log.Ltime | log.Lshortfile)
	debugLog := log.New(os.Stdout , "DEBUG\t" , log.Ldate | log.Ltime | log.Llongfile)
	
	db := &database.Database{}
	err := db.Connect(uri , dbName)
	if err != nil {
		errorLog.Fatal(err)
		return
	}
	defer db.Close(errorLog)

	mux := http.NewServeMux()
	app := &handlers.Application{
		InfoLog : infoLog,
		ErrorLog : log.New(os.Stderr , "ERROR\t", log.Ldate | log.Ltime | log.Lshortfile),
		Debug : true,
		DebugLog: debugLog,
		DB : db,
	}
	
	mux.HandleFunc("/register" , app.Register) // register
	mux.HandleFunc("/login" , app.Login) // login 
	mux.HandleFunc("/details" , app.Details) // get user details
	mux.HandleFunc("/userid" , app.UserID) // get user id for a username
	mux.HandleFunc("/username" , app.Username) // get username for userid
	mux.HandleFunc("/direct" , app.GetDirects) // get all directs
	mux.HandleFunc("/directsorted" , app.SortedDirects) // get all directs sorted by time
	mux.HandleFunc("/directs" , app.Directs) // websocket endtime for real time messaging
	

	infoLog.Printf("Starting server on %s\n" , PORT)
	
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:8081",
			"http://localhost:19006",
		},
		AllowCredentials: true,
		Debug: false,
	})
	
	server := &http.Server{
		Addr : ":" + PORT,
		ErrorLog : errorLog,
		Handler : c.Handler(mux),
	}
	
	err = server.ListenAndServe()
	errorLog.Fatal(err)
}

