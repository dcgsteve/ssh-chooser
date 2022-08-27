package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"fyne.io/systray"
	"github.com/gen2brain/beeep"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {

	systray.SetTitle("SSH Chooser")
	systray.SetIcon(getIcon("assets/ssh.ico"))

	// Add hosts
	for _, host := range getHosts() {
		mItem := systray.AddMenuItem(host, "Connect to "+host)
		go func(host string) {
			<-mItem.ClickedCh
			displayMessage(host)
			c := exec.Command("cmd", "/C", "wt.exe", "ssh", host)
			e := c.Start()
			if e != nil {
				fmt.Printf("Error %v", e)
			}
		}(host)
	}

	// Add exit
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit SSH chooser")
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

func getHosts() []string {

	var hosts []string
	var configFile string = "c:\\users\\me\\.ssh\\config"

	config, err := os.Open(configFile)
	if err != nil {
		os.Exit(1)
	}
	defer config.Close()
	s := bufio.NewScanner(config)
	for s.Scan() {
		if strings.HasPrefix(s.Text(), "Host ") {
			if len(s.Text()) > 5 {
				h := strings.TrimSpace(s.Text()[5:])
				if h != "*" && len(h) > 0 {
					hosts = append(hosts, h)
				}
			}
		}
	}

	return hosts
}

func displayMessage(msg string) {
	err := beeep.Notify("SSH Chooser", msg, "assets/ssh.png")
	if err != nil {
		panic(err)
	}
}
