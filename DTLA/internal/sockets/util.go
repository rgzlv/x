package sockets

import (
	"dtla/internal/util"

	"golang.org/x/net/websocket"
)

func errLogAndSend(ws *websocket.Conn, tr *TableRow, err error) {
	util.LogErrorDepth(3, err.Error())
	tr.Error = true
	tr.ErrMsg = err.Error()
	err = websocket.JSON.Send(ws, &tr)
	if err != nil {
		util.LogErrorDepth(3, err.Error())
	}
}
