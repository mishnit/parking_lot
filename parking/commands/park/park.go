package main

import (
	"fmt"
	"os"

	"parking_lot/parking"
)

func main() {
	args := os.Args
	if len(args) == 3 {
		cpn := os.Args[1]
		cc := os.Args[2]
		park := parking.DoPark(cpn, cc)
		fmt.Println(park)
	}
}