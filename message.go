package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/pkg/errors"
)

type Message struct {
	To   Friend   `json:"to"`
	From Friend   `json:"from"`
	CC   []Friend `json:"cc"`
	Body string   `json:"body"`
}

func (msg Message) String() string {
	return fmt.Sprintf("\n\tFrom: %s\n\tTo: %s\n\tCC: %s\n\tBody: %s\n\n", msg.From, msg.To, msg.CC, msg.Body)
}

func (msg *Message) Send() {
	if msg.Body == "" {
		log.Println("Ignoring empty message")
		return
	}

	log.Print(msg)
	// Connect to my friend
	log.Printf("Connecting to my friend %s(%s)\n", msg.To.Name, msg.To.Number)
	conn, err := net.DialTCP("tcp", nil, msg.To.Number.TCPAddr)
	if err != nil {
		err = errors.Wrapf(err, "Unable to connect to %s(%s)", msg.To.Name, msg.To.Number)
		log.Printf("%+v", err)
		return
	}
	defer conn.Close()

	// Marshal the message to bytes
	msgb, err := json.Marshal(msg)
	if err != nil {
		err = errors.Wrapf(err, "Unable to marshal the message %s", msg.To.Name, msg.To.Number)
		log.Printf("%+v", err)
		return
		log.Printf(": %v", err)
	}

	log.Println("Sending message")
	_, err = conn.Write(msgb)
	if err != nil {
		log.Fatalf("Unable to send the message: %v", err)
	}
}

func (msg *Message) Forward(body string) {
	if len(msg.CC) == 0 {
		// The message has finally been sent to everyone
		return
	}

	reply := Message{
		From: msg.To,
		To:   msg.CC[0],
		CC:   msg.CC[1:len(msg.CC)],
		Body: body,
	}

	reply.Send()
}

func readMessage(reader io.ReadCloser) (Message, error) {
	defer reader.Close()

	log.Println("Parsing incoming message...")

	msg := Message{}
	err := json.NewDecoder(reader).Decode(&msg)
	return msg, errors.Wrap(err, "Unable to parse incomming message")
}