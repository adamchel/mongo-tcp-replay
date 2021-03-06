package main

import (
	"encoding/binary"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"strconv"
	"sync"
	"time"
)

// Earliest packet timestamp
var packetMinTimestamp time.Time

// Map of sending hosts to MongoConnections
var mapHostConnection map[string]*MongoConnection

type MongoPacket struct {
	delta   time.Duration
	payload []byte
}

func ProcessPackets(pcapFile string,
	mongodHost string,
	mongodPort string) {
	if handle, err := pcap.OpenOffline(pcapFile); err != nil {
		panic(err)
	} else {
		var connectionWaitGroup sync.WaitGroup
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		firstPacket := <-packetSource.Packets()
		packetMinTimestamp = GetPacketTime(firstPacket)
		mapHostConnection = make(map[string]*MongoConnection)
		SendPacket(firstPacket,
			&connectionWaitGroup,
			mongodHost,
			mongodPort)
		for packet := range packetSource.Packets() {
			SendPacket(packet,
				&connectionWaitGroup,
				mongodHost,
				mongodPort)
		}
		for _, mConnection := range mapHostConnection {
			mConnection.EOF()
		}
		connectionWaitGroup.Wait()
	}
}

func GetPacketTime(packet gopacket.Packet) time.Time {
	return packet.Metadata().CaptureInfo.Timestamp
}

func SendPacket(packet gopacket.Packet,
	connectionWaitGroup *sync.WaitGroup,
	mongodHost string,
	mongodPort string) {
	// If packet contains a mongo message
	if packet.ApplicationLayer() != nil {
		payload := packet.ApplicationLayer().Payload()
		delta := GetPacketTime(packet).Sub(packetMinTimestamp)

		// Get timestamp's delta from first packet
		// Get mongo wire protocol payload
		mongoPacket := MongoPacket{
			payload: payload,
			delta:   delta,
		}

		transportLayer := packet.TransportLayer()
		networkLayer := packet.NetworkLayer()

		var srcIp string
		var srcPort string

		if networkLayer.LayerType() == layers.LayerTypeIPv4 {
			ip4header := networkLayer.LayerContents()
			// Convert binary to IP string
			srcIp = strconv.Itoa(int(ip4header[12])) + "." +
				strconv.Itoa(int(ip4header[13])) + "." +
				strconv.Itoa(int(ip4header[14])) + "." +
				strconv.Itoa(int(ip4header[15]))
		}
		if transportLayer.LayerType() == layers.LayerTypeTCP {
			tcpHeader := transportLayer.LayerContents()
			// Hack to be able to use convert what should be a uint16 to string
			tcpHeaderSrcPort := []byte{0, 0, tcpHeader[0], tcpHeader[1]}
			srcPort = strconv.Itoa(int(binary.BigEndian.Uint32(tcpHeaderSrcPort)))
		}

		src := srcIp + ":" + srcPort

		if mConnection, ok := mapHostConnection[src]; ok {
			mConnection.Send(mongoPacket)
		} else {
			connectionWaitGroup.Add(1)
			mConnection := NewMongoConnection(mongodHost, mongodPort, 100)
			mapHostConnection[src] = mConnection
			go mConnection.ExecuteConnection(connectionWaitGroup)
			mConnection.Send(mongoPacket)
		}
	}
}
