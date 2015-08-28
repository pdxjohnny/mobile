package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pdxjohnny/microsocket/client"
)

func NewClient() *client.Conn {
	ws := client.NewClient()
  ws.Recv = func (m []byte)  {}
  ws.ClientId, _ = os.Hostname()
	wsUrl := fmt.Sprintf("http://%s:%d/ws", "carpoolme.net", 8081)
	err := ws.Connect(wsUrl)
	if err != nil {
		log.Println(err)
	}
	go ws.Read()
	startMessage := fmt.Sprintf("Started %s", ws.ClientId)
	ws.Write([]byte(startMessage))
	return ws
}

func RedirectStdout() {
  r, w, _ := os.Pipe()
  os.Stdout = w
  os.Stderr = w

	go func(reader io.Reader) {
  	logger := NewClient()
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
      logger.Write([]byte(scanner.Text()))
		}
		if err := scanner.Err(); err != nil {
			message := fmt.Sprintf("There was an error with the logger %s", err)
			logger.Write([]byte(message))
		}
	}(r)
}
