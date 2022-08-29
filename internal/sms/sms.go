package sms

import (
	"bufio"
	supportFunctions "gomod/internal/support_functions"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type SMSData struct {
	Country      string `json:"country"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
	Provider     string `json:"provider"`
}

func GetSms(smsChan chan [][]SMSData) {

	smsDataFile, err := os.Open("simulator/skillbox-diploma/sms.data")
	if err != nil {
		log.Println("Файл отсутсвует или пуст")
		log.Println(err)
		smsChan <- nil
		return
	}
	defer smsDataFile.Close()

	unsortedSmsCollection := []SMSData{}

	scanner := bufio.NewScanner(smsDataFile)
	for scanner.Scan() {
		line := scanner.Text()
		lineSlice := strings.Split(line, ";")
		bandwidth, _ := strconv.Atoi(lineSlice[1])

		if len(lineSlice) == 4 &&
			lineSlice[2] != "" &&
			supportFunctions.CheckCountry(lineSlice[0]) &&
			supportFunctions.CheckProvider(lineSlice[3]) &&
			(bandwidth >= 0 && bandwidth <= 100) {

			correctLine := SMSData{lineSlice[0], lineSlice[1], lineSlice[2], lineSlice[3]}
			unsortedSmsCollection = append(unsortedSmsCollection, correctLine)
		}
	}
	var sortedSmsCollection [][]SMSData
	smsData := makeFullCountry(unsortedSmsCollection)

	smsSortedByProvider := sortByProvider(smsData)
	sortedSmsCollection = append(sortedSmsCollection, smsSortedByProvider)
	smsSortedByCountry := sortByCountry(smsData)
	sortedSmsCollection = append(sortedSmsCollection, smsSortedByCountry)
	smsChan <- sortedSmsCollection
}

func sortByProvider(unsortedSmsCollection []SMSData) []SMSData {
	smsSortedByProvider := make([]SMSData, len(unsortedSmsCollection))
	copy(smsSortedByProvider, unsortedSmsCollection)
	sort.Slice(smsSortedByProvider, func(i, j int) bool {
		return smsSortedByProvider[i].Provider < smsSortedByProvider[j].Provider
	})
	return smsSortedByProvider
}

func sortByCountry(unsortedSmsCollection []SMSData) []SMSData {
	smsSortedByCountry := make([]SMSData, len(unsortedSmsCollection))
	copy(smsSortedByCountry, unsortedSmsCollection)
	sort.Slice(smsSortedByCountry, func(i, j int) bool {
		return smsSortedByCountry[i].Country[:2] < smsSortedByCountry[j].Country[:2]
	})
	return smsSortedByCountry
}

func makeFullCountry(shortnameCountrySmsCollection []SMSData) []SMSData {
	for i := 0; i < len(shortnameCountrySmsCollection); i++ {
		switch {
		case shortnameCountrySmsCollection[i].Country == "RU":
			shortnameCountrySmsCollection[i].Country = "Russia"
		case shortnameCountrySmsCollection[i].Country == "US":
			shortnameCountrySmsCollection[i].Country = "USA"
		case shortnameCountrySmsCollection[i].Country == "GB":
			shortnameCountrySmsCollection[i].Country = "Great Britain"
		case shortnameCountrySmsCollection[i].Country == "FR":
			shortnameCountrySmsCollection[i].Country = "France"
		case shortnameCountrySmsCollection[i].Country == "BL":
			shortnameCountrySmsCollection[i].Country = "Saint-Barthelemy"
		case shortnameCountrySmsCollection[i].Country == "AT":
			shortnameCountrySmsCollection[i].Country = "Austria"
		case shortnameCountrySmsCollection[i].Country == "BG":
			shortnameCountrySmsCollection[i].Country = "Bulgaria"
		case shortnameCountrySmsCollection[i].Country == "DK":
			shortnameCountrySmsCollection[i].Country = "Denmark"
		case shortnameCountrySmsCollection[i].Country == "CA":
			shortnameCountrySmsCollection[i].Country = "Canada"
		case shortnameCountrySmsCollection[i].Country == "ES":
			shortnameCountrySmsCollection[i].Country = "Spain"
		case shortnameCountrySmsCollection[i].Country == "CH":
			shortnameCountrySmsCollection[i].Country = "Switzerland"
		case shortnameCountrySmsCollection[i].Country == "TR":
			shortnameCountrySmsCollection[i].Country = "Turkey"
		case shortnameCountrySmsCollection[i].Country == "PE":
			shortnameCountrySmsCollection[i].Country = "Peru"
		case shortnameCountrySmsCollection[i].Country == "NZ":
			shortnameCountrySmsCollection[i].Country = "New Zealand"
		case shortnameCountrySmsCollection[i].Country == "MC":
			shortnameCountrySmsCollection[i].Country = "Monaco"
		}
	}
	return shortnameCountrySmsCollection
}
