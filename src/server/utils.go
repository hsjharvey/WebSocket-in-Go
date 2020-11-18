package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func loadSubjectIds() {
	jsonFile, err := os.Open(subjectsIDFilePath)
	if err != nil {
		log.Fatalf("error in reading files! %s", err)
	} else {
		log.Println("finish loading user verification information")
	}
	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal([]byte(strings.TrimSuffix(string(byteValue), "\r\n")), &subjectIDs)
	if err != nil {
		log.Fatalf("error in unmarshal IDs! %s", err)
	}
}

func (c *Client) saveSubjectData() {
	jsonString, err := json.Marshal(c.playData)
	if err != nil {
		log.Fatalf("error in writing files! #{err}")
	} else {
		ioutil.WriteFile(strings.Join([]string{"src/", c.ID, "_play_data.json"}, ","), jsonString, os.ModePerm)
		log.Println("play data saved.")
	}

	jsonString, err = json.Marshal(c.gameInfo)
	if err != nil {
		log.Fatalf("error in writing files! #{err}")
	} else {
		ioutil.WriteFile(strings.Join([]string{"src/", c.ID, "_game_info.json"}, ","), jsonString, os.ModePerm)
		log.Println("game info data saved.")
	}
}
