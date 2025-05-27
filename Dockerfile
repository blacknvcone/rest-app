# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install make and git (if needed for dependencies)
RUN apk add --no-cache make git

# Set environment variables for the builder stage
ENV LIBRARY_PATH="${pkgs.tesseract}/lib:${pkgs.leptonica}/lib:${pkgs.opencv4}/lib"
ENV CPATH="${pkgs.tesseract}/include:${pkgs.leptonica}/include:${pkgs.opencv4}/include"
ENV PKG_CONFIG_PATH="${pkgs.tesseract}/lib/pkgconfig:${pkgs.leptonica}/lib/pkgconfig:${pkgs.opencv4}/lib/pkgconfig"
ENV CGO_CPPFLAGS="-I${pkgs.opencv4}/include/opencv4 -I${pkgs.opencv4}/include"
ENV CGO_CXXFLAGS="--std=c++11"
ENV CGO_LDFLAGS="-L${pkgs.opencv4}/lib -lopencv_core -lopencv_highgui -lopencv_imgproc -lopencv_videoio -lopencv_imgcodecs"

# Copy go.mod, go.sum, and Makefile first for better caching
COPY go.mod go.sum Makefile ./

# Install dependencies using Makefile target
RUN make install-deps

# Copy the rest of the source code
COPY . .

# Build the Go app (replace 'main.go' with your entrypoint if different)
RUN go build -o app

# Final stage
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/app .

# Expose port (change if your app uses a different port)
EXPOSE 8080

ENV LIBRARY_PATH="${pkgs.tesseract}/lib:${pkgs.leptonica}/lib:${pkgs.opencv4}/lib"
ENV CPATH="${pkgs.tesseract}/include:${pkgs.leptonica}/include:${pkgs.opencv4}/include"
ENV PKG_CONFIG_PATH="${pkgs.tesseract}/lib/pkgconfig:${pkgs.leptonica}/lib/pkgconfig:${pkgs.opencv4}/lib/pkgconfig"
ENV CGO_CPPFLAGS="-I${pkgs.opencv4}/include/opencv4 -I${pkgs.opencv4}/include"
ENV CGO_CXXFLAGS="--std=c++11"
ENV CGO_LDFLAGS="-L${pkgs.opencv4}/lib -lopencv_core -lopencv_highgui -lopencv_imgproc -lopencv_videoio -lopencv_imgcodecs"

# Run the app
ENTRYPOINT ["./app"]
