package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/gragas/woobloo-game/player"
	"github.com/gragas/woobloo-game/protocol"
	"github.com/gragas/woobloo-game/tile"
	"github.com/satori/go.uuid"
	"golang.org/x/net/websocket"
	"net/http"
)

var ip, port string
var mapWidth, mapHeight uint
var sockets map[uuid.UUID]*player.Player
var players []*player.Player
var tiles [][]*tile.Tile

func main() {
	// parse flags and args
	flag.StringVar(&ip, "ip", "localhost", "specifies the ip")
	flag.StringVar(&port, "port", "8080", "specifies the port")
	flag.UintVar(&mapWidth, "mapWidth", uint(50), "specifies the map width")
	flag.UintVar(&mapHeight, "mapHeight", uint(50), "specifies the map height")
	flag.Parse()
	nargs := flag.NArg()
	if nargs == 0 {
		panic(errors.New("No UUIDs found after parsing arguments."))
	}

	// initialize variables
	sockets = make(map[uuid.UUID]*player.Player)
	players = make([]*player.Player, 2)
	tiles = tile.NewDefaultMap(mapWidth, mapHeight)

	// setup handlers and listen
	for i := 0; i < nargs; i++ {
		http.Handle("/"+flag.Arg(i), websocket.Handler(accepter))
	}
	err := http.ListenAndServe(ip+":"+port, nil)
	if err != nil {
		panic(err)
	}
}

func accepter(ws *websocket.Conn) {
	buf := make([]byte, 512)
	_, err := ws.Read(buf)
	if err != nil {
		fmt.Printf("Error reading from websocket at first connection!")
		return
	}
	UUID, err := uuid.FromString(ws.Request().URL.Path[1:37])
	if err != nil {
		fmt.Printf("Invalid UUID!")
		return
	}
	p := player.NewDefaultPlayer()
	sockets[UUID] = p // TOOD: change this
	players = append(players, p)
	sendEverything(ws)
	if err != nil {
		fmt.Printf("Could not send everything at first connection!")
		return
	}
	go talk(ws, UUID)
}

func talk(ws *websocket.Conn, UUID uuid.UUID) {
	buf := make([]byte, 512)
	for {
		nbytes, err := ws.Read(buf)
		if err != nil {
			fmt.Printf("Error reading from websocket. Disconnecting.")
			return // stop talking to this websocket
		}
		if nbytes == 0 {
			fmt.Printf("Bad message from websocket. Continuing.")
			continue // bad message
		}
		err = handle(ws, UUID, buf[:nbytes])
		if err != nil {
			fmt.Printf("Error writing to websocket. Disconnecting.")
			return // stop talking to this websocket
		}
	}
}

func handle(ws *websocket.Conn, UUID uuid.UUID, msg []byte) error {
	var err error
	code := protocol.Code(msg[0])
	switch code {
	case protocol.GetEverything:
		sendEverything(ws)
	default:
		_, err = ws.Write(msg) // echo it back
	}
	return err
}

func sendEverything(ws *websocket.Conn) error {
	data := struct {
		Players []*player.Player
		Map     [][]*tile.Tile
	}{Players: players, Map: tiles}
	bytes, err := json.Marshal(&data)
	if err != nil {
		fmt.Printf("Error marshalling everything!")
		panic(err)
	}
	_, err = ws.Write(bytes)
	return err
}
