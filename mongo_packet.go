package main

import (
	"encoding/binary"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"strconv"
	"strings"
)

// earliest packet timestamp
var MIN_UNIX_TIMESTAMP int64 = 0

// map of host to packets
var HOST_PACKET_MAP map[string][]MongoPacket

type MongoPacket struct {
	unixTimestamp int64
	payload       []byte
}

func process_packets(filename string) {
	if handle, err := pcap.OpenOffline(filename); err != nil {
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		firstPacket := <-packetSource.Packets()
		MIN_UNIX_TIMESTAMP = get_unix_timestamp(firstPacket)
		HOST_PACKET_MAP = make(map[string][]MongoPacket)
		for packet := range packetSource.Packets() {
			handle_packet(packet)
		}
	}
}

func get_unix_timestamp(packet gopacket.Packet) int64 {
	return packet.Metadata().CaptureInfo.Timestamp.Unix()
}

func handle_packet(packet gopacket.Packet) MongoPacket {
	// if packet contains a mongo message
	if packet.ApplicationLayer() != nil {
		payload := packet.ApplicationLayer().Payload()
		unixTimestamp := get_unix_timestamp(packet) - MIN_UNIX_TIMESTAMP

		// get timestamp's delta from first packet
		// get mongo wire protocol payload
		mongoPacket := MongoPacket{
			payload:       payload,
			unixTimestamp: unixTimestamp,
		}

		transportLayer := packet.TransportLayer()
		networkLayer := packet.NetworkLayer()

		var srcIp string
		var srcPort string

		// TODO: other protocols?
		if networkLayer.LayerType() == layers.LayerTypeIPv4 {
			ip4header := networkLayer.LayerContents()
			// convert binary to IP string
			srcIp = strconv.Itoa(int(ip4header[12])) + "." +
				strconv.Itoa(int(ip4header[13])) + "." +
				strconv.Itoa(int(ip4header[14])) + "." +
				strconv.Itoa(int(ip4header[15]))
		}
		if transportLayer.LayerType() == layers.LayerTypeTCP {
			tcpHeader := transportLayer.LayerContents()
			// disgusting hack to be able to use convert what should be a uint16 to string
			tcpHeaderSrcPort := []byte{0, 0, tcpHeader[0], tcpHeader[1]}
			srcPort = strconv.Itoa(int(binary.BigEndian.Uint32(tcpHeaderSrcPort)))
		}

		src := srcIp + ":" + srcPort
		HOST_PACKET_MAP[src] = append(HOST_PACKET_MAP[src], mongoPacket)

		return mongoPacket
	}
	return MongoPacket{unixTimestamp: packet.Metadata().CaptureInfo.Timestamp.Unix()}
}

func make_connection() {
	mongoConnections := []MongoConnection{}
	if len(HOST_PACKET_MAP) != 0 {
		for k, v := range HOST_PACKET_MAP {
			src := strings.Split(k, ":")
			mongoConnection := MongoConnection{
				host:       src[0],
				port:       src[1],
				packetList: v,
			}
			mongoConnections = append(mongoConnections, mongoConnection)
		}
	}
	make_connections(mongoConnections)
}
