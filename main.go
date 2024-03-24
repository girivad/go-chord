package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	overlay "github.com/girivad/go-chord/Overlay"
)

func main() {
	// Read cmdline args (ip + contact address if joining an existing chord ring)
	// contactPtr := flag.String("contact", "None", "IP address of a contact in a Chord Ring you want to join.")
	// flag.Parse()

	if len(os.Args) < 4 {
		log.Println("Please provide this node's IP address and the capacity of its chord ring.")
		os.Exit(1)
	}

	ip := os.Args[1]
	capacity, err := strconv.ParseInt(os.Args[2], 10, 64)

	if err != nil {
		fmt.Println("Error Reading Capacity:", err)
		os.Exit(1)
	}

	contact := os.Args[3]

	// TO-DO: Implement IP Verification (verify that this is a valid IP address, at least via Regex)

	// Create a new ChordNode and join an existing chord ring if requested.
	chordServer, err := overlay.NewChordServer(ip, capacity)

	if err != nil {
		log.Printf("[FATAL] Failed to connect to self.")
		os.Exit(1)
	}

	if contact != "None" {
		// Make a client for the contact, and then run a join service on it.
		log.Printf("[INFO] Contact in the Chord Ring: %s", contact)
		contactNode, err := overlay.Connect(contact)

		if err != nil {
			log.Println("Unable to connect to the contact in the Chord Ring:", err)
			os.Exit(1)
		}

		err = chordServer.Join(contactNode)

		if err != nil {
			log.Println("Failed to join chord ring:", err)
			os.Exit(1)
		}
	}

	// Serve data from 8080, gRPC through 8081.
	go func() {
		err = chordServer.Serve()
		if err != nil {
			log.Println("Failed to Serve Data/Services:", err)
			os.Exit(1)
		}
	}()

	select {}
}
