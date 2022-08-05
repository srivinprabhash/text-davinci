default: run

run:
	go run main.go

build:
	make build-linux
	make build-windows

build-linux:
	rm -rf dist/linux/textdavinci
	go build -o dist/linux/textdavinci -ldflags "-s -w" main.go

build-windows:
	rm -rf dist/windows/textdavinci-win.exe
	env GOOS=windows GOARCH=amd64 go build -o dist/windows/textdavinci-win.exe -ldflags "-s -w" main.go
	cp dist/windows/textdavinci-win.exe .
