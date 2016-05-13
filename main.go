package main

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/woobloo/woobloo-game/player"
	"github.com/woobloo/woobloo-game/protocol"
	"github.com/woobloo/woobloo-game/tile"
	"github.com/satori/go.uuid"
	"golang.org/x/net/websocket"
	"net/http"
)

const TransferSize = 4096

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
	buf := make([]byte, TransferSize)
	_, err := ws.Read(buf)
	if err != nil {
		fmt.Println("Error reading from websocket at first connection!")
		fmt.Println(err)
		return
	}
	UUID, err := uuid.FromString(ws.Request().URL.Path[1:37])
	if err != nil {
		fmt.Println("Invalid UUID!")
		return
	}
	p := player.NewDefaultPlayer()
	sockets[UUID] = p // TOOD: change this
	players = append(players, p)
	sendEverything(ws)
	if err != nil {
		fmt.Println("Could not send everything at first connection!")
		return
	}
	go talk(ws, UUID)
}

func talk(ws *websocket.Conn, UUID uuid.UUID) {
	buf := make([]byte, TransferSize)
	for {
		nbytes, err := ws.Read(buf)
		if err != nil {
			fmt.Println("Error reading from websocket. Disconnecting.")
			return // stop talking to this websocket
		}
		if nbytes == 0 {
			fmt.Println("Bad message from websocket. Continuing.")
			continue // bad message
		}
		err = handle(ws, UUID, buf[:nbytes])
		if err != nil {
			fmt.Println("Error writing to websocket. Disconnecting.")
			return // stop talking to this websocket
		}
	}
}

func handle(ws *websocket.Conn, UUID uuid.UUID, msg []byte) error {
	var err error
	code := protocol.Code(msg[0])
	switch code {
	case protocol.GetEverything:
		err = sendEverything(ws)
	case protocol.GetTile:
		err = sendTile(ws, msg)
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
		fmt.Println("Error marshalling everything!")
		panic(err)
	}
	gzws := gzip.NewWriter(ws)
	if _, err = gzws.Write(bytes); err != nil {
		return err
	}
	if err = gzws.Flush(); err != nil {
		return err
	}
	err = gzws.Close()
	return err
}

func sendTile(ws *websocket.Conn, msg []byte) error {
	pos := struct{X uint; Y uint}{}
	err := json.Unmarshal(msg[1:], pos)
	if err != nil {
		fmt.Println("Bad GetTile message.")
		return nil
	}
	bytes, err := json.Marshal(tiles[pos.Y][pos.X])
	if err != nil {
		panic(err)
	}
	_, err = ws.Write(bytes)
	return err
}
