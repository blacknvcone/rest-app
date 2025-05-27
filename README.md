# rest-app
rest-app

## Purpose
This application processes uploaded images and generates JSON data based on a defined structure using OCR and image processing libraries.

## API Documentation

### OCR Receipt Endpoint

**POST** `http://localhost:8089/v1/public-api/ocr/receipt`

- **Description:** Upload an image of a receipt to extract structured JSON data.
- **Request:** `multipart/form-data` with an `image` file field.
- **Response:** JSON object containing extracted data.

#### Example Request (using curl)
```sh
curl -X POST http://localhost:8089/v1/public-api/ocr/receipt \
  -F "image=@/path/to/your/receipt.jpg"
```

#### Example Response
```json
{
  "status": "success",
  "data": {
    // ...extracted fields based on defined structure...
  }
}
```

## Installation

### Required Local Dependencies for OCR
This app requires the following libraries to be installed on your local machine for OCR processing:
- [Tesseract](https://github.com/tesseract-ocr/tesseract)
- [Leptonica](http://www.leptonica.org/)
- [OpenCV4](https://opencv.org/)
- [libglvnd](https://github.com/NVIDIA/libglvnd)

#### Install on macOS (using Homebrew)
```sh
brew install tesseract leptonica opencv libglvnd
```

#### Install on Ubuntu/Debian
```sh
sudo apt-get update
sudo apt-get install tesseract-ocr libleptonica-dev libopencv-dev libglvnd-dev
```

### Go-migrate CLI
```sh
#mac
$ brew install golang-migrate

#linux
$ curl -L https://github.com/golang-migrate/migrate/releases/download/$version/migrate.$platform-amd64.tar.gz | tar xvz
```

## How To Run
#### Using Makefile
```sh
#already install swag and air
$ make run 
```

## Technologies
- [Golang](https://go.dev/)
- [Gorm](https://gorm.io/index.html)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [Swaggo](https://github.com/swaggo/swag)
- PostgreSQL
- Tesseract
- Leptonica
- OpenCV4
- libglvnd

## Accessing Swagger
```
localhost:8080/swagger/index.html
```