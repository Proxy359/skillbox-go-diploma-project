package billing

import (
	"bufio"
	"log"
	"math"
	"os"
	"strconv"
)

type BillingData struct {
	CreateCustomer bool
	Purchase       bool
	Payout         bool
	Recurring      bool
	FraudControl   bool
	CheckoutPage   bool
}

func GetBilling(billingChan chan BillingData) {

	billingDataFile, err := os.Open("simulator/skillbox-diploma/billing.data")
	if err != nil {
		log.Println("Файл отсутсвует или пуст")
		log.Println(err)
		billingChan <- BillingData{}
		return
	}
	defer billingDataFile.Close()

	scanner := bufio.NewScanner(billingDataFile)
	var elem BillingData
	for scanner.Scan() {
		line := scanner.Text()

		newSlice := []int{}
		for i := len(line); i > 0; i-- {
			cercleInt, _ := strconv.Atoi(line[i-1:])
			newSlice = append(newSlice, cercleInt)
			line = line[:i-1]
		}

		var digit uint8
		for i := 0; i < len(newSlice); i++ {
			digit = digit + uint8(newSlice[i])*uint8(math.Pow(2, float64(i)))
		}

		CreateCustomer := digit&1 == 1
		Purchase := digit>>1&1 == 1
		Payout := digit>>2&1 == 1
		Recurring := digit>>3&1 == 1
		FraudControl := digit>>4&1 == 1
		CheckoutPage := digit>>5&1 == 1

		elem = BillingData{
			CreateCustomer: CreateCustomer,
			Purchase:       Purchase,
			Payout:         Payout,
			Recurring:      Recurring,
			FraudControl:   FraudControl,
			CheckoutPage:   CheckoutPage,
		}
	}
	billingChan <- elem
}
