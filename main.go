package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("Hello, world.\n")

	simulate_mongo_connection("localhost", 27017)

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 1 {
		process_packets(argsWithoutProg[0])
	} else {
		fmt.Println("please provide a *.pcap filename")
	}
}
