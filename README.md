## 🚗 PARKING LOT 1.4.2:

A Golang GRPC Microservice with REST and command line interface

## 🦉 REQUIREMENTS:
```
go version go1.13.1 linux/amd64 (for installing dependencies in vendor folder and building commands in build folder)
docker-compose version 1.21.2
Docker version 19.03.2
```

## 🙈 RUN UNIT TESTS + BUILD:
```
$ ./bin/setup
```

## 😼 RUN FUNCTIONAL SPECS CHECK:
```
$ ./bin/run_functional_tests
```

## 🦄 RUN INTERACTIVE SHELL:
```
$ ./bin/parking_lot
```

## 🏆 COMMANDS USAGE:
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

## 🐶 PASS FILE AS AN ARGUMENT:
```
$ ./build/parking_lot bin/fixtures/file_input.txt
```

## 🚀 REST APIS:
```
http://localhost:3569/swagger-parking/
```

## 🐞 ERROR CODES AND MESSAGES:
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
TODO: while sending tarball verify .git is there and (vendor and build are not there)

## 😎 AUTHOR:

- Nitin Mishra
- geekymishnit@gmail.com
