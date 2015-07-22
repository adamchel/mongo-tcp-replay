package main

import (
	"fmt"
	"net"
	"sync"
	"time"
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

var simulationEpoch time.Duration = 0
var epochSet = false

type MongoConnection struct {
	mongodHost string
	mongodPort string
	packets chan MongoPacket
	done chan bool
}

func NewMongoConnection(mongodHost, mongodPort string, bufSize int) (*MongoConnection) {
	packets := make(chan MongoPacket, bufSize)
	return &MongoConnection{mongodHost, mongodPort, packets, make(chan bool)}
}
    
func (connection *MongoConnection) Send(packet MongoPacket) {
	connection.packets <- packet
}

func (connection *MongoConnection) EOF() {
	connection.done <- true
}
    
func (connection *MongoConnection) ExecuteConnection(waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	packetSend := make(chan MongoPacket)
	defer close(packetSend)
	go startMongoTCPConnection(connection.mongodHost, connection.mongodPort, packetSend)

    for {
		select {
		case packet := <- connection.packets:
			packetSend <- packet
		case done := <- connection.done:
			if(done || !done) {
				break
			}
		}    		
    }

	// Drain the packets
	for {
		select {
		case packet := <- connection.packets:
			packetSend <- packet
		default:
			// No more packets :)
			break
		}
	}
}

func startMongoTCPConnection(host, port string, packetChan chan MongoPacket) {
	var conn, error = net.Dial("tcp", host + ":" + port)
	if error != nil {
		fmt.Printf("Failed to connect to the mongod...\n")
		panic(error)
	}
	defer conn.Close()

	var readBuffer [4096]byte

	for {
		packet, isOpen :=<-packetChan
		if !isOpen {
			return
		}

		if !epochSet {
			simulationEpoch = time.Duration(time.Now().UnixNano())
			epochSet = true
		} else {
			time.Sleep((simulationEpoch + packet.delta) - time.Duration(time.Now().UnixNano()))
		}

		fmt.Printf("Sending packet with delta %d to mongod!", packet.delta)

		conn.Write(packet.payload)

		// Read the tcp reply into a buffer to discard
		conn.Read(readBuffer[0:])
	}
}

/*func make_connections(mConnection []MongoConnection) {
	for _, connection := range mConnection {
		connectionWaitGroup.Add(1)
		simulate_mongo_connection(connection)
	}
}*/

/*func simulate_mongo_connection(mConnection MongoConnection) {
	defer connectionWaitGroup.Done()
	// TODO: from command line args
	var conn, error = net.Dial("tcp", "localhost:27017")
	if error != nil {
		fmt.Printf("Failed to connect to the mongod...\n")
		return
	}
	var packetWaitGroup sync.WaitGroup
	for _, mPacket := range mConnection.packetList {
		packetWaitGroup.Add(1)
		go replay(conn, mPacket, packetWaitGroup)
	}
}*/

/*func replay(conn net.Conn,
	mPacket MongoPacket,
	wg sync.WaitGroup) {
	
	// Calculate our wait time as the TCP packet unix timestamp delta
	waitTime := mPacket.unixTimestamp
	timer := time.NewTimer(time.Duration(waitTime) * time.Millisecond)
	// Done with set up, but wait for all other replay go routines to be
	wg.Done()
	wg.Wait()
	// Delay write by the time delta
	<-timer.C
	conn.Write(mPacket.payload)
	// Read the tcp reply into a buffer to discard
	conn.Read(readBuffer[0:])
}*/
