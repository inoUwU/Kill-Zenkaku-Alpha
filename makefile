run:
	go run main.go

debug:
	go build -o debug.exe main.go

build:
	go build -ldflags -H=windowsgui -o KillZenkakuAlpha.exe main.go