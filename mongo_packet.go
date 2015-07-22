package main

import (
	"encoding/binary"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

// mongo wire protocol header field offset
var OFFSET int = 4

func process_packets(filename string) {
	if handle, err := pcap.OpenOffline(filename); err != nil {
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			handle_packet(packet)
		}
	}
}

// TODO: separate by src host/port
// TODO: get timestamps, delta from first packet
func handle_packet(packet gopacket.Packet) {
	// if packet contains a mongo message
	if packet.ApplicationLayer() != nil {
		payload := packet.ApplicationLayer().Payload()

		messageLength := binary.LittleEndian.Uint32(payload[0:OFFSET])
		requestID := binary.LittleEndian.Uint32(payload[OFFSET : 2*OFFSET])
		responseTo := binary.LittleEndian.Uint32(payload[2*OFFSET : 3*OFFSET])
		opCode := binary.LittleEndian.Uint32(payload[3*OFFSET : 4*OFFSET])

		fmt.Println(messageLength)
		fmt.Println(requestID)
		fmt.Println(responseTo)
		fmt.Println(opCode)
		fmt.Println()
	}
}
