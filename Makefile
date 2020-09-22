test:
	cd index && go test
	cd node  && go test
	cd prop  && go test
	cd sgp4  && go test
	cd tle   && go test
	cd cmd/spibatch && go test
	cd cmd/spipipe  && go test
	cd cmd/spitool  && go test

all: test

install: test
	cd cmd/spibatch && go install
	cd cmd/spipipe  && go install
	cd cmd/spitool  && go install

