package main

import (
	"fmt"
	"github.com/getlantern/systray"
	"golang.org/x/sys/windows"
	"os"
	"unsafe"
)

// https://dev.to/osuka42/building-a-simple-system-tray-app-with-go-899

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

// kill zenkaku
func killing() {
	mKill := systray.AddMenuItem("Start kill", "Start kill")
	user32 := windows.NewLazyDLL("user32.dll")
	getWindow := user32.NewProc("GetWindowThreadProcessId")
	// getContext := windows.NewLazyDLL("imm32.dll").NewProc("ImmGetContext")
	// getTickCount := windows.NewLazyDLL("imm32.dll").NewProc("ImmGetConversionStatus")

	// 起動時間を取得
	// r, b, _ := getTickCount.Call()
	a, _, _ := getWindow.Call(uintptr(0))
	fmt.Println(a)

	go func() {

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
		toggle()

		for {
			<-mKill.ClickedCh
			toggle()
			fmt.Println("Hello")
		}
	}()
}

type HWND uintptr

func (t HWND) uintptr() uintptr {
	return uintptr(t)
}

type LPCTSTR string

func (t LPCTSTR) uintptr() uintptr {
	return uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(string(t))))
}

type UINT uint32
type MBType UINT

func (t MBType) uintptr() uintptr {
	return uintptr(t)
}

const (
	MBTypeOK MBType = 0x00000000
	// 以下略
)
