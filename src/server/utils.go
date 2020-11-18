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
	subjectsIDFilePath = "src/id_verification_actual.json"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}

	subjectIDs  = map[string]interface{}{}
	practiceIDs = map[string]interface{}{}

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024 * 256,
		WriteBufferSize: 1024 * 256,
	}
)

type msgToWebGL struct {
	MsgType string      `json:"msg_type"`
	Msg     interface{} `json:"msg"`
}

type playData struct {
	currentSession          string `json:"currentSession"`
	currentPart             string `json:"currentPart"`
	randomSeed              string `json:"randomSeed"`
	eventType               string `json:"eventType"`
	participantChoice       string `json:"participantChoice"`
	displaySampleLeftValue  string `json:"displaySampleLeftValue"`
	displaySampleRightValue string `json:"displaySampleRightValue"`
	leftSliderMin           string `json:"leftSliderMin"`
	leftSliderMax           string `json:"leftSliderMax"`
	rightSliderMin          string `json:"rightSliderMin"`
	rightSliderMax          string `json:"rightSliderMax"`
	sliderStartPosition     string `json:"sliderStartPosition"`
	participantEstimation   string `json:"participantEstimation"`
	predictionError         string `json:"predictionError"`
	riskFreeMean            string `json:"riskFreeMean"`
	time                    string `json:"time"`
}

type gameData struct {
	gameOrder                        string `json:"gameOrder"`
	gaussianFinalChoiceReward        string `json:"gaussianFinalChoiceReward"`
	student_tFinalChoiceReward       string `json:"student_tFinalChoiceReward"`
	exponentialFinalChoiceReward     string `json:"exponentialFinalChoiceReward"`
	exponentialAvgPredictionError    string `json:"exponentialAvgPredictionError"`
	gaussianRewardCalculation        string `json:"gaussianRewardCalculation"`
	student_tRewardCalculation       string `json:"student_tRewardCalculation"`
	exponentialRewardCalculation     string `json:"exponentialRewardCalculation"`
	finalScaledReward                string `json:"finalScaledReward"`
	choiceRewardToPaymentScaleFactor string `json:"choiceRewardToPaymentScaleFactor"`
	finalPayment                     string `json:"finalPayment"`
	gaussianpart_1                   string `json:"gaussianpart_1"`
	gaussianpart_2                   string `json:"gaussianpart_2"`
	student_tpart_1                  string `json:"student_tpart_1"`
	student_tpart_2                  string `json:"student_tpart_2"`
	exponentialpart_1                string `json:"exponentialpart_1"`
	exponentialpart_2                string `json:"exponentialpart_2"`
}

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
		log.Fatalf("error in writing play data files! #{err}")
	} else {
		ioutil.WriteFile(strings.Join([]string{"output/", c.ID, "_play_data.json"}, ","), jsonString, os.ModePerm)
		log.Println("play data saved.")
	}

	jsonString, err = json.Marshal(c.gameInfo)
	if err != nil {
		log.Fatalf("error in writing game info data files! #{err}")
	} else {
		ioutil.WriteFile(strings.Join([]string{"output/", c.ID, "_game_info.json"}, ","), jsonString, os.ModePerm)
		log.Println("game info data saved.")
	}
}
