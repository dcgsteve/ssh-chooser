package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"fyne.io/systray"
	"github.com/gen2brain/beeep"
)

//go:embed winres/chooser.ico
var ico []byte

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {

	systray.SetTitle("SSH Chooser")
	systray.SetTemplateIcon(ico, ico)

	// Add hosts
	for _, host := range getHosts() {
		mItem := systray.AddMenuItem(host, "Connect to "+host)
		go func(host string) {
			<-mItem.ClickedCh
			e := triggerTerminal(host)
			if e != nil {
				displayMessage(e.Error())
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

// func getIcon(s string) []byte {
// 	b, err := os.ReadFile(s)
// 	if err != nil {
// 		fmt.Print(err)
// 	}
// 	return b
// }

func getHosts() []string {

	var hosts []string

	configFile := filepath.Join(fmt.Sprintf("%s%s", os.Getenv("USERPROFILE"), "\\.ssh\\config"))
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

	sort.Strings(hosts)

	return hosts
}

func displayMessage(msg string) {
	err := beeep.Notify("SSH Chooser", msg, "assets/ssh.png")
	if err != nil {
		panic(err)
	}
}

func triggerTerminal(host string) error {
	tPath := filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft\\WindowsApps\\wt.exe")
	procAttr := new(os.ProcAttr)
	procAttr.Files = []*os.File{nil, nil, nil}

	args := []string{"ssh", host}
	_, e := os.StartProcess(tPath, append([]string{tPath}, args...), procAttr)

	return e
}
