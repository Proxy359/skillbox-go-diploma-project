package app

import (
	"encoding/json"
	"gomod/internal/billing"
	"gomod/internal/email"
	"gomod/internal/incidents"
	"gomod/internal/mms"
	"gomod/internal/sms"
	"gomod/internal/support"
	voice "gomod/internal/voice_call"
	"net/http"

	"github.com/gorilla/mux"
)

type ResultT struct {
	Status bool       `json:"status"`
	Data   ResultSetT `json:"data"`
	Error  string     `json:"error"`
}

type ResultSetT struct {
	SMS       [][]sms.SMSData                `json:"sms"`
	MMS       [][]mms.MMSData                `json:"mms"`
	VoiceCall []voice.VoiceCallData          `json:"voice_call"`
	Email     map[string][][]email.EmailData `json:"email"`
	Billing   billing.BillingData            `json:"billing"`
	Support   []int                          `json:"support"`
	Incidents []incidents.IncidentData       `json:"incident"`
}

var (
	smsChan       = make(chan [][]sms.SMSData)
	mmsChan       = make(chan [][]mms.MMSData)
	voiceCallChan = make(chan []voice.VoiceCallData)
	emailChan     = make(chan map[string][][]email.EmailData)
	billingChan   = make(chan billing.BillingData)
	supportChan   = make(chan []int)
	incidentsChan = make(chan []incidents.IncidentData)
)

func Start() {
	router := mux.NewRouter()
	router.HandleFunc("/", handleConnection).Methods("GET", "OPTIONS")
	http.ListenAndServe("127.0.0.1:8282", router)
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	resultAns := getResultData()
	byteResultAns, err := json.Marshal(resultAns)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(byteResultAns)
}

func getResultData() ResultT {

	resultAns := ResultT{false, ResultSetT{}, "Error on collect data"}

	go sms.GetSms(smsChan)
	go mms.GetMms(mmsChan)
	go voice.GetVoice(voiceCallChan)
	go email.GetEmail(emailChan)
	go billing.GetBilling(billingChan)
	go support.GetSupport(supportChan)
	go incidents.GetIncidents(incidentsChan)

	if incidentsChan != nil {

		smsData := <-smsChan
		if len(smsData) == 0 {
			return resultAns
		}

		mmsData := <-mmsChan
		if len(mmsData) == 0 {
			return resultAns
		}

		voiceData := <-voiceCallChan
		if len(voiceData) == 0 {
			return resultAns
		}

		emailData := <-emailChan
		if len(emailData) == 0 {
			return resultAns
		}

		checkBillingData := billing.BillingData{}
		billingData := <-billingChan
		if billingData == checkBillingData {
			return resultAns
		}

		supportData := <-supportChan
		if len(supportData) == 0 {
			return resultAns
		}

		incidentsData := <-incidentsChan
		if len(incidentsData) == 0 {
			return resultAns
		}

		resultSetAns := ResultSetT{smsData, mmsData, voiceData, emailData, billingData, supportData, incidentsData}

		resultAns = ResultT{true, resultSetAns, ""}
	}
	return resultAns
}
