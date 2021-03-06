//go:generate protoc -I/usr/local/include -I. -I${GOPATH}/src -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:./pb ./parking.proto
//go:generate protoc -I/usr/local/include -I. -I${GOPATH}/src -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --plugin=protoc-gen-grpc-gateway=$GOPATH/bin/protoc-gen-grpc-gateway  --grpc-gateway_out=logtostderr=true:./pb ./parking.proto
//go:generate protoc -I/usr/local/include -I. -I${GOPATH}/src -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --plugin=protoc-gen-swagger=$GOPATH/bin/protoc-gen-swagger  --swagger_out=logtostderr=true:./static ./parking.proto
package parking

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	pb "parking_lot/parking/pb"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service Service
}

func ListenGRPC(s Service, port int) error {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Println(err)
		return err
	}

	serv := grpc.NewServer()
	pb.RegisterParkingServiceServer(serv, &grpcServer{s})
	reflection.Register(serv)
	return serv.Serve(lis)
}

func serveSwaggerJsonHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/parking.swagger.json")
}

func ListenREST(s Service, restport int, grpcport int) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	rmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterParkingServiceHandlerFromEndpoint(ctx, rmux, fmt.Sprintf(":%d", grpcport), opts)
	if err != nil {
		log.Println(err)
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/api/", rmux)
	mux.HandleFunc("/parking.swagger.json", serveSwaggerJsonHandler)
	fs := http.FileServer(http.Dir("./static/swagger"))
	mux.Handle("/swagger-parking/", http.StripPrefix("/swagger-parking", fs))
	return http.ListenAndServe(fmt.Sprintf(":%d", restport), mux)
}

func (s *grpcServer) CreateLot(ctx context.Context, p *pb.CreateLotRequest) (*pb.CreateLotResponse, error) {
	if p.MaxSlotsCount <= 0 {
		return &pb.CreateLotResponse{Status: "ErrLotSizeLTEZero"}, nil
	}

	err := s.service.CreateLot(ctx, p.MaxSlotsCount)

	if err != nil {
		log.Println(err)
		return &pb.CreateLotResponse{Status: "Error"}, nil
	}
	return &pb.CreateLotResponse{Status: "Success"}, nil
}

func (s *grpcServer) PostPark(ctx context.Context, p *pb.PostParkRequest) (*pb.PostParkResponse, error) {
	r, err := s.service.PostPark(ctx, p.CarReg, p.CarColour)

	if err == ErrInvalidCarNumber {
		return &pb.PostParkResponse{Park: &pb.Park{}, Status: "ErrInvalidCarNumber"}, nil
	}

	if err == ErrNoLotFound {
		return &pb.PostParkResponse{Park: &pb.Park{}, Status: "ErrNoLotFound"}, nil
	}

	if err == ErrParkingFull {
		return &pb.PostParkResponse{Park: &pb.Park{}, Status: "ErrParkingFull"}, nil
	}

	if err != nil {
		log.Println(err)
		return &pb.PostParkResponse{Park: &pb.Park{}, Status: "Error"}, nil
	}
	return &pb.PostParkResponse{Park: &pb.Park{SlotNum: r.SlotNum, CarReg: r.CarReg, CarColour: r.CarColour}, Status: "Success"}, nil
}

func (s *grpcServer) PostUnpark(ctx context.Context, p *pb.PostUnparkRequest) (*pb.PostUnparkResponse, error) {
	if p.SlotNum <= 0 {
		return &pb.PostUnparkResponse{Status: "ErrInvalidSlot"}, nil
	}

	err := s.service.PostUnpark(ctx, p.SlotNum)

	if err == ErrInvalidSlot {
		return &pb.PostUnparkResponse{Status: "ErrInvalidSlot"}, nil
	}

	if err == ErrNoLotFound {
		return &pb.PostUnparkResponse{Status: "ErrNoLotFound"}, nil
	}

	if err == ErrParking {
		return &pb.PostUnparkResponse{Status: "ErrParking"}, nil
	}

	if err != nil {
		log.Println(err)
		return &pb.PostUnparkResponse{Status: "Error"}, nil
	}
	return &pb.PostUnparkResponse{Status: "Success"}, nil
}

func (s *grpcServer) GetParks(ctx context.Context, p *pb.GetParksRequest) (*pb.GetParksResponse, error) {
	r, err := s.service.GetParks(ctx)

	if err == ErrNoLotFound {
		return &pb.GetParksResponse{Parks: []*pb.Park{}, Status: "ErrNoLotFound"}, nil
	}

	if err == ErrNotFound {
		return &pb.GetParksResponse{Parks: []*pb.Park{}, Status: "ErrNotFound"}, nil
	}

	if err != nil {
		log.Println(err)
		return &pb.GetParksResponse{Parks: []*pb.Park{}, Status: "Error"}, nil
	}

	parks := []*pb.Park{}
	for _, a := range r {
		parks = append(
			parks,
			&pb.Park{
				SlotNum:   a.SlotNum,
				CarReg:    a.CarReg,
				CarColour: a.CarColour,
			},
		)
	}

	return &pb.GetParksResponse{Parks: parks, Status: "Success"}, nil
}

func (s *grpcServer) GetCarRegsByColour(ctx context.Context, p *pb.GetCarRegsByColourRequest) (*pb.GetCarRegsByColourResponse, error) {
	r, err := s.service.GetCarRegsByColour(ctx, p.CarColour)

	if err == ErrNoLotFound {
		return &pb.GetCarRegsByColourResponse{Cars: []string{}, Status: "ErrNoLotFound"}, nil
	}

	if err == ErrNotFound {
		return &pb.GetCarRegsByColourResponse{Cars: []string{}, Status: "ErrNotFound"}, nil
	}

	if err != nil {
		log.Println(err)
		return &pb.GetCarRegsByColourResponse{Cars: []string{}, Status: "Error"}, nil
	}

	cars := []string{}
	for _, a := range r {
		cars = append(
			cars,
			a)
	}

	return &pb.GetCarRegsByColourResponse{Cars: cars, Status: "Success"}, nil
}

func (s *grpcServer) GetSlotsByColour(ctx context.Context, p *pb.GetSlotsByColourRequest) (*pb.GetSlotsByColourResponse, error) {
	r, err := s.service.GetSlotsByColour(ctx, p.CarColour)

	if err == ErrNoLotFound {
		return &pb.GetSlotsByColourResponse{Slots: []uint32{}, Status: "ErrNoLotFound"}, nil
	}

	if err == ErrNotFound {
		return &pb.GetSlotsByColourResponse{Slots: []uint32{}, Status: "ErrNotFound"}, nil
	}

	if err != nil {
		log.Println(err)
		return &pb.GetSlotsByColourResponse{Slots: []uint32{}, Status: "Error"}, nil
	}

	slots := []uint32{}
	for _, a := range r {
		slots = append(
			slots, a)
	}

	return &pb.GetSlotsByColourResponse{Slots: slots, Status: "Success"}, nil
}

func (s *grpcServer) GetSlotByCarReg(ctx context.Context, p *pb.GetSlotByCarRegRequest) (*pb.GetSlotByCarRegResponse, error) {
	r, err := s.service.GetSlotByCarReg(ctx, p.CarReg)

	var a uint32
	if err == ErrInvalidCarNumber {
		return &pb.GetSlotByCarRegResponse{SlotNum: a, Status: "ErrInvalidCarNumber"}, nil
	}

	if err == ErrNotFound {
		return &pb.GetSlotByCarRegResponse{SlotNum: a, Status: "ErrNotFound"}, nil
	}

	if err != nil {
		log.Println(err)
		return &pb.GetSlotByCarRegResponse{SlotNum: a, Status: "Error"}, nil
	}
	return &pb.GetSlotByCarRegResponse{SlotNum: r.SlotNum, Status: "Success"}, nil
}
