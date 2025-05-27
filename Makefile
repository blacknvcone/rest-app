include .env
.PHONY: test clean build

install-deps:
    brew update
    brew install tesseract leptonica opencv pkg-config libglvnd
    echo '#!/usr/bin/env bash' > env.sh
    echo 'export LIBRARY_PATH="${pkgs.tesseract}/lib:${pkgs.leptonica}/lib:${pkgs.opencv4}/lib"' >> env.sh
    echo 'export CPATH="${pkgs.tesseract}/include:${pkgs.leptonica}/include:${pkgs.opencv4}/include"' >> env.sh
    echo 'export PKG_CONFIG_PATH="${pkgs.tesseract}/lib/pkgconfig:${pkgs.leptonica}/lib/pkgconfig:${pkgs.opencv4}/lib/pkgconfig"' >> env.sh
    echo 'export CGO_CPPFLAGS="-I${pkgs.opencv4}/include/opencv4 -I${pkgs.opencv4}/include"' >> env.sh
    echo 'export CGO_CXXFLAGS="--std=c++11"' >> env.sh
    echo 'export CGO_LDFLAGS="-L${pkgs.opencv4}/lib -lopencv_core -lopencv_highgui -lopencv_imgproc -lopencv_videoio -lopencv_imgcodecs"' >> env.sh
    echo "Run 'source env.sh' to set environment variables."

build:
	go mod download
	go build -o setup main.go

test:
	go get github.com/newm4n/goornogo
	go test ./... -v -covermode=count -coverprofile=coverage.out
	goornogo -c 20 -i coverage.out

clean: 
	go clean

run:
	go mod download
	go run main.go

migrateup:
	migrate -path migrations -database "${DB_DSN}" -verbose up

migratedown:
	migrate -path migrations -database "${DB_DSN}" -verbose down