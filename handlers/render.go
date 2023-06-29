/*
	
		THIS FILE CONTAINS HANDLERS FOR SERVER SIDE RENDERING AND PROCESSING OF DATA
		
		GET						/rendered/directs					returns rendered directs for the specific user
*/

package handlers

import (
	"net/http"

	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/momerlk/avaara-server/database"
	"github.com/momerlk/avaara-server/structs"
)

type Primitive interface {
	string | int64 | int32 | float64 | float32
}

type Set[T Primitive] struct {
	Values 			[]T
}
func (s *Set[T]) Insert(val T){
	for i := 0;i < len(s.Values);i++ {
		if s.Values[i] == val {
			return
		}
	}
	s.Values = append(s.Values , val)
}

type RenderedDirectSelf struct {
	Userid 				string 								`json:"userid"`
	Username 			string 								`json:"username"`
	Name 				string 								`json:"name"`
}
type RenderedDirect struct {
	Username 			string 								`json:"username"`
	Name 				string 								`json:"name"`
	Text 				string 								`json:"text"`
	Type 				string 								`json:"type"` // can be sent or received
}
type RenderedDirects struct {
	Self				RenderedDirectSelf					`json:"self"`
	Keys 				[]string							`json:"keys"`
	Directs 			map[string][]RenderedDirect			`json:"directs"`
}

// RenderDirects : internally render directs
func (a *Application) RenderDirects(userId string , msgs []structs.DirectMessage) RenderedDirects{
	var directs = map[string][]RenderedDirect{} // key : username of the other user , value : rendered directs
	keys := &Set[string]{} // exists due to javascript's lack of getting object keys
	var self *RenderedDirectSelf = nil // user details of the client

	// looping through each of the sorted msgs
	for _ , msg := range msgs {
		Type := "" // type of the msgs 'sent' or 'received'
		username := "" // username of the other user
		name := ""	 // name of the other user
		if userId == msg.Sender {
			Type = "sent"

			// as sender is self , other user is receiver
			username = msg.Meta.ReceiverUsername
			name = msg.Meta.ReceiverName

			// if self not initialized initialize self
			if self == nil {
				self = &RenderedDirectSelf{}
				self.Name = msg.Meta.SenderName
				self.Username = msg.Meta.SenderUsername
				self.Userid = msg.Sender
			}

		} else {
			Type = "received"

			// self is receiver so other user is sender
			username = msg.Meta.SenderUsername
			name = msg.Meta.SenderName

			// repeating this code block as there may be no messages that have been sent
			if self == nil {
				self = &RenderedDirectSelf{}
				self.Name = msg.Meta.ReceiverName
				self.Username = msg.Meta.ReceiverUsername
				self.Userid = msg.Receiver
			}

		}

		keys.Insert(username) // inserts username into key which doesnt accept duplicate values

		// append a direct into the direct slice of a specific user
		directs[username]  = append(directs[username] , RenderedDirect{
			Text : msg.Content,
			Username: username,
			Name : name,
			Type: Type,
			})
	}

	// response of the direct
	resp := RenderedDirects {
		Self : *self,
		Directs : directs,
		Keys : keys.Values,
	}

	return resp
}

func (a *Application) HandlerRenderDirects(w http.ResponseWriter , r *http.Request){
GET(w , r , func(w http.ResponseWriter , r *http.Request){
	a.Debugln("RenderDirects called")
	sess , ok := a.Verify(w , r)// verify the user with using cookies
	if !ok {
		a.Debugln("RenderDirects : not verified")
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

	// sorting the msgs using the quick sort algorithm
	structs.QuickSort[structs.DirectMessage](msgs , 0 , len(msgs)-1 , func (a structs.DirectMessage , b structs.DirectMessage) int {
		return structs.CompareTime(structs.GetGoTime(a.TimeSent) , structs.GetGoTime(b.TimeSent)) * -1
	})
	

	resp := a.RenderDirects(sess.userId , msgs)	

	a.Debugln("RenderDirects body =" , resp)

	w.Header().Set("Content-Type" , "application/json")
	json.NewEncoder(w).Encode(resp)
})
}