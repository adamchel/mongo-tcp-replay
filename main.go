package main

import (
	"flag"
)

func main() {
	pcapFilePtr := flag.String("pcapFile", "workload.pcap", "pcap workload file")
	mongodHostPtr := flag.String("mongodHost", "localhost", "host address for mongod")
	mongodPortPtr := flag.String("mongodPort", "27017", "port for mongod")
	interactiveMode := flag.Bool("interactiveMode", false, "run the pcap file in a single-threaded interactive mode.")

	flag.Parse()

	ProcessPackets(*pcapFilePtr, *mongodHostPtr, *mongodPortPtr, *interactiveMode)

}
