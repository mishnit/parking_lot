package parking

type Park struct {
	SlotNum   uint64 `json:"SlotNum"`
	CarReg    string `json:"CarReg"`
	CarColour string `json:"CarColour"`
}

type Slot struct {
	SlotNum uint64 `json:"SlotNum"`
}

type Car struct {
	CarReg string `json:"CarReg"`
}
