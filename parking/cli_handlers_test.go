package parking

import "testing"

// positive test cases

func TestCreateParkingLot(t *testing.T) {
	total := CreateParkingLot(9)
	if total != "Created a parking lot with 9 slots" {
		t.Errorf("CreateParkingLot was incorrect, got: %s, want: %s.", total, "Created a parking lot with 9 slots")
	}
}

func TestDoPark(t *testing.T) {
	CreateParkingLot(3)
	DoPark("KA-09-HH-3333", "White")
	park := DoPark("KA-09-HH-3334", "White")
	if park != "Allocated slot number: 2" {
		t.Errorf("Park was incorrect, got: %s, want: %s.", park, "Allocated slot number: 2")
	}
}

func TestLeave(t *testing.T) {
	leave := Leave(1)
	if leave != "Slot number 1 is free" {
		t.Errorf("Leave was incorrect, got: %s, want: %s.", leave, "Slot number 1 is free")
	}
}

func TestStatus(t *testing.T) {
	CreateParkingLot(5)
	DoPark("KA-01-HH-1234", "White")
	DoPark("KA-01-HH-1235", "Black")
	status := Status()
	if len(status) != 2 {
		t.Errorf("Status was incorrect, got: %d, want: %d.", len(status), 2)
	}
}

func TestRegistrationNumbersForCarsWithColour(t *testing.T) {
	carPlateNumber := RegistrationNumbersForCarsWithColour("White")
	if carPlateNumber != "KA-01-HH-1234" {
		t.Errorf("RegistrationNumbersForCarsWithColour was incorrect, got: %s, want: %s.", carPlateNumber, "KA-01-HH-1234")
	}
}

func TestSlotNumbersForCarsWithColour(t *testing.T) {
	slotNumber := SlotNumbersForCarsWithColour("White")
	if slotNumber != "1" {
		t.Errorf("SlotNumbersForCarsWithColour was incorrect, got: %s, want: %s.", slotNumber, "1")
	}
}

func TestSlotNumberForRegistrationNumber(t *testing.T) {
	slotNumber := SlotNumberForRegistrationNumber("KA-01-HH-1234")
	if slotNumber != "1" {
		t.Errorf("Park was incorrect, got: %s, want: %s.", slotNumber, "1")
	}
}

// TODO: negative test cases
