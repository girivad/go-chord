package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	overlay "github.com/girivad/go-chord/Overlay"
)

func main() {
	// Read cmdline args (ip + contact address if joining an existing chord ring)
	contactPtr := flag.String("contact", "None", "IP address of a contact in a Chord Ring you want to join.")

	if len(os.Args) < 3 {
		fmt.Println("Please provide this node's IP address and the capacity of its chord ring.")
		os.Exit(1)
	}

	flag.Parse()
	ip := os.Args[1]
	capacity, err := strconv.ParseInt(os.Args[2], 2, 64)

	if err != nil {
		fmt.Println("Error Reading Capacity:", err)
		os.Exit(1)
	}

	// TO-DO: Implement IP Verification (verify that this is a valid IP address, at least via Regex)

	// Create a new ChordNode and join an existing chord ring if requested.
	chordServer := overlay.NewChordServer(ip, capacity)

	if *contactPtr != "None" {
		// Make a client for the contact, and then run a join service on it.
		contactNode, err := overlay.Connect(*contactPtr)

		if err != nil {
			fmt.Println("Unable to connect to the contact in the Chord Ring:", err)
			os.Exit(1)
		}

		err = chordServer.Join(contactNode)

		if err != nil {
			fmt.Println("Failed to join chord ring:", err)
			os.Exit(1)
		}
	}

	// Serve data from 8080, gRPC through 8081.
	err = chordServer.Serve()

	if err != nil {
		fmt.Println("Failed to Serve Data/Services:", err)
		os.Exit(1)
	}
}
