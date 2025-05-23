package main

import (
	_ "embed"
	"github.com/getlantern/systray"
	"syscall"
	"time"
)

//go:embed assets/icon.ico
var iconData []byte

var (
	user32              = syscall.NewLazyDLL("user32.dll")
	imm32               = syscall.NewLazyDLL("imm32.dll")
	sendMessage         = user32.NewProc("SendMessageA")
	immGetDefaultIMeWNd = imm32.NewProc("ImmGetDefaultIMEWnd")
	getForegroundWindow = user32.NewProc("GetForegroundWindow")
)

const (
	IME_CONTROL       int    = 0x283
	GETCONVERSIONMODE int    = 1
	GETOPENSTATUS     int    = 5
	SETCONVERSIONMODE int    = 6
	ZENKAKU_ALPHA     int    = 8
	APP_NAME          string = "KillZenkakuAlpha"
	START             string = "Start Kill"
	STOP              string = "Stop kill"
	QUIT              string = "Quit"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(iconData)
	systray.SetTitle(APP_NAME)
	systray.SetTooltip(APP_NAME)
	kill()
	mQuit := systray.AddMenuItem(QUIT, QUIT)
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func onExit() {
	// Cleaning stuff here.
}

func kill() {
	mKill := systray.AddMenuItem(STOP, STOP)

	handleImm := func() {
		hwnd, _, _ := getForegroundWindow.Call()
		imwd, _, _ := immGetDefaultIMeWNd.Call(hwnd)
		imeMode, _, _ := sendMessage.Call(imwd, uintptr(IME_CONTROL), uintptr(GETCONVERSIONMODE), uintptr(0))
		imeState, _, _ := sendMessage.Call(imwd, uintptr(IME_CONTROL), uintptr(GETOPENSTATUS), uintptr(0))
		imeEnabled := imeState != 0

		if imeEnabled && imeMode == uintptr(ZENKAKU_ALPHA) {
			sendMessage.Call(imwd, uintptr(IME_CONTROL), uintptr(SETCONVERSIONMODE), uintptr(0))
		}
	}

	run := true
	toggle := func() {
		if run {
			mKill.SetTitle(START)
			mKill.SetTooltip(START)

			run = false
		} else {
			mKill.SetTitle(STOP)
			mKill.SetTooltip(STOP)
			run = true
		}
	}

	go func() {
		for {
			<-mKill.ClickedCh
			toggle()
		}
	}()

	go func() {
		for {
			if run {
				handleImm()
			}
			time.Sleep(time.Millisecond * 200)
		}
	}()
}
