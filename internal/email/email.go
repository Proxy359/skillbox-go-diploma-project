package email

import (
	"bufio"
	supportFunctions "gomod/internal/support_functions"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type EmailData struct {
	Country      string `json:"country"`
	Provider     string `json:"provider"`
	DeliveryTime int    `json:"delivery_time"`
}

func GetEmail(emailChan chan map[string][][]EmailData) {

	smsDataFile, err := os.Open("simulator/skillbox-diploma/email.data")
	if err != nil {
		log.Println("Файл отсутсвует или пуст")
		log.Println(err)
		emailChan <- nil
		return
	}
	defer smsDataFile.Close()

	emailCollection := []EmailData{}

	scanner := bufio.NewScanner(smsDataFile)
	for scanner.Scan() {
		line := scanner.Text()
		lineSlice := strings.Split(line, ";")

		if len(lineSlice) == 3 &&
			lineSlice[2] != "" &&
			supportFunctions.CheckCountry(lineSlice[0]) &&
			checkProvider(lineSlice[1]) {

			deliveryTime, _ := strconv.Atoi(lineSlice[2])
			correctLine := EmailData{lineSlice[0], lineSlice[1], deliveryTime}
			emailCollection = append(emailCollection, correctLine)
		}
	}
	countrysMap := make(map[string][]EmailData)
	countrysMap = createCountrysMap(countrysMap, emailCollection)
	for key := range countrysMap {
		countrysMap[key] = filterByProvider(countrysMap[key])
	}

	newMap := make(map[string][][]EmailData)
	for key := range countrysMap {
		fastestProviders := countrysMap[key][:3]
		slovestProviders := countrysMap[key][len(countrysMap[key])-4 : len(countrysMap[key])-1]
		cercleSloce := [][]EmailData{}
		cercleSloce = append(cercleSloce, fastestProviders)
		cercleSloce = append(cercleSloce, slovestProviders)

		newMap[key] = cercleSloce
	}

	emailChan <- newMap
}

func filterByProvider(emailSlise []EmailData) []EmailData {
	sort.SliceStable(emailSlise, func(i, j int) bool {
		return emailSlise[i].DeliveryTime > emailSlise[j].DeliveryTime
	})
	return emailSlise
}

func createCountrysMap(mapOfCountrys map[string][]EmailData, collectionOfEmail []EmailData) map[string][]EmailData {
	for i := 0; i < len(collectionOfEmail); i++ {
		cercleCountry := collectionOfEmail[i].Country
		mapOfCountrys[cercleCountry] = append(mapOfCountrys[cercleCountry], collectionOfEmail[i])
	}
	return mapOfCountrys
}

func checkProvider(str string) bool {
	providersList := [13]string{"Gmail", "Yahoo", "Hotmail", "MSN", "Orange", "Comcast", "AOL", "Live", "RediffMail", "GMX", "Protonmail",
		"Yandex", "Mail.ru"}
	for i := 0; i < len(providersList); i++ {
		if str == providersList[i] {
			return true
		}
	}
	return false
}
