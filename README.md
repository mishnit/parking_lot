Parking lot - A GRPC Microservice with REST and Command line interface (including test cases)

DEVELOPMENT REQUIREMENTS:
```
go version go1.13.1 linux/amd64
docker-compose version 1.21.2
Docker version 19.03.2
```

TEST (docker):
```
$ ./bin/setup
```

RUN (docker):
```
$ ./bin/parking_lot
```

REST/CURL Endpoints
```
http://localhost:8000/swagger-parking/
```

TODO: Two CLI interfaces for grpc (first for shell launch and second with file input as argument)
TODO: create test files with seed data (positive and Negative tests remaining)
TODO: update parking_lot file for cmd application


send tarball zip and no code public being made

Check this : https://github.com/george-e-shaw-iv/integration-tests-example/blob/master/cmd/listd/tests/item_test.go
