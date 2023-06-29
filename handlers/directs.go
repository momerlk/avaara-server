/*

		THIS FILE IS FOR ALL FUNCTIONALITY RELATING TO DIRECT MESSAGING
		that is between two users 1-1 communication

		handlers :
		GET 			/directs  			client gets all of its direct messages
		GET				/directs/sorted		clients gets all directs sorted by time sent
		WEBSOCKET		/direct
*/

package handlers

import (
	"net/http"
	
	
	"bytes"
	"time"
	
	"encoding/json"
	
	"github.com/momerlk/avaara-server/database"
	"github.com/momerlk/avaara-server/structs"
	
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	
	"go.mongodb.org/mongo-driver/bson"
)

const directColl = "directs"


type GetDirectsBody struct {
	Sent 			[]structs.DirectMessage		`json:"sent"`
	Received 		[]structs.DirectMessage		`json:"received"`
}
func (a *Application) HandlerGetDirects(w http.ResponseWriter , r *http.Request){
GET(w , r , func(w http.ResponseWriter , r *http.Request){
	a.Debugln("GetDirects called")
	sess , ok := a.Verify(w , r) // verify the user with using cookies
	if !ok {
		a.Debugln("GetDirects : not verified")
		a.ClientError(w , http.StatusUnauthorized)
		return
	}
	
	
	sent , err := database.Get[structs.DirectMessage](a.DB , directColl , bson.D{{"sender" , sess.userId}})
	if err != nil {
		a.Debugln("GetDirects : could not get sent , err=" , err)
		a.ServerError(w , err)
		return
	}
	
	// get all direct messages with sender with given users id and receiver with clients users id
	received , err := database.Get[structs.DirectMessage](a.DB , directColl , bson.D{{"receiver" , sess.userId}})
	if err != nil {
		a.Debugln("GetDirects : could not get received , err=" , err)
		a.ServerError(w , err)
		return
	}
	
	resp := GetDirectsBody{
		Sent : sent,
		Received: received,
	}
	
	
	w.Header().Set("Content-Type" , "application/json")
	json.NewEncoder(w).Encode(resp)
})
}

// SortDirects : sorts the given directs
func (a  *Application) SortDirects(msgs []structs.DirectMessage){
	// sorting the array using the quick sort algorithm
	structs.QuickSort[structs.DirectMessage](msgs , 0 , len(msgs)-1 , func (a structs.DirectMessage , b structs.DirectMessage) int {
		return structs.CompareTime(structs.GetGoTime(a.TimeSent) , structs.GetGoTime(b.TimeSent)) * -1
	})
}

func (a *Application) HandlerSortedDirects(w http.ResponseWriter , r *http.Request){
GET(w , r , func(w http.ResponseWriter , r *http.Request){
	a.Debugln("SortedDirects called")
	sess , ok := a.Verify(w , r) // verify the user with using cookies
	if !ok {
		a.Debugln("SortedDirects : not verified")
		a.ClientError(w , http.StatusUnauthorized)
		return
	}


	sent , err := database.Get[structs.DirectMessage](a.DB , directColl , bson.D{{"sender" , sess.userId}})
	if err != nil {
		a.Debugln("GetDirects : could not get sent , err=" , err)
		a.ServerError(w , err)
		return
	}

	// get all direct messages with sender with given users id and receiver with clients users id
	received , err := database.Get[structs.DirectMessage](a.DB , directColl , bson.D{{"receiver" , sess.userId}})
	if err != nil {
		a.Debugln("GetDirects : could not get received , err=" , err)
		a.ServerError(w , err)
		return
	}
	
	// join received and sent into a single array
	msgs := make([]structs.DirectMessage , 0 , len(sent) + len(received))
	msgs = append(msgs , sent... )
	msgs = append(msgs , received...)	
	
	a.SortDirects(msgs)
	
	
	w.Header().Set("Content-Type" , "application/json")
	json.NewEncoder(w).Encode(msgs)
})
}


// DirectsReceive : format for receiving data from websocket
type DirectsReceive struct {
	ReceiverUsername 		string 		`json:"receiver"` 		//username of the receiver
	Content 				string 		`json:"content"`		//content of the message as text
}
// DirectsSend : format for sending data from websocket
type DirectsSend struct {
	Directs 				RenderedDirects					`json:"directs"` 		// all directs to be received
	Status 					int 							`json:"status"`			// status of operation
}

// Directs : handler for direct messaging using websockets
func (a *Application) Directs(w http.ResponseWriter , r *http.Request){
	a.Debugln("Directs called")
	
	// verifying the user
	sess , ok := a.Verify(w , r)
	if !ok {
		a.ClientError(w , http.StatusUnauthorized)
		return
	}
	
	// upgraades the connection to the websocket protocol
	conn , _ , _ , err := ws.UpgradeHTTP(r , w)
	if err != nil {
		a.ServerError(w , err)
		return
	}
	
	
	quit := make(chan bool) // a channel that will signal to the goroutines whether to stop or not
	a.Debugf("Directs : connected to %s!" , sess.userId)
	
	// a wrapper around websocket connection close
	closeConn := func(){
		a.Debugln("closeConn called")
		if <-quit {
			a.Debugln("quit is closed already!")
		} else {
		}
		close(quit) // closes quit channel and all go routines spawned with this connection
		err := conn.Close()
		if err != nil {
			a.Debugln("Directs : errored while trying to close the goroutine")
			a.ServerError(w , err)
			return
		}
		a.Debugf("Directs : connection with %s ended!" , sess.userId)
	}
	
	
	
	// go routine for receiving data 
	go func(){
		defer close(quit)
		for {
			select {
			case <-quit:
			default:
				// work
				raw , _ , err := wsutil.ReadClientData(conn) // reads data into bytes
				if err != nil {
					a.ServerError(w , err)
					return
				}
				decoder := json.NewDecoder(bytes.NewReader(raw)) // for reading bytes and then decoding into json
				
				
				
				var data DirectsReceive		// data received from websocket
				err = decoder.Decode(&data) // reades bytes then create a decoder to decode into json
				if err != nil {
					a.Debugln("Directs : couldn't decode received data into json")
					wsutil.WriteServerText(w , []byte("400"))
					continue
				}
				
				receiver , err := database.GetOne[structs.User](a.DB , usersColl , bson.D{{"username" , data.ReceiverUsername}})
				if err != nil {
					a.Debugln("failed to retrieve user with username=" , data.ReceiverUsername)
					a.ServerError(w , err)
					return
				}
				
				// get details of the sender from database
				sender , err := database.GetOne[structs.User](a.DB , usersColl , bson.D{{"id" , sess.userId}})
				if err != nil {
					a.Debugln("failed to retrieve user with username=" , data.ReceiverUsername)
					a.ServerError(w , err)
					return
				}
				
				// create a message object
				msg := structs.DirectMessage{
					Id : structs.GenerateID(),
					Sender : sess.userId,
					Receiver : receiver.Id,
					TimeSent : structs.CurrentTime(),
					Received: false,
					Content : data.Content,
					Attachment: "",
					Meta : structs.DirectMessageMeta{
						ReceiverName: receiver.Name,
						ReceiverUsername: receiver.Username,
						SenderName: sender.Name,
						SenderUsername: sender.Username,
					},
					Type : "direct",
				}
				
				// stores the data in the database
				err = a.DB.Store(directColl , msg)
				if err != nil {
					a.ServerError(w , err)
					return
				}
				
			}
		}
	}()
	
	// go routine for sending data
	go func(){
		defer close(quit) 
		for {
			select {
			case <-quit:
				
			default:
				// work
				var toSend DirectsSend
				unReceived , err := database.Get[structs.DirectMessage](a.DB , directColl ,
					bson.D{{"receiver" , sess.userId} , {"received" , false}},
				)
				if err != nil {
					a.ServerError(w , err)
					return
				}
				
				// if there are no unreceived messages
				if len(unReceived) == 0 || unReceived == nil {
					time.Sleep(1 * time.Second) // wait 1 second to limit the number of queries to database
					continue
				}
				
				// sorts the unreceived directs according to time
				a.SortDirects(unReceived)

				// render the directs
				rendered := a.RenderDirects(sess.userId , unReceived)

				// initializing data to send
				toSend.Directs = rendered
				toSend.Status = http.StatusOK
				
				if err != nil{
					a.Debugln("Directs : failed to encode data")
					a.ServerError(w ,err)
					return
				}
				
				b , err := json.Marshal(toSend)
				if err != nil {
					a.Debugln("Directs : failed to marshal json to send")
					a.ServerError(w , err)
					return
				}
				
				err = wsutil.WriteServerText(conn , b)
				if err != nil {
					a.Debugln("Directs : failed to write to websocket")
					a.ServerError(w , err)
					return
				}
				
				// updating database to match websocket client
				for _ , msg := range unReceived {
					msg.Received = true
					err = a.DB.Update(directColl , msg , bson.D{{"id" , msg.Id}})
					if err != nil {
						a.ServerError(w , err)
						return
					}
				}
				
				
			}
		}
	}()
	
	// continuing main loop so that send and receive goroutines don't close
	for {
		select {
		case <-quit:
			closeConn()
			return
		default:
		}
	}
}

