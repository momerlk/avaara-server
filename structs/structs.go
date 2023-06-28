package structs



type Time struct {
	Year 				uint		`json:"year"`
	Month 				uint		`json:"month"`
	Day 				uint		`json:"day"`
	Hour 				uint		`json:"hour"`
	Minute				uint		`json:"minute"`
	Second 				uint		`json:"second"`
}

type File struct {
	Id 					string		`json:"id"` 	// fileId or id of the file
	Data 				[]byte 		`json:"data"`	// Data in bytes
}

type Date struct {
	Year 				uint 			`json:"year"`
	Month				uint 			`json:"month"`
	Day 				uint 			`json:"day"`
}


type DirectMessageMeta struct {
	SenderUsername			string			`json:"sender_username"`
	SenderName				string			`json:"sender_name"`
	ReceiverUsername		string			`json:"receiver_username"`
	ReceiverName			string			`json:"receiver_name"`
}
type DirectMessage struct {
	Id					string 				`json:"id"`					
	
	Sender 				string 				`json:"sender"`				// userId of the sender
	Receiver			string 				`json:"receiver"`			// chatId of the receiver 		*
	
	TimeSent 			Time				`json:"time_sent"`			// the time message was sent	
	Received 			bool 				`json:"received"`			// * 
	
	Content 			string 				`json:"content"`			// the content of the message *
	Attachment 			string 				`json:"attachment"` 		// fileId of the attachment
	
	Meta 				DirectMessageMeta	`json:"meta"`				// meta information	
	
	Type 				string 				`json:"type"`				// type of the message either "direct" or "group"
}

type User struct {
	Id 					string 			`json:"id"`				// the userId of the user prefixed with 'user-'
	
	Name 				string 			`json:"name"`			// the name of the user *
	Age 				int				`json:"age"`			// the age of the user
	DOB 				string 			`json:"dob"`			// the date of birth of the user *
	Gender 				string 			`json:"gender"`			// the gender of the user *
	Pfp					string 			`json:"pfp"`			// the fileid of profile picture of the user

	Username 			string 			`json:"username"`		// the username of the user *
	Email 				string 			`json:"email"`			// the email of the user *
	Password 			string 			`json:"password"` 		// the password stored as a hash *
}



// Chat  conversation between two or more users
type Chat struct {
	Id 					string 				`json:"id"`			// chatId
	Members 			[]string			`json:"members"`	//userId of the members
}