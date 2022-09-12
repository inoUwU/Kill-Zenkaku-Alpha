package main

import (
	"fmt"
	"io/ioutil"

	"github.com/getlantern/systray"
)

// https://dev.to/osuka42/building-a-simple-system-tray-app-with-go-899
func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(getIcon("assets/icon128.ico"))
	systray.SetTitle("KillZenkakuAlpha")
	systray.SetTooltip("Look at me, I'm a tooltip!")
}

func onExit() {
	// Cleaning stuff here.
}

func getIcon(s string) []byte {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		fmt.Print(err)
	}
	return b
}
