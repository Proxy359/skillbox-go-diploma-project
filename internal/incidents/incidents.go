package incidents

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type IncidentData struct {
	Topic  string `json:"topic"`
	Status string `json:"status"`
}

func GetIncidents(incidentsChan chan []IncidentData) {
	respAns, err := http.Get("http://127.0.0.1:8383/accendent")
	if err != nil {
		log.Println("Выполнить запрос не вышло")
		log.Println(err)
		incidentsChan <- nil
		return
	}

	if respAns.StatusCode != 200 {
		var IncidentDataCollection []IncidentData
		log.Println(IncidentDataCollection)
		incidentsChan <- nil
		return
	}

	byteAns, err := ioutil.ReadAll(respAns.Body)
	if err != nil {
		log.Println("Распаковать запрос не удалось")
		log.Println(err)
		incidentsChan <- nil
		return
	}

	var incidentDataCollection []IncidentData
	if err := json.Unmarshal(byteAns, &incidentDataCollection); err != nil {
		log.Println("Заанмаршалить запрос не удалось")
		log.Println(err)
		incidentsChan <- nil
		return
	}

	incidentsChan <- incidentDataCollection
}
