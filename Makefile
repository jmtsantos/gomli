# This how we want to name the binary output
BINARY=./build/gomli

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-X github.com/jmtsantos/gomli/cmd.Version=1.0.0"

# Builds the project
build:
	go build ${LDFLAGS} -o ${BINARY} ./gomli/main.go

# Installs our project: copies binaries
install:
	go install ${LDFLAGS} github.com/jmtsantos/gomli/gomli

# Cleans our project: deletes binaries
clean:
	rm -rf ./build

.PHONY: clean install