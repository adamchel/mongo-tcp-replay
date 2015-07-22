package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
	"sync"
)

/*

type MsgHeader struct {
	messageLength int32 // total message size, including this
	requestID     int32 // identifier for this message
	responseTo    int32 // requestID from the original request (used in reponses from db)
	opCode        int32 // request type - see table in Wire protocol docs
}

type OpMsg struct {
    header 		MsgHeader	// standard message header
    message 	string 		// message for the database
}

*/

type MongoConnection struct {
	host string
	port string
	packetList []MongoPacket
}

var connectionWaitGroup sync.WaitGroup

func make_connections(mConnection []MongoConnection) {
	for _, connection in mConnection {
		connectionWaitGroup.Add(1)
		simulate_mongo_connection(mConnection)
	}
}

func simulate_mongo_connection(mConnection MongoConnection) {
	defer connectionWaitGroup.Done()
	var conn, error = net.Dial("tcp", mConnection.host + ":" + mConnection.port)
	if error != nil {
		fmt.Printf("Failed to connect to the mongod...\n")
		return
	}
	var packetWaitGroup sync.WaitGroup
	for _, mPacket := range mConnection.packetList {
		packetWaitGroup.Add(1)
		go replay(conn, mPacket, packetWaitGroup)
	}
}

func replay(conn net.Conn,
		    mPacket MongoPacket,
		    wg sync.WaitGroup) {
	var readBuffer [4096]byte
	// Calculate our wait time as the TCP packet unix timestamp delta
	waitTime := mPacket.unixTimestamp
	timer := time.NewTimer(time.Duration(waitTime) * time.Millisecond)
	// Done with set up, but wait for all other replay go routines to be
	wg.Done()
	wg.Wait()
	// Delay write by the time delta
	<- timer.C
	conn.Write(mPacket.payload)
	// Read the tcp reply into a buffer to discard
	conn.Read(readBuffer[0:])
}