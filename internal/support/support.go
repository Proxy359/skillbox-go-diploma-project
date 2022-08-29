package support

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type SupportData struct {
	Topic         string `json:"topic"`
	ActiveTickets int    `json:"active_tickets"`
}

func GetSupport(supportChan chan []int) {
	respAns, err := http.Get("http://127.0.0.1:8383/support")
	if err != nil {
		log.Println("Выполнить запрос не вышло")
		log.Println(err)
		supportChan <- nil
		return
	}

	if respAns.StatusCode != 200 {
		var supportDataCollection []SupportData
		log.Println(supportDataCollection)
		supportChan <- nil
		return
	}

	byteAns, err := ioutil.ReadAll(respAns.Body)
	if err != nil {
		log.Println("Распаковать запрос не удалось")
		log.Println(err)
		supportChan <- nil
		return
	}

	var supportDataCollection []SupportData
	if err := json.Unmarshal(byteAns, &supportDataCollection); err != nil {
		log.Println("Заанмаршалить запрос не удалось")
		log.Println(err)
		supportChan <- nil
		return
	}

	load := 0
	for i := 0; i < len(supportDataCollection); i++ {
		load = load + supportDataCollection[i].ActiveTickets
	}
	log.Println(load)
	loadLvl := 0
	if load < 9 {
		loadLvl = 1
	} else if load > 16 {
		loadLvl = 3
	} else {
		loadLvl = 2
	}

	supportStatus := []int{}
	supportStatus = append(supportStatus, loadLvl)

	if load%2 == 1 {
		ans := (float64(load)/2 + 1) * 3.333333333333333
		if ans-float64(int(ans)) == 0 {
			supportStatus = append(supportStatus, int(ans))
		} else {
			supportStatus = append(supportStatus, int(ans+1))
		}
	} else {
		ans := (float64(load)/2 + 0.5) * 3.333333333333333
		if ans-float64(int(ans)) == 0 {
			supportStatus = append(supportStatus, int(ans))
		} else {
			supportStatus = append(supportStatus, int(ans+1))
		}
	}
	supportChan <- supportStatus
}
