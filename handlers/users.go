package handlers

import (
	"encoding/json"
	"net/http"
	
	"fmt"
	
	"crypto/sha512"
	
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"oserver/database"
	"oserver/structs"
)

const usersColl = "users"

func Hash(x string) string{
	hasher := sha512.New()
	res := hasher.Sum([]byte(x))
	return fmt.Sprintf("%x" , res)
}

// Register 'user/register' endpoint
func (a *Application) Register(w http.ResponseWriter , r *http.Request){
POST(w , r , func (w http.ResponseWriter , r *http.Request){
	var user structs.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		a.ServerError(w , err)
		return
	}
	
	a.Debugln("Registed user =" , user)
	
	user.Id = structs.GenerateUserId()
	user.Age , err = structs.GetAge(user.DOB)
	if err != nil {
		a.ServerError(w , err)
		return
	}
	user.Password = Hash(user.Password)
	
	err = a.DB.Store(usersColl , user)
	if err != nil {
		a.ServerError(w , err)
		return
	}
	
	
	w.Header().Set("Content-Type" , "application/json")
	
	json.NewEncoder(w).Encode(bson.D{{"msg" , "successfully registered user!"}})
})
}



func (a *Application) Login(w http.ResponseWriter , r *http.Request){
POST(w , r , func(w http.ResponseWriter , r *http.Request){
	var body Credentials
	
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		a.Debugln("Error occured while decoding json body , err =" , err)
		a.ServerError(w , err)
		return 
	}
	
	var user structs.User
	
	if body.Email == "" {
		user , err = database.GetOne[structs.User](a.DB , usersColl , bson.D{{"username" , body.Username}})
	} else {
		user , err = database.GetOne[structs.User](a.DB , usersColl , bson.D{{"email" , body.Email}})
	}
	if err != nil {
		a.ServerError(w , err)
		return 
	}
	
	w.Header().Set("Content-Type" , "application/json")
	
	if Hash(body.Password) == user.Password { 
		// authenticated
		
		token := GenerateSessionToken()
		expiresAt := GenerateSessionExpiry(12 * time.Hour)
		Sessions[token] = Session{
			userId : user.Id,
			expiry : expiresAt,
		}
		
		a.Debugln("Login : session =" , Sessions[token])
		
		http.SetCookie(w , &http.Cookie{
			Name : "session_token",
			Value : token,
			Expires : expiresAt,
		})
		
		a.Debugln("Set cookie with token :" , token)
		
		json.NewEncoder(w).Encode(bson.D{{"msg" , "work hard!"}})
		
	} else {
		a.ClientError(w , http.StatusUnauthorized)
	}
})
}

func (a *Application) Logout(w http.ResponseWriter , r *http.Request){
	session := store.Load(r)
	err := session.PutString(w , "authenticated" , "false")
	if err != nil {
		a.ServerError(w , err)
		return
	}
}

type UserIDResponse struct {
	UserId 			string 		`json:"userid"`
}
func (a *Application) UserID(w http.ResponseWriter , r *http.Request){
GET(w , r , func(w http.ResponseWriter , r *http.Request){
	sess , ok := a.Verify(w , r)
	if !ok {
		a.ClientError(w , http.StatusUnauthorized)
		return
	}
	
	a.Debugln("UserId : sess =" , sess , "ok =" , ok)

	username := r.URL.Query().Get("username")
	res , err := database.GetOne[structs.User](a.DB , usersColl , bson.D{{"username" , username}})
	if err != nil {
		a.ServerError(w , err)
		return
	}
	w.Header().Set("Content-Type" , "application/json")
	json.NewEncoder(w).Encode(UserIDResponse{
		UserId : res.Id,
	})
})
}

type UsernameResponse struct {
	Username 			string 			`json:"username"`
	Name 				string 			`json:"name"`
}
func (a *Application) Username(w http.ResponseWriter , r *http.Request){
	GET(w , r , func(w http.ResponseWriter , r *http.Request){
		_ , ok := a.Verify(w , r)
		if !ok {
			a.ClientError(w , http.StatusUnauthorized)
			return
		}
		
		a.Debugln("/username called")

		userid := r.URL.Query().Get("userid")
		a.Debugln("/username userid =" , userid)
		res , err := database.GetOne[structs.User](a.DB , usersColl , bson.D{{"id" , userid}})
		if err != nil {
			a.ServerError(w , err)
			return
		}
		a.Debugln("/username final username =" , res.Username)
		w.Header().Set("Content-Type" , "application/json")
		json.NewEncoder(w).Encode(UsernameResponse{
			Username : res.Username,
			Name : 	res.Name,
			})
	})
}


func (a *Application) Details(w http.ResponseWriter , r *http.Request){
GET(w , r , func(w http.ResponseWriter , r *http.Request) {
	sess , ok := a.Verify(w , r)
	if ok == false {
		a.ClientError(w , http.StatusUnauthorized)
		return
	}
	
	res , err := database.GetOne[structs.User](a.DB , usersColl , bson.D{{"id" , sess.userId}})
	if err != nil {
		a.ServerError(w , err)
		return
	}
	
	res.Password = ""
	
	w.Header().Set("Content-Type" , "application/json")
	json.NewEncoder(w).Encode(res)
})
}

func (a *Application) Verify(w http.ResponseWriter , r *http.Request) (Session , bool) {
	
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return Session{} , false
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return Session{} , false
	}
	sessionToken := c.Value
	
	a.Debugln("Verify : token =" , sessionToken)

	// We then get the session from our session map
	userSession, exists := Sessions[sessionToken]
	if !exists {
		// If the session token is not present in session map, return an unauthorized error
		w.WriteHeader(http.StatusUnauthorized)
		return Session{} , false
	}
	// If the session is present, but has expired, we can delete the session, and return
	// an unauthorized status
	if userSession.IsExpired() {
		delete(Sessions, sessionToken)
		w.WriteHeader(http.StatusUnauthorized)
		return Session{} , false
	}
	
	return userSession , true
}
