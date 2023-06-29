package main

import (
	"log"
	"os"

	"net/http"

	"github.com/momerlk/avaara-server/database"
	"github.com/momerlk/avaara-server/handlers"
	
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

	mux.HandleFunc("/directs" , app.HandlerGetDirects) // get all directs
	mux.HandleFunc("/directs/sorted" , app.HandlerSortedDirects) // get all directs sorted by time

	mux.HandleFunc("/rendered/directs" , app.HandlerRenderDirects) // gets rendered form of the directs

	mux.HandleFunc("/direct" , app.Directs) // websocket endtime for real time messaging
	

	infoLog.Printf("Starting server on %s\n" , PORT)
	
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://192.168.8.100:5173",
			"http://localhost:5173",
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

