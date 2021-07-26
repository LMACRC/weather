# Needs https://github.com/tsandmann/armv6l-toolchain-mac

export CGO_ENABLED=1
export GOOS=linux
export GOARCH=arm
export GOARM=6
export CC=armv6l-linux-gnueabihf-gcc
export CGO_CFLAGS=-march=armv6
export CGO_CXXFLAGS=-march=armv6
export PATH := /Users/stuartcarnie/Downloads/armv6l-toolchain-mac-master/bin:$(PATH)

arm-weather:
	go build -o bin/rpi/weather ./cmd/weather

arm-weatherctl:
	go build -o bin/rpi/weatherctl ./cmd/weatherctl

env:
	go env