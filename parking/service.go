package parking

import (
	"context"
	"errors"
	"regexp"
)

var (
	ErrParkingFull      = errors.New("Lot is Full!")
	ErrParkingEmpty     = errors.New("Lot is Empty!")
	ErrParking          = errors.New("Parking slot is Empty!")
	regexCarNumber      = regexp.MustCompile(`^[A-Z]{2}-[0-9]{2}-[A-Z]{2}-[0-9]{4}$`)
	ErrInvalidCarNumber = errors.New("Invalid Car Number Plate!")
)

type Service interface {
	CreateLot(ctx context.Context, maxslotscount uint32) error
	PostPark(ctx context.Context, carreg string, carcolour string) (*Park, error)
	PostUnpark(ctx context.Context, slotnum uint32) error
	GetParks(ctx context.Context) ([]Park, error)
	GetCarRegsByColour(ctx context.Context, carcolour string) ([]string, error)
	GetSlotsByColour(ctx context.Context, carcolour string) ([]uint64, error)
	GetSlotByCarReg(ctx context.Context, carreg string) (*Slot, error)
}

type Park struct {
	SlotNum   uint32 `json:"SlotNum"`
	CarReg    string `json:"CarReg"`
	CarColour string `json:"CarColour"`
}

type Slot struct {
	SlotNum uint32 `json:"SlotNum"`
}

type Car struct {
	CarReg string `json:"CarReg"`
}

type ParkingService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &ParkingService{r}
}

func (s *ParkingService) CreateLot(ctx context.Context, maxslotscount uint32) error {
	err := s.repository.CreateLot(ctx, maxslotscount)
	if err != nil {
		return err
	}
	return nil
}

func (s *ParkingService) PostPark(ctx context.Context, carreg string, carcolour string) (*Park, error) {
	if !regexCarNumber.MatchString(carreg) {
		return nil, ErrInvalidCarNumber
	}
	park, err := s.repository.PostPark(ctx, carreg, carcolour)
	if err != nil {
		return nil, err
	}
	return park, nil
}

func (s *ParkingService) PostUnpark(ctx context.Context, slotnum uint32) error {
	err := s.repository.PostUnpark(ctx, slotnum)
	if err != nil {
		return err
	}
	return nil
}

func (s *ParkingService) GetParks(ctx context.Context) ([]Park, error) {
	parks, err := s.repository.GetParks(ctx)
	if err != nil {
		return nil, err
	}
	return parks, nil
}

func (s *ParkingService) GetCarRegsByColour(ctx context.Context, carcolour string) ([]string, error) {
	cars, err := s.repository.GetCarRegsByColour(ctx, carcolour)
	if err != nil {
		return nil, err
	}
	return cars, nil
}

func (s *ParkingService) GetSlotsByColour(ctx context.Context, carcolour string) ([]uint64, error) {
	slots, err := s.repository.GetSlotsByColour(ctx, carcolour)
	if err != nil {
		return nil, err
	}
	return slots, nil
}

func (s *ParkingService) GetSlotByCarReg(ctx context.Context, carreg string) (*Slot, error) {
	slot, err := s.repository.GetSlotByCarReg(ctx, carreg)
	if err != nil {
		return nil, err
	}
	return slot, nil
}
