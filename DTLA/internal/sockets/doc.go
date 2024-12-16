package sockets

/*
The JSON string received from the client should look like
```
{
	"Sender": {
		"Name": ("A" | "B"),
	},

	"Receiver": {
		"Name": ("A" | "B")
	},

	"Payload": {
		"Encrypted": (true | false),
		"Msg": "..."
	}
}
```
*/
