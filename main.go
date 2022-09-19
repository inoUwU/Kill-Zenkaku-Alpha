package main

import (
	"fmt"
	"github.com/getlantern/systray"
	"os"
	"syscall"
	"time"
)

const (
	WM_IME_CONTROL        int = 0x283
	IMC_GETCONVERSIONMODE int = 1
	IMC_GETOPENSTATUS     int = 5
	IMC_SETCONVERSIONMODE int = 6
	IME_ZENKAKU_Alpha     int = 8
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(getIcon("assets/icon128.ico"))
	systray.SetTitle("KillZenkakuAlpha")
	systray.SetTooltip("KillZenkakuAlpha")
	killing()
	mQuit := systray.AddMenuItem("Quit", "Quit app")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func onExit() {
	// Cleaning stuff here.
}

func getIcon(s string) []byte {
	b, err := os.ReadFile(s)
	if err != nil {
		fmt.Print(err)
	}
	return b
}

func killing() {
	mKill := systray.AddMenuItem("Stop Kill", "Stop Kill")

	user32, err := syscall.LoadDLL("user32.dll")
	if err != nil {
		panic(err)
	}
	defer user32.Release()

	imm32, err := syscall.LoadDLL("imm32.dll")
	if err != nil {
		panic(err)
	}
	defer imm32.Release()

	procGetForegroundWindow, err := user32.FindProc("GetForegroundWindow")
	if err != nil {
		panic(err)
	}

	hwnd, _, _ := procGetForegroundWindow.Call()

	immGetDefaultIMeWNd, err := imm32.FindProc("ImmGetDefaultIMEWnd")
	if err != nil {
		panic(err)
	}

	imwd, _, _ := immGetDefaultIMeWNd.Call(hwnd)

	sendMessage, err := user32.FindProc("SendMessageA")
	if err != nil {
		panic(err)
	}

	handleImm := func() {
		imeConvMode, _, _ := sendMessage.Call(imwd, uintptr(WM_IME_CONTROL), uintptr(IMC_GETCONVERSIONMODE), uintptr(0))
		imeState, _, _ := sendMessage.Call(imwd, uintptr(WM_IME_CONTROL), uintptr(IMC_GETOPENSTATUS), uintptr(0))

		imeEnabled := imeState != 0

		if imeEnabled && imeConvMode == uintptr(IME_ZENKAKU_Alpha) {
			sendMessage.Call(imwd, uintptr(WM_IME_CONTROL), uintptr(IMC_SETCONVERSIONMODE), uintptr(0))
		}
	}

	run := true
	toggle := func() {
		if run {
			mKill.SetTitle("Stop Kill")
			mKill.SetTooltip("Stop kill")

			run = false
		} else {
			mKill.SetTitle("Start Kill")
			mKill.SetTooltip("Start kill")
			run = true
		}
	}

	go func() {
		for {
			<-mKill.ClickedCh
			toggle()
		}
	}()

	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()
	go func() {
		for {
			if run {
				handleImm()
			}
		}
	}()
}
