package parking

import (
	"context"
	"errors"
)

var (
	ErrParkingFull = errors.New("parking is full")
)

type Service interface {
	CreateLot(ctx context.Context, maxslotscount uint64) error
	PostPark(ctx context.Context, carreg string, carcolour string) (*Park, error)
	PostUnpark(ctx context.Context, slotnum uint64) error
	GetParks(ctx context.Context) ([]Park, error)
	GetCarRegsByColour(ctx context.Context, carcolour string) ([]Car, error)
	GetSlotsByColour(ctx context.Context, carcolour string) ([]Slot, error)
	GetSlotByCarReg(ctx context.Context, carreg string) (*Slot, error)
}

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

type ParkingService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &ParkingService{r}
}

func (s *ParkingService) CreateLot(ctx context.Context, maxslotscount uint64) error {
	err := s.repository.CreateLot(ctx, maxslotscount)
	if err != nil {
		return err
	}
	return nil
}

func (s *ParkingService) PostPark(ctx context.Context, carreg string, carcolour string) (*Park, error) {
	park, err := s.repository.PostPark(ctx, carreg, carcolour)
	if err != nil {
		return nil, err
	}
	return park, nil
}

func (s *ParkingService) PostUnpark(ctx context.Context, slotnum uint64) error {
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

func (s *ParkingService) GetCarRegsByColour(ctx context.Context, carcolour string) ([]Car, error) {
	cars, err := s.repository.GetCarRegsByColour(ctx, carcolour)
	if err != nil {
		return nil, err
	}
	return cars, nil
}

func (s *ParkingService) GetSlotsByColour(ctx context.Context, carcolour string) ([]Slot, error) {
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
