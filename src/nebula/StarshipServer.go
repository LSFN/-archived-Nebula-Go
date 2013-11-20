package nebula

import (
	"net"

	"code.google.com/p/uuid"
	"code.google.com/p/goprotobuf/proto"

	"lsfn/common"
)

type StarshipServer struct {
	unjoinedClients []*StarshipListener
	orphanShipIDs []string
	joinedClients map[string]*StarshipListener
	networkChannels map[string]chan
	listenConnection net.TCPListener
	gameID string
	allowJoin bool
}

func (s *StarshipServer) handleConnectingStarship(conn) {
	conn.SetKeepAlive(true)
	starship := &StarshipListener{conn}
	append(s.unjoinedClients, starship)
	
	starship.SendMessage(joinInfoMessage)
	go starship.Listen()
}

func (s *StarshipServer) Listen() {
	s.listenConnection, err := net.Listen("tcp", ":39461")
	if err != nil {
		return false
	}
	gameID = uuid.New()
	allowJoin = true
	for {
		conn, err := ln.Accept()
		if err != nil {
			s.shutDown()
			break
		}
		handleConnectingStarship(conn)
	}
}

func (s *StarshipServer) shutDown() {
	s.listenConnection.Close()
	for client := range s.unjoinedClients {
		client.Disconnect()
	}

}

func (s *StarshipServer) processIncomingMessages() {
	
}