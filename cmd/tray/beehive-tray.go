package main

import (
	"os"
	"os/exec"

	"log"

	"github.com/getlantern/systray"
)

var beehiveCmd *exec.Cmd
var webViewCmd *exec.Cmd

func main() {
	systray.Run(do, onExit)
}

func do() {
	systray.SetIcon(Icon)
	systray.SetTooltip("Beehive")
	open := systray.AddMenuItem("Open Beehive", "Open Beehive")
	go func() {
		for {
			<-open.ClickedCh
			webViewCmd = exec.Command("beehive-web")
			webViewCmd.Run()
		}
	}()

	quit := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-quit.ClickedCh
		systray.Quit()
	}()
	for {
		beehiveCmd = exec.Command("beehive")
		beehiveCmd.Run()
		log.Fatal("Beehive process died")
	}
}

func onExit() {
	if err := beehiveCmd.Process.Signal(os.Interrupt); err != nil {
		log.Fatal("Failed to kill process: ", err)
	}
}
