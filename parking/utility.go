package parking

import (
	"time"
)

type Func func(int) error

func find(slice []uint32, val uint32) (int, bool) {
	for idx, item := range slice {
		if uint32(item) == val {
			return idx, true
		}
	}
	return -1, false
}

func ForeverSleep(d time.Duration, f Func) {
	for i := 0; ; i++ {
		err := f(i)
		if err == nil {
			return
		}
		time.Sleep(d)
	}
}

func nextslot(MaxSlotsCount uint32, UsedSlots []uint32) uint32 {
	var i uint32
	i = 1
	for i <= MaxSlotsCount {
		_, found := find(UsedSlots, i)
		if !found {
			break
		}
		i = i + 1
	}
	return i
}
