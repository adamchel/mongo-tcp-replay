package main

import (
	"encoding/binary"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

// mongo wire protocol header field offset
const OFFSET int = 4

// earliest packet timestamp
var MIN_UNIX_TIMESTAMP int64 = 0

type MongoPacket struct {
	unixTimestamp int64
	messageLength uint32
	requestID     uint32
	responseTo    uint32
	opCode        uint32
	message       []byte
}

func process_packets(filename string) {
	if handle, err := pcap.OpenOffline(filename); err != nil {
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		firstPacket := <-packetSource.Packets()
		MIN_UNIX_TIMESTAMP = get_unix_timestamp(firstPacket)
		for packet := range packetSource.Packets() {
			handle_packet(packet)
		}
	}
}

func get_unix_timestamp(packet gopacket.Packet) int64 {
unixTimestamp:
	packet.Metadata().CaptureInfo.Timestamp.Unix()
}

// TODO: separate by src host/port
func handle_packet(packet gopacket.Packet) MongoPacket {
	// if packet contains a mongo message
	if packet.ApplicationLayer() != nil {
		payload := packet.ApplicationLayer().Payload()

		// get timestamp's delta from first packet
		// get mongo wire protocol payload
		unixTimestamp := get_unix_timestamp(packet) - MIN_UNIX_TIMESTAMP
		messageLength := binary.LittleEndian.Uint32(payload[0:OFFSET])
		requestID := binary.LittleEndian.Uint32(payload[OFFSET : 2*OFFSET])
		responseTo := binary.LittleEndian.Uint32(payload[2*OFFSET : 3*OFFSET])
		opCode := binary.LittleEndian.Uint32(payload[3*OFFSET : 4*OFFSET])
		message := payload[4*OFFSET : messageLength]

		mongoPacket := MongoPacket{
			unixTimestamp: unixTimestamp,
			messageLength: messageLength,
			requestID:     requestID,
			responseTo:    responseTo,
			opCode:        opCode,
			message:       message,
		}

		return mongoPacket
	}
	return MongoPacket{unixTimestamp: packet.Metadata().CaptureInfo.Timestamp.Unix()}
}
