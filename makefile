run:
	go run main.go

debug:
	go build -o debug.exe main.go

clean:
	if exist icon.syso del icon.syso
	if exist .\build\KillZenkakuAlpha.exe del .\build\KillZenkakuAlpha.exe

make-icon:
	windres --output-format=coff -o icon.syso icon.rc

build: clean make-icon
	go build -ldflags -H=windowsgui -o ./build/KillZenkakuAlpha.exe