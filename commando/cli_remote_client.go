package commando

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Park struct {
	SlotNum   uint32 `json:"SlotNum"`
	CarReg    string `json:"CarReg"`
	CarColour string `json:"CarColour"`
}

type StatusRes struct {
	Parks  []Park `json:"Parks"`
	Status string `json:"Status"`
}

type RegistrationNumbersRes struct {
	Cars   []string `json:"Cars"`
	Status string   `json:"Status"`
}

type SlotNumbersRes struct {
	Slots  []int  `json:"Slots"`
	Status string `json:"Status"`
}

type SlotNumberReq struct {
	SlotNum int `json:"SlotNum"`
}

type SlotNumberRes struct {
	SlotNum int    `json:"SlotNum"`
	Status  string `json:"Status"`
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
	Park   Park   `json:"Park"`
	Status string `json:"Status"`
}

func CreateParkingLot(n int) string {
	if n <= 0 {
		return "Lot size cannot be <= zero"
	}

	lot := ParkingLotReq{n}
	jsonReq, err := json.Marshal(lot)
	resp, err := http.Post("http://localhost:3569/api/v1/createlot", "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))

	if err != nil {
		log.Fatalln(err)
		return "Error"
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var response Res
	json.Unmarshal(bodyBytes, &response)

	if response.Status == "ErrLotSizeLTEZero" {
		return "Lot size cannot be <= zero"
	}

	if response.Status == "Error" {
		return "Unexpected error occured"
	}

	return "Created a parking lot with " + strconv.Itoa(n) + " slots"
}

func DoPark(cpn string, cc string) string {
	park := ParkReq{cpn, cc}
	jsonReq, err := json.Marshal(park)
	resp, err := http.Post("http://localhost:3569/api/v1/postpark", "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
	if err != nil {
		log.Fatalln(err)
		return "Error"
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var response ParkRes
	json.Unmarshal(bodyBytes, &response)

	if response.Status == "ErrInvalidCarNumber" {
		return "Invalid indian car number plate format"
	}

	if response.Status == "ErrNoLotFound" {
		return "No lot available, please create a lot first"
	}

	if response.Status == "ErrParkingFull" {
		return "Sorry, parking lot is full"
	}

	if response.Status == "Error" {
		return "Unexpected error occured"
	}

	ppc := response.Park.SlotNum
	sn := strconv.Itoa(int(ppc))
	return "Allocated slot number: " + sn
}

func Leave(slotNumber int) string {
	if slotNumber <= 0 {
		return "Slot invalid"
	}

	slot := SlotNumberReq{slotNumber}
	jsonReq, err := json.Marshal(slot)
	resp, err := http.Post("http://localhost:3569/api/v1/postunpark", "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))

	if err != nil {
		log.Fatalln(err)
		return "Error"
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var response Res
	json.Unmarshal(bodyBytes, &response)

	if response.Status == "ErrInvalidSlot" {
		return "Slot invalid"
	}

	if response.Status == "ErrNoLotFound" {
		return "No lot available, please create a lot first"
	}

	if response.Status == "ErrParking" {
		return "Parking slot is empty"
	}

	if response.Status == "Error" {
		return "Unexpected error occured"
	}

	sn := strconv.Itoa(slotNumber)

	return "Slot number " + sn + " is free"
}

func Status() []string {
	var status []string
	resp, err := http.Get("http://localhost:3569/api/v1/getparks")
	if err != nil {
		log.Fatalln(err)
		status = append(status, "Error")
		return status
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var response StatusRes
	json.Unmarshal(bodyBytes, &response)

	if response.Status == "ErrNoLotFound" {
		status = append(status, "No lot available, please create a lot first")
		return status
	}

	if response.Status == "ErrNotFound" {
		status = append(status, "Not found")
		return status
	}

	if response.Status == "Error" {
		status = append(status, "Unexpected error occured")
		return status
	}

	for _, value := range response.Parks {
		slotNumber := strconv.FormatUint(uint64(value.SlotNum), 10)
		s := slotNumber + "           " + value.CarReg + "      " + value.CarColour
		status = append(status, s)
	}

	return status
}

func RegistrationNumbersForCarsWithColour(carColor string) string {
	resp, err := http.Get("http://localhost:3569/api/v1/getcarregs/colour/" + carColor)
	if err != nil {
		log.Fatalln(err)
		return "Error"
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var response RegistrationNumbersRes
	json.Unmarshal(bodyBytes, &response)

	if response.Status == "ErrNoLotFound" {
		return "No lot available, please create a lot first"
	}

	if response.Status == "ErrNotFound" {
		return "Not found"
	}

	if response.Status == "Error" {
		return "Unexpected error occured"
	}

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
		return "Error"
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var response SlotNumbersRes
	json.Unmarshal(bodyBytes, &response)

	if response.Status == "ErrNoLotFound" {
		return "No lot available, please create a lot first"
	}

	if response.Status == "ErrNotFound" {
		return "Not found"
	}

	if response.Status == "Error" {
		return "Unexpected error occured"
	}

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
		return "Error"
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var response SlotNumberRes
	json.Unmarshal(bodyBytes, &response)

	if response.Status == "ErrInvalidCarNumber" {
		return "Invalid indian car number plate format"
	}

	if response.Status == "ErrNotFound" {
		return "Not found"
	}

	if response.Status == "Error" {
		return "Unexpected error occured"
	}

	sn := strconv.Itoa(response.SlotNum)

	return sn
}

func trimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}
