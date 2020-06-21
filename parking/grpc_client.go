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

func (c *Client) CreateLot(ctx context.Context, maxslotscount uint32) (string, error) {
	r, err := c.service.CreateLot(
		ctx,
		&pb.CreateLotRequest{MaxSlotsCount: maxslotscount},
	)

	if err != nil {
		return r.Status, err
	}

	return r.Status, nil
}

func (c *Client) PostPark(ctx context.Context, carreg string, carcolour string) (*Park, string, error) {
	r, err := c.service.PostPark(
		ctx,
		&pb.PostParkRequest{CarReg: carreg, CarColour: carcolour},
	)

	if err != nil {
		return nil, r.Status, nil
	}

	park := &Park{SlotNum: r.Park.SlotNum, CarReg: r.Park.CarReg, CarColour: r.Park.CarColour}
	return park, r.Status, nil
}

func (c *Client) PostUnpark(ctx context.Context, slotnum uint32) (string, error) {
	r, err := c.service.PostUnpark(
		ctx,
		&pb.PostUnparkRequest{SlotNum: slotnum},
	)

	if err != nil {
		return r.Status, nil
	}

	return r.Status, nil
}

func (c *Client) GetParks(ctx context.Context) ([]Park, string, error) {
	r, err := c.service.GetParks(
		ctx,
		&pb.GetParksRequest{},
	)

	if err != nil {
		return nil, r.Status, nil
	}

	parks := []Park{}
	for _, a := range r.Parks {
		parks = append(parks, Park{
			SlotNum:   a.SlotNum,
			CarReg:    a.CarReg,
			CarColour: a.CarColour,
		})
	}
	return parks, r.Status, nil
}

func (c *Client) GetCarRegsByColour(ctx context.Context, carcolour string) ([]string, string, error) {
	r, err := c.service.GetCarRegsByColour(
		ctx,
		&pb.GetCarRegsByColourRequest{CarColour: carcolour},
	)

	if err != nil {
		return nil, r.Status, nil
	}

	cars := []string{}
	for _, a := range r.Cars {
		cars = append(cars, a)
	}

	return cars, r.Status, nil
}

func (c *Client) GetSlotsByColour(ctx context.Context, carcolour string) ([]uint32, string, error) {
	r, err := c.service.GetSlotsByColour(
		ctx,
		&pb.GetSlotsByColourRequest{CarColour: carcolour},
	)

	if err != nil {
		return nil, r.Status, nil
	}

	slots := []uint32{}
	for _, a := range r.Slots {
		slots = append(slots, a)
	}
	return slots, r.Status, nil
}

func (c *Client) GetSlotByCarReg(ctx context.Context, carreg string) (*Slot, string, error) {
	r, err := c.service.GetSlotByCarReg(
		ctx,
		&pb.GetSlotByCarRegRequest{CarReg: carreg},
	)

	if err != nil {
		return nil, r.Status, nil
	}

	slot := &Slot{SlotNum: r.SlotNum}
	return slot, r.Status, nil
}
