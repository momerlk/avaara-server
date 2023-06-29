package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"net/http"
	"runtime/debug"

	"github.com/alexedwards/scs"

	"github.com/momerlk/avaara-server/database"
	"github.com/momerlk/avaara-server/structs"
)

var (
	key = string(structs.GenerateFileId())
	store = scs.NewCookieManager(key)
)


type Application struct {
	InfoLog 		*log.Logger
	ErrorLog 		*log.Logger
	DebugLog 		*log.Logger
	Debug 			bool
	DB				*database.Database
}
func (a *Application) ServerError(w http.ResponseWriter , err error){
	trace := fmt.Sprintf("%s\n%s" , err.Error() , debug.Stack())
	a.ErrorLog.Println(trace)
	http.Error(w , http.StatusText(http.StatusInternalServerError) , http.StatusInternalServerError)
}
func (a *Application) ClientError(w http.ResponseWriter , status int){
	http.Error(w , http.StatusText(status) , status)
}
func (a *Application) NotFound(w http.ResponseWriter){
	a.ClientError(w , http.StatusNotFound)
}

func POST(w http.ResponseWriter , r *http.Request , f func (http.ResponseWriter, *http.Request)){
	if r.Method != "POST"{
		w.Header().Set("Allow" , "POST")
		http.Error(w , "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	f(w,r)
}

func GET(w http.ResponseWriter , r *http.Request , f func (http.ResponseWriter, *http.Request)){
	if r.Method != "GET"{
		w.Header().Set("Allow" , "GET")
		http.Error(w , "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	
	f(w,r)
	
}

func (a *Application) Debugf(format string , v ...any){
	if a.Debug == true {
		a.DebugLog.Printf(format , v)
	}
}
func (a *Application) Debugln(v ...any){
	if a.Debug == true {
		a.DebugLog.Println(v)
	}
}

func WriteMsg(w http.ResponseWriter , r *http.Request , msg string){
	w.Header().Set("Content-Type" , "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message" : msg})
}




