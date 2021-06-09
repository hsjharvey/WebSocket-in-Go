package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	// Maximum message size
	maxMessageSize     = 1024 * 10
	subjectsIDFilePath = "./src/server/input/id_verification_actual.json"
	practiceIDFilePath = "./src/server/input/id_verification_practice.json"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}

	subjectIDs  = map[string]interface{}{}
	practiceIDs = map[string]interface{}{}

	playData     = make(map[string]interface{})
	gameInfoData = make(map[string]interface{})

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024 * 256,
		WriteBufferSize: 1024 * 256,
	}
)

type msgToClient struct {
	MsgType string      `json:"msg_type"`
	Msg     interface{} `json:"msg"`
}

func loadSubjectIds() {
	jsonFile, err := os.Open(subjectsIDFilePath)
	if err != nil {
		log.Fatalf("error in reading files! %s", err)
	}
	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal([]byte(strings.TrimSuffix(string(byteValue), "\r\n")), &subjectIDs)
	if err != nil {
		log.Fatalf("error in unmarshal IDs! %s", err)
	}

	jsonFile, err = os.Open(practiceIDFilePath)
	if err != nil {
		log.Fatalf("error in reading files! %s", err)
	}
	// read our opened jsonFile as a byte array.
	byteValue, _ = ioutil.ReadAll(jsonFile)

	err = json.Unmarshal([]byte(strings.TrimSuffix(string(byteValue), "\r\n")), &practiceIDs)
	if err != nil {
		log.Fatalf("error in unmarshal IDs! %s", err)
	}

	log.Println("finish loading user verification information")
}

func (c *Client) saveClientData() {
	jsonString, err := json.MarshalIndent(c.playData, "", "  ")
	if err != nil {
		log.Fatalf("error in writing play data files! #{err}")
	} else {
		ioutil.WriteFile("output/"+c.ID+"_play_data.json", jsonString, os.ModePerm)
		log.Println("play data saved.")
	}

	jsonString, err = json.MarshalIndent(c.gameInfo, "", "   ")
	if err != nil {
		log.Fatalf("error in writing game info data files! #{err}")
	} else {
		ioutil.WriteFile("output/"+c.ID+"_game_info.json", jsonString, os.ModePerm)
		log.Println("game info data saved.")
	}
}
