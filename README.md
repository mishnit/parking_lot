# Parking lot - Golang GRPC Microservice with REST and Command line interface

## REQUIREMENTS:
```
go version go1.13.1 linux/amd64
docker-compose version 1.21.2
Docker version 19.03.2
```

## RUN UNIT TESTS:
```
$ ./bin/setup
```

## TEST FUNCTIONAL SPECS:
```
$ ./bin/run_functional_tests
```

## RUN INTERACTIVE SHELL:
```
$ ./bin/parking_lot
```

## PASS File as an argument:
```
$ ./build/parking_lot bin/fixtures/file_input.txt
```

## REST APIs
```
http://localhost:3569/swagger-parking/
```

## Commands Usage
```
create_parking_lot <max_slots_num>                        Create Parking lot of size n
park <car_reg_number> <car_colour>                        Park car in the first available slot from entry gate
leave <slot_num>                                          Unpark car from given slot num
status                                                    List all occupied slots along with car details
registration_numbers_for_cars_with_colour <car_colour>    List car registration numbers having specific colour
slot_numbers_for_cars_with_colour <car_colour>            List slot numbers parked with car having specific colour
slot_number_for_registration_number <car_reg_number>      Display slot number for given car registration number
exit                                                      Exit from shell
```

## Error Codes & Formats (check all muct be part of commands as well as rest)
```
ErrNoParkingLotCreated  No parking lot available. Create parking lot first
ErrLotSizeZero          Lot size cannot be zero!
ErrParkingFull          Sorry, parking lot is full
ErrParkingEmpty         Lot is Empty!
ErrInvalidSlot          Slot Invalid!
ErrParking              Parking slot is Empty!
ErrInvalidCarNumber     Invalid Indian Car Number Plate Format!
regexCarNumber          ^[A-Z]{2}-[0-9]{2}-[A-Z]{1,2}-[0-9]{4}$
```

TODO: update cli_handlers_test.go for positive and error code checks unit test cases
TODO: send tarball zip and no code public being made
