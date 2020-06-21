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

## PASS File as an argument:
```
$ ./build/parking_lot bin/fixtures/file_input.txt
```

## REST APIs
```
http://localhost:3569/swagger-parking/
```

## Error Codes & Messages
```
ErrLotSizeZero          Lot size cannot be zero
ErrNoLotFound           No lot available, please create a lot first
ErrParkingFull          Sorry, parking lot is full
ErrInvalidSlot          Slot invalid
ErrParking              Parking slot is empty
ErrInvalidCarNumber     Invalid indian car number plate format
ErrNotFound             Not found
Error                   Unexpected error occured
regexCarNumber          ^[A-Z]{2}-[0-9]{2}-[A-Z]{1,2}-[0-9]{1,4}$
```

TODO: add unit test cases for parking package only (not for commando package)
TODO: verify all make commands and bin functions once again
