package mms

import (
	"encoding/json"
	supportFunctions "gomod/internal/support_functions"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
)

type MMSData struct {
	Country      string `json:"country"`
	Provider     string `json:"provider"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
}

func GetMms(mmsChan chan [][]MMSData) {
	respAns, err := http.Get("http://127.0.0.1:8383/mms")
	if err != nil {
		log.Println("Выполнить запрос не вышло")
		log.Println(err)
		mmsChan <- nil
		return
	}

	if respAns.StatusCode != 200 {
		var mmsData []MMSData
		log.Println(mmsData)
		mmsChan <- nil
		return
	}

	byteAns, err := ioutil.ReadAll(respAns.Body)
	if err != nil {
		log.Println("Распаковать запрос не удалось")
		log.Println(err)
		mmsChan <- nil
		return
	}

	var unsortedMmsCollection []MMSData
	if err := json.Unmarshal(byteAns, &unsortedMmsCollection); err != nil {
		log.Println("Заанмаршалить запрос не удалось")
		log.Println(err)
		mmsChan <- nil
		return
	}

	for i := 0; i < len(unsortedMmsCollection); i++ {
		if supportFunctions.CheckCountry(unsortedMmsCollection[i].Country) &&
			supportFunctions.CheckProvider(unsortedMmsCollection[i].Provider) &&
			unsortedMmsCollection[i].ResponseTime != "" &&
			unsortedMmsCollection[i].Bandwidth != "" {
			unsortedMmsCollection = unsortedMmsCollection[1 : len(unsortedMmsCollection)-1]
		}
	}
	sortedMmsCollection := [][]MMSData{}
	unsortedMmsCollection = makeFullCountry(unsortedMmsCollection)

	mmsSortedByProvider := sortByProvider(unsortedMmsCollection)
	sortedMmsCollection = append(sortedMmsCollection, mmsSortedByProvider)
	mmsSortedByCountry := sortByCountry(unsortedMmsCollection)
	sortedMmsCollection = append(sortedMmsCollection, mmsSortedByCountry)

	mmsChan <- sortedMmsCollection
}

func sortByProvider(unsortedMmsCollection []MMSData) []MMSData {
	mmsSortedByProvider := make([]MMSData, len(unsortedMmsCollection))
	copy(mmsSortedByProvider, unsortedMmsCollection)
	sort.Slice(mmsSortedByProvider, func(i, j int) bool {
		return mmsSortedByProvider[i].Provider < mmsSortedByProvider[j].Provider
	})
	return mmsSortedByProvider
}

func sortByCountry(unsortedMmsCollection []MMSData) []MMSData {
	mmsSortedByCountry := make([]MMSData, len(unsortedMmsCollection))
	copy(mmsSortedByCountry, unsortedMmsCollection)
	sort.Slice(mmsSortedByCountry, func(i, j int) bool {
		return mmsSortedByCountry[i].Country[:2] < mmsSortedByCountry[j].Country[:2]
	})
	return mmsSortedByCountry
}

func makeFullCountry(shortnameCountrySmsCollection []MMSData) []MMSData {
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
