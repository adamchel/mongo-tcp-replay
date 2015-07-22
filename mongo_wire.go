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
func simulate_mongo_connection(host string, port int64) {
	var conn, error = net.Dial("tcp",  host + ":" + strconv.FormatInt(port, 10))
	if error != nil {
		fmt.Printf("Failed to connect to the mongod..\n")
		return
	}
	const BUFFER_LENGTH = 1024
	var buf [BUFFER_LENGTH]byte

	// Temporary hacky test.
	var message = "Are you alive?\x00"

	var messageLength uint32 = 16 + uint32(len(message))
	var requestID 	uint32 	= 90135
	var responseTo  uint32  = 0
	var opCode 		uint32 	= 1000

	var i = 0
	binary.LittleEndian.PutUint32(buf[i:i+4], messageLength)
	i += 4
	binary.LittleEndian.PutUint32(buf[i:i+4], requestID)
	i += 4
	binary.LittleEndian.PutUint32(buf[i:i+4], responseTo)
	i += 4
	binary.LittleEndian.PutUint32(buf[i:i+4], opCode)
	i += 4

	copy(buf[i:], message[:])
	i += len(message)

	fmt.Printf("Message to mongod: %s\n", buf[:i])

	conn.Write(buf[:i])

	var n_read, read_error = conn.Read(buf[0:])

	if (read_error != nil) {
		fmt.Printf("Failed to read from the mongod.\n")
		return
	}

	fmt.Printf("Response from mongod: response: %s\n", buf[:n_read])
	// End temporary hacky test.

	// TODO: send a collection of payloads in proper time order

}