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
		fmt.Println("		- help")
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

	if err != nil {
		fmt.Println("file not found")
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		command := splitCommand(scanner.Text())
		runCommand(command)
	}
}

func runCommand(command []string) {
	switch command[0] {
	case "create_parking_lot":
		if len(command) == 2 {
			maxSlots, err := strconv.Atoi(command[1])

			if err != nil {
				panic(err.Error())
			}

			cpl := commando.CreateParkingLot(maxSlots)
			fmt.Println(cpl)
		} else {
			fmt.Println("Invalid command")
		}
	case "park":
		if len(command) == 3 {
			p := commando.DoPark(command[1], command[2])
			fmt.Println(p)
		} else {
			fmt.Println("Invalid command")
		}
	case "leave":
		if len(command) == 2 {
			no, err := strconv.Atoi(command[1])
			if err != nil {
				panic(err.Error())
			}
			l := commando.Leave(no)
			fmt.Println(l)
		} else {
			fmt.Println("Invalid command")
		}
	case "status":
		if len(command) == 1 {
			status := commando.Status()
			fmt.Println("Slot No.    Registration No    Colour")
			for _, parkingCar := range status {
				fmt.Println(parkingCar)
			}
		} else {
			fmt.Println("Invalid command")
		}
	case "registration_numbers_for_cars_with_colour":
		if len(command) == 2 {
			r := commando.RegistrationNumbersForCarsWithColour(command[1])
			fmt.Println(r)
		} else {
			fmt.Println("Invalid command")
		}
	case "slot_numbers_for_cars_with_colour":
		if len(command) == 2 {
			ss := commando.SlotNumbersForCarsWithColour(command[1])
			fmt.Println(ss)
		} else {
			fmt.Println("Invalid command")
		}
	case "slot_number_for_registration_number":
		if len(command) == 2 {
			s := commando.SlotNumberForRegistrationNumber(command[1])
			fmt.Println(s)
		} else {
			fmt.Println("Invalid command")
		}
	case "help":
		fmt.Println("Available Commands:")
		fmt.Println("		- create_parking_lot <max_slots_num>")
		fmt.Println("		- park <car_reg_number> <car_colour>")
		fmt.Println("		- leave <slot_num>")
		fmt.Println("		- status")
		fmt.Println("		- registration_numbers_for_cars_with_colour <car_colour>")
		fmt.Println("		- slot_numbers_for_cars_with_colour <car_colour>")
		fmt.Println("		- slot_number_for_registration_number <car_reg_number>")
		fmt.Println("		- help")
		fmt.Println("		- exit")
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
