package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"parking_lot/commando"
	"strconv"
	"strings"
)

func main() {
	flag.Parse()
	if len(flag.Args()) > 0 {
		executeFile(flag.Args()[0])
	} else {
		fmt.Println("Welcome to parking lot 1.4.2")
		fmt.Println("Available Commands:")
		fmt.Println("		- create_parking_lot <max_slots_num>")
		fmt.Println("		- park <car_reg_number> <car_colour>")
		fmt.Println("		- leave <slot_num>")
		fmt.Println("		- status")
		fmt.Println("		- registration_numbers_for_cars_with_colour <car_colour>")
		fmt.Println("		- slot_numbers_for_cars_with_colour <car_colour>")
		fmt.Println("		- slot_number_for_registration_number <car_reg_number>")
		fmt.Println("		- exit")
		executeInlineCommands()
	}
}

func executeInlineCommands() {
	exitCommand := false
	var buffReader *bufio.Reader
	buffReader = bufio.NewReader(os.Stdin)

	for !exitCommand {
		inputText, _ := buffReader.ReadString('\n')
		inputText = strings.TrimRight(inputText, "\r\n")
		if inputText == "exit" {
			fmt.Println("Thanks for using Parking Lot 1.4.2")
			break
		}
		runCmdInput(inputText)
	}
}

func runCmdInput(inputText string) {
	command := splitCommand(inputText)
	//fmt.Println(command)
	runCommand(command)
}

func executeFile(path string) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		fmt.Println("file not found")
		return
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		command := splitCommand(scanner.Text())
		runCommand(command)
	}
}

func runCommand(command []string) {
	switch command[0] {
	case "create_parking_lot":
		maxSlots, err := strconv.Atoi(command[1])
		if err != nil {
			panic(err.Error())
		}
		cpl := commando.CreateParkingLot(maxSlots)
		fmt.Println(cpl)
	case "park":
		p := commando.DoPark(command[1], command[2])
		fmt.Println(p)
	case "leave":
		no, err := strconv.Atoi(command[1])
		if err != nil {
			panic(err.Error())
		}
		l := commando.Leave(no)
		fmt.Println(l)
	case "status":
		status := commando.Status()
		fmt.Println("Slot No.    Registration No    Colour")
		for _, parkingCar := range status {
			fmt.Println(parkingCar)
		}
	case "registration_numbers_for_cars_with_colour":
		r := commando.RegistrationNumbersForCarsWithColour(command[1])
		fmt.Println(r)
	case "slot_numbers_for_cars_with_colour":
		ss := commando.SlotNumbersForCarsWithColour(command[1])
		fmt.Println(ss)
	case "slot_number_for_registration_number":
		s := commando.SlotNumberForRegistrationNumber(command[1])
		fmt.Println(s)
	default:
		fmt.Println("There is no such command available")
	}
}

func splitCommand(command string) []string {
	splitCommand := []string{}

	for _, s := range strings.Split(command, " ") {
		if s != "" {
			splitCommand = append(splitCommand, s)
		}
	}
	return splitCommand
}
