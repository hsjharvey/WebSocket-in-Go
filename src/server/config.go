package main

import (
	"github.com/gorilla/websocket"
)

const (
	// Maximum message size
	maxMessageSize = 1024 * 10
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 256,
	WriteBufferSize: 1024 * 256,
}

const (
	subjectsIDFilePath = "src/id_verification_actual.json"
)

type msgToSubject struct {
	MsgType string      `json:"msg_type"`
	Msg     interface{} `json:"msg"`
}

type message struct {
	ID   string
	data []byte
}

var subjectIDs = map[string]interface{}{}
var practiceIDs = map[string]interface{}{}
