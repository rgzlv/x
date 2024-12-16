package sockets

import (
	"dtla/internal/util"
	"io"
	"math"

	"golang.org/x/net/websocket"
)

func handler(ws *websocket.Conn) {
	var err error

	sockWS := &sockWSCon{ws}

	err = sockWS.listenAndServeJSON("127.0.0.1", 1024, math.MaxUint16, []string{"A", "B"})
	if err != nil && err != io.EOF {
		var row TableRow
		row.Error = true
		row.ErrMsg = err.Error()
		err = websocket.JSON.Send(ws, &row)
		if err != nil {
			util.LogError(err.Error())
			return
		}
	}
}
