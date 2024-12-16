package sockets

import (
	"dtla/internal/util"
	"errors"
	"fmt"
	"io"
	"math"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"golang.org/x/net/websocket"
)

type TCPPayload struct {
	Encrypted bool
	Msg       string
}

type TCPWire struct {
	Packet     gopacket.Packet
	PacketDump string
}

type TCPHost struct {
	Name string

	Addr string
	Port int
}

// type wsCon websocket.Conn

type sockWSCon struct {
	*websocket.Conn
}

var dbg = fmt.Println

// Listen for JSON requests on `host`
// `portMin` and `portMax` values define the range for randomly generated port
// numbers for the sockets with names defined in `socketNames`
func (ws *sockWSCon) listenAndServeJSON(host string, portMin, portMax int, socketNames []string) error {
	var err error
	var wg sync.WaitGroup

	ws.MaxPayloadBytes = 32 << 20 // 32MB

	table := Table{
		rows:    make([]TableRow, 64, 512),
		sockets: make(map[string]*TCPHost, len(socketNames)),
	}
	row := TableRow{
		table: &table,
	}

	for {
		err = ws.SetDeadline(time.Now().Add(time.Hour))
		if err != nil {
			util.LogError(err.Error())
			break
		}

		err = websocket.JSON.Receive(ws.Conn, &row)
		if err != nil {
			if err != io.EOF {
				util.LogError(err.Error())
			}
			break
		}

		sockOpts := socketOptions{
			host:        host,
			portMin:     portMin,
			portMax:     portMax,
			socketNames: socketNames,
			senderName:  row.Sender.Name,
			recvName:    row.Receiver.Name,
		}
		err = table.retryInitSockets(sockOpts, 3, 3)
		if err != nil {
			util.LogError(err.Error())
			break
		}
		defer table.con.Close()
		defer table.ln.Close()
		row.setAddrPort()

		fmt.Printf("SenderPort: %d, RecvPort: %d\n", row.Sender.Port, row.Receiver.Port)

		packetCh := make(chan gopacket.Packet, 64)
		packetErrCh := make(chan error, 1)
		wg.Add(1)
		go listenPackets(packetCh, packetErrCh, &wg)
		wg.Wait()

		err = row.execute()
		if err != nil {
			break
		}

		var packet gopacket.Packet
		packet, err = waitForPacket(packetCh, packetErrCh, row.Sender.Port, row.Receiver.Port)
		if err != nil {
			break
		}
		packetErrCh <- errStopListening

		row.Wire = TCPWire{
			Packet:     packet,
			PacketDump: packet.Dump(),
		}

		table.rows = append(table.rows, row)

		err = websocket.JSON.Send(ws.Conn, &row)
		if err != nil {
			util.LogError(err.Error())
			break
		}
	}

	return err
}

var errStopListening error = errors.New("listenPackets should stop listening")

func listenPackets(packetCh chan gopacket.Packet, errCh chan error, wg *sync.WaitGroup) {
	handle, err := pcap.OpenLive("lo", math.MaxInt32, false, pcap.BlockForever)
	if err != nil {
		wg.Done()
		packetCh <- nil
		errCh <- err
		return
	}
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	wg.Done()

	// TODO: Use BPF to filter packets?
	// If there's too many irrelevant packets that aren't discarded from
	// `packetCh` buffer, it will be filled, block and won't be able to
	// respond to `errCh`.

	for packet := range packetSource.Packets() {
		select {
		case err = <-errCh:
			if err == errStopListening {
				return
			}
		default:
		}
		packetCh <- packet
	}
}

func waitForPacket(packetCh chan gopacket.Packet, errCh chan error, senderPort, recvPort int) (gopacket.Packet, error) {
	var err error
	var packet gopacket.Packet

	envV, err := util.EnvGetInt("V")
	if err != nil {
		return nil, err
	}

WAIT_FOR_PACKETS:
	for {
		// These are recoverable errors, fatal errors are output at the end
		if err != nil && envV > 0 {
			util.LogInfo(err.Error())
		}

		packet = <-packetCh

		select {
		case err = <-errCh:
			break WAIT_FOR_PACKETS
		default:
		}

		TCPLayer := packet.Layer(layers.LayerTypeTCP)
		if TCPLayer == nil {
			err = errors.New("Not a TCP packet")
			continue
		}

		srcPort := TCPLayer.(*layers.TCP).SrcPort
		dstPort := TCPLayer.(*layers.TCP).DstPort

		switch srcPort {
		case layers.TCPPort(senderPort):
		case layers.TCPPort(recvPort):
		default:
			err = errors.New("srcPort didn't match sender or receiver port")
			continue
		}

		switch dstPort {
		case layers.TCPPort(senderPort):
		case layers.TCPPort(recvPort):
		default:
			err = errors.New("dstPort didn't match sender or receiver port")
			continue
		}

		if len(TCPLayer.LayerPayload()) == 0 {
			err = errors.New("Packet empty")
			continue
		}

		util.LogInfo("Packet matched")

		err = nil
		break
	}
	// In case fatal errors are possible in the future
	if err != nil {
		util.LogError(err.Error())
	}

	return packet, err
}
