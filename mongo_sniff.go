package main

import (
	"encoding/binary"
)

// mongo wire protocol header field offset
var OFFSET32 int = 4
var OFFSET64 int = 8

type MongoMessage struct {
	messageLength uint32
	requestID     uint32
	responseTo    uint32
	opCode        uint32
	message       []byte
}

/*type struct OP_UPDATE {
    ZERO uint32                 // 0 - reserved for future use
    fullCollectionName string  // "dbname.collectionname"
    flags uint32            	   // bit vector. see below
    documents []byte           // the query to select the document
              				   // specification of the update to perform
}

type struct OP_INSERT {
    flags uint32           		// bit vector - see below
    fullCollectionName string   // "dbname.collectionname"
    documents []byte     		// one or more documents to insert into the collection
}

struct OP_QUERY {
    flags uint32 		  // bit vector of query options.  See below for details.
    fullCollectionName string  // "dbname.collectionname"
    numberToSkip uint32    // number of documents to skip
    numberToReturn uint32  // number of documents to return
                                      //  in the first OP_REPLY batch
    documents []byte        // query object.  See below for details.
 							// Optional. Selector indicating the fields
                           //  to return.  See below for details.
}

struct OP_GET_MORE {
    ZERO uint32           // 0 - reserved for future use
    fullCollectionName string    // "dbname.collectionname"
    numberToReturn uint32 // number of documents to return
    cursorID uint64       // cursorID from the OP_REPLY
}

struct OP_DELETE {
    ZERO uint32       		  // 0 - reserved for future use
    fullCollectionName string // "dbname.collectionname"
    flags uint32         	  // bit vector - see below for details.
    selectorDoc []byte        // query object.  See below for details.
}

struct OP_KILL_CURSORS {
    ZERO uint32   			// 0 - reserved for future use
    numberOfCursorIDs uint32 // number of cursorIDs in message
    cursorIDs []uint64       // sequence of cursorIDs to close
}
*/
type OP_REPLY struct {
    responseFlags uint32    // bit vector - see details below
    cursorID uint64        // cursor id if client needs to do get more's
    startingFrom uint32    // where in the cursor this reply is starting
    numberReturned uint32   // number of documents in the reply
    documents []byte     // documents
}

// TODO: separate by src host/port
func handle_packet(packet MongoPacket) (MongoMessage, string) {
	// if packet contains a mongo message
	payload := packet.payload

	messageLength := binary.LittleEndian.Uint32(payload[0:OFFSET32])
	requestID := binary.LittleEndian.Uint32(payload[OFFSET32 : 2*OFFSET32])
	responseTo := binary.LittleEndian.Uint32(payload[2*OFFSET32 : 3*OFFSET32])
	opCode := binary.LittleEndian.Uint32(payload[3*OFFSET32 : 4*OFFSET32])
	message := payload[4*OFFSET32 : messageLength]

	mongoMessage := MongoMessage{
		messageLength: messageLength,
		requestID:     requestID,
		responseTo:    responseTo,
		opCode:        opCode,
		message:       message,
	}

	switch (opCode) {
		case 1001:
			return mongoMessage, "OP_MSG"
			break
		case 2001:
			return mongoMessage, "OP_UPDATE"
			break
		case 2002:
			return mongoMessage, "OP_INSERT"
			break
		case 2003:
			return mongoMessage, "RESERVED"
			break
		case 2004:
			return mongoMessage, "OP_QUERY"
			break
		case 2005:
			return mongoMessage, "OP_GET_MORE"
			break
		case 2006:
			return mongoMessage, "OP_DELETE"
			break
		case 2007:
			return mongoMessage, "OP_KILL_CURSORS"
			break
		default:
	
	}
	return mongoMessage, "NO!"
}

func parseOpReply (message []byte, messageLength int) OP_REPLY {

	responseFlags 	:= binary.LittleEndian.Uint32(message[0:OFFSET32])
	cursorID 		:= binary.LittleEndian.Uint64(message[OFFSET32 : 3*OFFSET32])
	startingFrom 	:= binary.LittleEndian.Uint32(message[3*OFFSET32 : 4*OFFSET32])
	numberReturned 	:= binary.LittleEndian.Uint32(message[4*OFFSET32 : 5*OFFSET32])
	documents 		:= message[5*OFFSET32 : (messageLength - 16)]

	opReply := OP_REPLY {
		responseFlags:   responseFlags,
		cursorID:        cursorID,
		startingFrom:    startingFrom,
		numberReturned:  numberReturned,
		documents:       documents,
	}

	return opReply
}

/*
func parseOpMsg (message []byte) OP_MSG {

}

func parseOpUpdate (message []byte) OP_UPDATE {
	
}

func parseOpInsert (message []byte) OP_INSERT {
	
}

func parseOpQuery (message []byte) OP_QUERY {
	
}

func parseOpGetMore (message []byte) OP_GET_MORE {
	
}

func parseOpDelete (message []byte) OP_DELETE {
	
}

func parseOpKillCursors (message []byte) OP_KILL_CURSORS {
	
}*/