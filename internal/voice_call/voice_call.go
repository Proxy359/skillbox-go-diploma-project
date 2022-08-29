package voice

import (
	"bufio"
	supportFunctions "gomod/internal/support_functions"
	"log"
	"os"
	"strconv"
	"strings"
)

type VoiceCallData struct {
	Country             string  `json:"country"`
	Bandwidth           string  `json:"bandwidth"`
	ResponseTime        string  `json:"response_time"`
	Provider            string  `json:"provider"`
	ConnectionStability float32 `json:"connection_stability"`
	TTFB                int     `json:"ttfb"`
	VoicePurity         int     `json:"voice_purity"`
	MedianOfCallsTime   int     `json:"median_of_calls_time"`
}

func GetVoice(voiceCallChan chan []VoiceCallData) {
	voiceDataFile, err := os.Open("simulator/skillbox-diploma/voice.data")
	if err != nil {
		log.Println("Файл отсутсвует или пуст")
		log.Println(err)
		voiceCallChan <- nil
		return
	}
	defer voiceDataFile.Close()

	voiceCallCollection := []VoiceCallData{}

	scanner := bufio.NewScanner(voiceDataFile)
	for scanner.Scan() {
		line := scanner.Text()
		lineSlice := strings.Split(line, ";")
		bandwidth, _ := strconv.Atoi(lineSlice[1])
		connectionStability, _ := strconv.ParseFloat(lineSlice[4], 32)
		connectionStability32 := float32(connectionStability)
		TTFB, _ := strconv.Atoi(lineSlice[5])
		VoicePurity, _ := strconv.Atoi(lineSlice[6])
		MedianOfCallsTime, _ := strconv.Atoi(lineSlice[7])

		if len(lineSlice) == 8 &&
			lineSlice[2] != "" &&
			supportFunctions.CheckCountry(lineSlice[0]) &&
			supportFunctions.CheckProvider(lineSlice[3]) &&
			(bandwidth >= 0 && bandwidth <= 100) &&
			connectionStability32 != 0 &&
			TTFB != 0 &&
			VoicePurity != 0 &&
			MedianOfCallsTime != 0 {
			correctLine := VoiceCallData{lineSlice[0], lineSlice[1], lineSlice[2], lineSlice[3], connectionStability32, TTFB, VoicePurity, MedianOfCallsTime}
			voiceCallCollection = append(voiceCallCollection, correctLine)
		}
	}
	voiceCallChan <- voiceCallCollection
}
