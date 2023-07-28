package main

import (
	"log"
	"net"
	"strings"
	"time"
)

const (
	udpAddr         = "192.168.178.1:12345" // Destination IP address and port
	secretMessage   = "NOCWARE ROCKS"       // Message to be encoded in Morse code
	morseUnitLength = 100                   // Message Payload length
	payloadMessage  = "edocesrom"           // morsecode spelled backwards
)

var morseCode = map[string]string{
	"A": ". -",
	"B": "- . . .",
	"C": "- . - .",
	"D": "- . .",
	"E": ".",
	"F": ". . - .",
	"G": "- - .",
	"H": ". . . .",
	"I": ". .",
	"J": ". - - -",
	"K": "- . -",
	"L": ". - . .",
	"M": "- -",
	"N": "- .",
	"O": "- - -",
	"P": ". - - .",
	"Q": "- - . -",
	"R": ". - .",
	"S": ". . .",
	"T": "-",
	"U": ". . -",
	"V": ". . . -",
	"W": ". - -",
	"X": "- . . -",
	"Y": "- . - -",
	"Z": "- - . .",
	" ": "    ", // represent a space with 7 Units
}

func main() {
	// Resolve UDP address
	udpAddr, err := net.ResolveUDPAddr("udp", udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	// Create UDP connection
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Encode the message in Morse code with word separation
	encodedMessage := encodeMessage(secretMessage)

	// Send each Morse code sign as a separate UDP packet
	for _, char := range encodedMessage {
		var packetCount int
		var payload bool = false
		if string(char) == "." {
			packetCount = 1
			payload = true
		}
		if string(char) == "-" {
			packetCount = 3
			payload = true
		}
		if string(char) == " " {
			packetCount = 1
			payload = false
		}

		sendPacket(conn, packetCount, payload)
		// Sleep for a short duration before sending the next packet
		time.Sleep(100 * time.Millisecond)

	}
	log.Println(encodedMessage)
	log.Println("Message sent successfully!")
}

func encodeMessage(secretMessage string) string {
	encoded := ""
	messageLength := len(secretMessage)

	for i, char := range secretMessage {
		strChar := string(char)
		encodedChar, found := morseCode[strChar]
		if !found {
			// Handle unsupported characters by logging an error
			log.Printf("Unsupported character '%s' found. Skipping.", strChar)
			continue
		}
		encoded += encodedChar
		if i < messageLength-1 && string(char) != " " {
			encoded += "   "
		}
	}

	// Remove the trailing space using strings.TrimSpace()
	encoded = strings.TrimSpace(encoded)

	return encoded
}

func sendPacket(conn *net.UDPConn, packetCount int, payload bool) error {
	// Check payload boolean
	var bytes []byte
	message := payloadMessage

	if payload {
		// If message is shorter than 100 bytes, repeat the message until it fills up 100 bytes
		for len(message) < morseUnitLength {
			message += message
		}
		// Slice the message down to 100 bytes if it's too long
		bytes = []byte(message[:morseUnitLength])
	} else {
		bytes = []byte("")
	}

	// Send the bytes as a UDP packet
	for packet := 0; packet < packetCount; packet++ {
		_, err := conn.Write(bytes)
		if err != nil {
			log.Println("Error sending packet:", err)
			return err
		}
	}

	return nil
}
