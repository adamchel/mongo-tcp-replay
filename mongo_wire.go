package main

import (
	"fmt"
	"net"
	"strconv"
	"encoding/binary"
)

/*type MsgHeader struct {
	messageLength int32 // total message size, including this
	requestID     int32 // identifier for this message
	responseTo    int32 // requestID from the original request (used in reponses from db)
	opCode        int32 // request type - see table in Wire protocol docs
}

type OpMsg struct {
    header 		MsgHeader	// standard message header
    message 	string 		// message for the database
}*/

// TODO: Eventually take a list of packets with time deltas and an epoch, and play them in the
//		 correct order and time. Currently sends the mongod on host:port an OP_MSG and gets 
//		 the response.
func simulate_mongo_connection(host string, port int64, packetList mongoPacket[]) {
	var conn, error = net.Dial("tcp",  host + ":" + strconv.FormatInt(port, 10))
	if error != nil {
		fmt.Printf("Failed to connect to the mongod..\n")
		return
	}
	for _, mPacket := range packetList {
		conn.Write(mPacket.payload)
	}

	var buf = [4096]byte
	// Read the tcp reply into a buffer to discard
	conn.Read(buf[0:])
}