package sockets

import (
	"dtla/internal/util"
	"errors"
	"io"
	"math/rand"
	"net"
)

type TableRow struct {
	Error  bool
	ErrMsg string

	Sender   TCPHost
	Wire     TCPWire
	Receiver TCPHost
	Payload  TCPPayload

	table *Table
}

type Table struct {
	sockets map[string]*TCPHost
	ln      *net.TCPListener
	con     *net.TCPConn
	rows    []TableRow
}

type socketOptions struct {
	host                 string
	portMin, portMax     int
	socketNames          []string
	senderName, recvName string
}

// Try to initialize `Table.sockets` map `retryMap` number of times which
// might fail because the randomly generated port number for a socket might
// collide with another socket in the map.
// Try to initialize `Table.con` and `Table.ln` `retryCon` number of times
// which might fail because the port is unavailable.
func (t *Table) retryInitSockets(opt socketOptions, retryMap, retryCon int) error {
	var err error

	for retryMap > 0 && retryCon > 0 {
		retryMap--
		err = t.initSocketMap(opt)
		if err != nil {
			continue
		}

		retryCon--
		err = t.initLnCon(opt)
		if err != nil {
			continue
		}
		break
	}

	return err
}

func (t *Table) initSocketMap(opt socketOptions) error {
	var err error

	envV, err := util.EnvGetInt("V")
	if err != nil {
		return err
	}

	for _, v := range opt.socketNames {
		t.sockets[v] = &TCPHost{
			Name: v,
			Addr: opt.host,
			Port: rand.Intn(opt.portMax-opt.portMin) + opt.portMin,
		}
	}

LOOP2:
	for i, v := range opt.socketNames {
		for ii, vv := range opt.socketNames {
			if i == ii {
				continue
			}

			if t.sockets[v].Port == t.sockets[vv].Port {
				err = errors.New("At least two sockets have the same port")
				break LOOP2
			}
		}
	}

	if err != nil && envV > 0 {
		util.LogError(err.Error())
	}

	return err
}

func (t *Table) initLnCon(opt socketOptions) error {
	var err error

	envV, err := util.EnvGetInt("V")
	if err != nil {
		return err
	}

	senderAddr := net.TCPAddr{
		IP:   net.ParseIP(t.sockets[opt.senderName].Addr),
		Port: t.sockets[opt.senderName].Port,
	}

	recvAddr := net.TCPAddr{
		IP:   net.ParseIP(t.sockets[opt.recvName].Addr),
		Port: t.sockets[opt.recvName].Port,
	}

	t.ln, err = net.ListenTCP("tcp4", &recvAddr)
	if err != nil {
		if envV > 0 {
			util.LogError(err.Error())
		}
		return err
	}

	t.con, err = net.DialTCP("tcp4", &senderAddr, &recvAddr)
	if err != nil && envV > 0 {
		util.LogError(err.Error())
	}

	return err
}

func (row *TableRow) setAddrPort() {
	row.Sender.Addr = row.table.sockets[row.Sender.Name].Addr
	row.Sender.Port = row.table.sockets[row.Sender.Name].Port
	row.Receiver.Addr = row.table.sockets[row.Receiver.Name].Addr
	row.Receiver.Port = row.table.sockets[row.Receiver.Name].Port
}

func (row *TableRow) execute() error {
	var err error

	recvErrCh := make(chan error)
	go row.recv(recvErrCh)

	_, err = row.table.con.Write([]byte(row.Payload.Msg))
	if err != nil {
		util.LogError(err.Error())
		return err
	}
	row.table.con.Close()
	util.LogInfof("Sent[%s]:\t\t%s\n", row.Sender.Name, row.Payload.Msg)

	return <-recvErrCh
}

func (row *TableRow) recv(ch chan error) {
	var err error

	con, err := row.table.ln.AcceptTCP()
	if err != nil {
		util.LogError(err.Error())
		ch <- err
		return
	}
	defer con.Close()

	buf, err := io.ReadAll(con)
	if err != nil && err != io.EOF {
		util.LogError(err.Error())
		ch <- err
		return
	}
	row.Payload.Msg = string(buf)
	util.LogInfof("Received[%s]: \t%s\n", row.Receiver.Name, row.Payload.Msg)

	ch <- nil
}
