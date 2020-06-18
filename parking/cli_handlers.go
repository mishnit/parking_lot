package parking

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type StatusRes struct {
	Parks []Park `json:"Parks"`
}

type RegistrationNumbersRes struct {
	Cars []string `json:"Cars"`
}

type SlotNumbersRes struct {
	Slots []int `json:"Slots"`
}

type SlotNumberReqRes struct {
	SlotNum int `json:"SlotNum"`
}

type ParkingLotReq struct {
	MaxSlotsCount int `json:"MaxSlotsCount"`
}

type Res struct {
	Status string `json:"Status"`
}

type ParkReq struct {
	CarReg    string `json:"CarReg"`
	CarColour string `json:"CarColour"`
}

type ParkRes struct {
	SlotNum   int    `json:"SlotNum"`
	CarReg    string `json:"CarReg"`
	CarColour string `json:"CarColour"`
}

func Status() []string {
	resp, err := http.Get("http://localhost:3569/api/v1/getparks")
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var response StatusRes
	json.Unmarshal(bodyBytes, &response)

	var status []string
	for _, value := range response.Parks {
		slotNumber := strconv.FormatUint(uint64(value.SlotNum), 10)
		s := slotNumber + "\t" + value.CarReg + "\t" + value.CarColour
		status = append(status, s)
	}

	return status
}

func RegistrationNumbersForCarsWithColour(carColor string) string {
	resp, err := http.Get("http://localhost:3569/api/v1/getcarregs/colour/" + carColor)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var response RegistrationNumbersRes
	json.Unmarshal(bodyBytes, &response)

	//KA-01-HH-1234, KA-01-HH-9999, KA-01-P-333
	var s string
	for _, value := range response.Cars {
		carPlateNumber := value
		add := carPlateNumber + ", "
		s = s + add
	}
	status := trimSuffix(s, ", ")

	return status
}

func SlotNumbersForCarsWithColour(carColor string) string {
	resp, err := http.Get("http://localhost:3569/api/v1/getslots/colour/" + carColor)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var response SlotNumbersRes
	json.Unmarshal(bodyBytes, &response)

	var s string
	for _, value := range response.Slots {
		slotNumber := strconv.Itoa(value)
		add := slotNumber + ", "
		s = s + add
	}
	status := trimSuffix(s, ", ")

	return status
}

func SlotNumberForRegistrationNumber(carPlateNumber string) string {
	resp, err := http.Get("http://localhost:3569/api/v1/getslot/reg/" + carPlateNumber)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var response SlotNumberReqRes
	json.Unmarshal(bodyBytes, &response)

	var slotNumber string
	if response != (SlotNumberReqRes{}) {
		slotNumber = strconv.Itoa(response.SlotNum)
	} else {
		slotNumber = "Not found"
	}
	return slotNumber
}

func CreateParkingLot(s string) string {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalln(err)
	}
	lot := ParkingLotReq{n}
	jsonReq, err := json.Marshal(lot)
	resp, err := http.Post("http://localhost:3569/api/v1/createlot", "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var response Res
	json.Unmarshal(bodyBytes, &response)

	if response.Status != "created" {
		log.Fatalln(err)
	}

	return "Created a parking lot with " + s + " slots"
}

func DoPark(cpn string, cc string) string {
	park := ParkReq{cpn, cc}
	jsonReq, err := json.Marshal(park)
	resp, err := http.Post("http://localhost:3569/api/v1/postpark", "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var response ParkRes
	json.Unmarshal(bodyBytes, &response)

	ppc := response.SlotNum
	sn := strconv.Itoa(ppc)
	return "Allocated slot number: " + sn
}

func Leave(slotNumber int) string {

	slot := SlotNumberReqRes{slotNumber}
	jsonReq, err := json.Marshal(slot)
	resp, err := http.Post("http://localhost:3569/api/v1/postunpark", "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var response Res
	json.Unmarshal(bodyBytes, &response)

	if response.Status != "deleted" {
		log.Fatalln(err)
	}

	sn := strconv.Itoa(slotNumber)

	return "Slot number " + sn + " is free"
}
