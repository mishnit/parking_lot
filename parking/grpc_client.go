package parking

import (
	"context"
	"log"
	pb "parking_lot/parking/pb"

	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.ParkingServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		log.Println(err)
		return nil, err
	}
	c := pb.NewParkingServiceClient(conn)
	return &Client{conn, c}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) CreateLot(ctx context.Context, maxslotscount uint64) (string, error) {
	r, err := c.service.CreateLot(
		ctx,
		&pb.CreateLotRequest{MaxSlotsCount: maxslotscount},
	)
	if err != nil {
		log.Println(err)
		return "Error", err
	}
	return r.Status, nil
}

func (c *Client) PostPark(ctx context.Context, carreg string, carcolour string) (*Park, error) {
	r, err := c.service.PostPark(
		ctx,
		&pb.PostParkRequest{CarReg: carreg, CarColour: carcolour},
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	park := &Park{SlotNum: r.Park.SlotNum, CarReg: r.Park.CarReg, CarColour: r.Park.CarColour}
	return park, nil
}

func (c *Client) PostUnpark(ctx context.Context, slotnum uint64) (string, error) {
	r, err := c.service.PostUnpark(
		ctx,
		&pb.PostUnparkRequest{SlotNum: slotnum},
	)
	if err != nil {
		log.Println(err)
		return "Error", err
	}
	return r.Status, nil
}

func (c *Client) GetParks(ctx context.Context) ([]Park, error) {
	r, err := c.service.GetParks(
		ctx,
		&pb.GetParksRequest{},
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	parks := []Park{}
	for _, a := range r.Parks {
		parks = append(parks, Park{
			SlotNum:   a.SlotNum,
			CarReg:    a.CarReg,
			CarColour: a.CarColour,
		})
	}
	return parks, nil
}

func (c *Client) GetCarRegsByColour(ctx context.Context, carcolour string) ([]Car, error) {
	r, err := c.service.GetCarRegsByColour(
		ctx,
		&pb.GetCarRegsByColourRequest{CarColour: carcolour},
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	cars := []Car{}
	for _, a := range r.Cars {
		cars = append(cars, Car{
			CarReg: a.CarReg,
		})
	}
	return cars, nil
}

func (c *Client) GetSlotsByColour(ctx context.Context, carcolour string) ([]Slot, error) {
	r, err := c.service.GetSlotsByColour(
		ctx,
		&pb.GetSlotsByColourRequest{CarColour: carcolour},
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	slots := []Slot{}
	for _, a := range r.Slots {
		slots = append(slots, Slot{
			SlotNum: a.SlotNum,
		})
	}
	return slots, nil
}

func (c *Client) GetSlotByCarReg(ctx context.Context, carreg string) (*Slot, error) {
	r, err := c.service.GetSlotByCarReg(
		ctx,
		&pb.GetSlotByCarRegRequest{CarReg: carreg},
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	slot := &Slot{SlotNum: r.SlotNum}
	return slot, nil
}
