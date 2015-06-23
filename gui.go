package gui

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
)

var (
	tcpConn    net.Conn          // Can be global as we don't intend to server multiple connections
	IP         = "10.0.1.8:8888" //"127.0.0.1:8888" // Feel free to serve across Network / LAN
	ScreenSize = 1024            // If we stick to a power of 2, integer division is easier
	html       string
)

// -- Expose Setup API
func Screen(s int) {
	ScreenSize = s
}
func Address(s string) {
	IP = s
}
func Launch() {

	// -- Listen & serve the GUIdisplay
	fmt.Println("\nWaiting for Display : Please navigate to " + IP + " to commence....")
	listener, _ := net.Listen("tcp", IP)
	for webSocReady := false; webSocReady != true; {
		tcpConn, _ = listener.Accept()
		webSocReady = handleTCP()
	}
	// --
}

// --

// -- Expose runtime API
func Wipe() {
	wsFrame := byte(1*128 + 1*2) // Simplified : FIN bit & Binary Type
	wsPayload := byte(1)         // 0 guiData bytes + 1 guiCmd byte
	guiCmd := byte(0 + 1<<3)     // guiCmd code 0 + single packet
	tcpConn.Write([]byte{wsFrame, wsPayload, guiCmd})
}
func Close() {
	wsFrame := byte(1*128 + 1*2) // Simplified : FIN bit & Binary Type
	wsPayload := byte(1)         // 0 guiData bytes + 1 guiCmd byte
	guiCmd := byte(7 + 1<<3)     // guiCmd code 0 + single packet
	tcpConn.Write([]byte{wsFrame, wsPayload, guiCmd})
	tcpConn.Close()
}
func Plot(x int, y int, r uint8, g uint8, b uint8, a uint8) {
	wsFrame := byte(1*128 + 1*2) // Simplified : FIN bit & Binary Type
	wsPayload := byte(9)         // 8 guiData bytes + 1 guiCmd byte
	guiCmd := byte(2 + 1<<3)     // guiCmd code 2 + single packet
	guiData := []byte{hiByte(x), lowByte(x), hiByte(y), lowByte(y), r, g, b, a}
	tcpConn.Write(append([]byte{wsFrame, wsPayload, guiCmd}, guiData...))
}
func FillRect(x int, y int, w int, h int, r uint8, g uint8, b uint8, a uint8) {
	wsFrame := byte(1*128 + 1*2) // Simplified : FIN bit & Binary Type
	wsPayload := byte(13)        // 12 guiData bytes + 1 guiCmd byte
	guiCmd := byte(4 + 1<<3)     // guiCmd code 4 + single packet
	guiData := []byte{hiByte(x), lowByte(x), hiByte(y), lowByte(y), hiByte(w), lowByte(w), hiByte(h), lowByte(h), r, g, b, a}
	tcpConn.Write(append([]byte{wsFrame, wsPayload, guiCmd}, guiData...))
}
func Circle(x int, y int, radius int, r uint8, g uint8, b uint8, a uint8) {
	wsFrame := byte(1*128 + 1*2) // Simplified : FIN bit & Binary Type
	wsPayload := byte(11)        // 10 guiData bytes + 1 guiCmd byte
	guiCmd := byte(5 + 1<<3)     // guiCmd code 5 + single packet
	guiData := []byte{hiByte(x), lowByte(x), hiByte(y), lowByte(y), hiByte(radius), lowByte(radius), r, g, b, a}
	tcpConn.Write(append([]byte{wsFrame, wsPayload, guiCmd}, guiData...))
}

// --

func hashWithMagicKey(clientKey string) string {
	hasher := sha1.New()
	hasher.Write([]byte(clientKey + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11")) // MagicKey
	return "Sec-WebSocket-Accept: " + base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

func readBytesOnWire() string {
	messageBuffer := make([]byte, 1024)
	tcpConn.Read(messageBuffer)
	return string(messageBuffer)
}

func handleTCP() bool {
	guiDisplay := strings.Replace(html, "GUI_SCREEN_SIZE", strconv.Itoa(ScreenSize), -1)
	guiDisplay = strings.Replace(guiDisplay, "GUI_IP", IP, -1)
	var Upgrade, clientKey string

	// -- Assume incoming HTTP GET request for WebSocket Upgrade on TCP connection. Parse Upgrade & Key if present
	request := readBytesOnWire()
	if UpgradeIndex := strings.Index(request, "Upgrade:"); UpgradeIndex != -1 {
		Upgrade = request[UpgradeIndex+9 : UpgradeIndex+9+9]
	}
	if clientKeyIndex := strings.Index(request, "Sec-WebSocket-Key:"); clientKeyIndex != -1 {
		clientKey = request[clientKeyIndex+19 : clientKeyIndex+19+24]
	}

	// -- Serve GUIdisplay if not WebSocket upgrade request
	if Upgrade == "" {
		tcpConn.Write([]byte(guiDisplay))
		fmt.Println("*** Serving GUIdisplay Page ***")
		tcpConn.Close()
		// Otherwise handshake and start sedning display...
	} else if Upgrade == "websocket" {
		acceptKey := hashWithMagicKey(clientKey)
		tcpConn.Write([]byte(wsUpgrade + acceptKey + "\r\n\r\n"))
		fmt.Println("*** GUIdisplay Opened WebSocket ***")
		return true
	}
	return false
	//---
}

//
// -- Useful functions Todo : Tidy up
//

func lowByte(i int) uint8 {
	return uint8(i & 0xff)
}

func hiByte(i int) uint8 {
	return uint8((i & 0xff00) >> 8)
}

func wsWrite(guiPacket [9]byte) {
	wsFrame := []byte{byte(1*128 + 1*2), 9} // FIN bit + Binary Type
	tcpConn.Write([]byte{wsFrame[0], wsFrame[1], guiPacket[0], guiPacket[1], guiPacket[2], guiPacket[3], guiPacket[4], guiPacket[5], guiPacket[6], guiPacket[7], guiPacket[8]})
}

// -- Setup
func init() {
	htmlB, err := ioutil.ReadFile("../GUIdisplay.html") // Relative from where running the main app from. Assumes sub dir of gui
	if err != nil {
		panic(err)
	}
	html = string(htmlB)
}

const wsUpgrade = "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nOrigin: null\r\nSec-WebSocket-Protocol: guiSocket-protocol\r\n"

// --
