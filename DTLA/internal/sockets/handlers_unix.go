//go:build unix

package sockets

import "golang.org/x/net/websocket"

func Handler(ws *websocket.Conn) {
	handler(ws)
}
