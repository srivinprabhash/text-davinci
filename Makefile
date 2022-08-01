default: run

run:
	go run main.go

build:
	rm -rf main
	go build -o textdavinci -ldflags "-s -w" main.go