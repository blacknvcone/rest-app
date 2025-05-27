# rest-app
rest-app

## Installation
### Go-migrate CLI
```sh
#mac
$ brew install golang-migrate

#linux
$ curl -L https://github.com/golang-migrate/migrate/releases/download/$version/migrate.$platform-amd64.tar.gz | tar xvz

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

## Accessing Swagger
```
localhost:8080/swagger/index.html
```