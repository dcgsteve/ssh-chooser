package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Menu struct {
	Label string
	Items []*fyne.MenuItem
}

var a fyne.App

func main() {

	a = app.New()
	w := a.NewWindow("Click a Host!")

	// List box with host details
	hosts := getHosts()

	lstHosts := widget.NewList(
		func() int {
			return len(hosts)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("example-length")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(hosts[i])
		})

	lstHosts.OnSelected = func(id widget.ListItemID) {

		// c := cmd.NewCmd("bash", "-c", "ssh", hosts[id])
		c := exec.Command("bash", "-c", "echo", "hello")
		err := c.Start()
		if err != nil {
			showError(err)
		}
		lstHosts.UnselectAll()
	}

	// Build window
	w.SetContent(
		container.NewMax(
			lstHosts,
		),
	)
	w.Resize(fyne.NewSize(300, 500))

	// Start GUI
	w.ShowAndRun()

}

func getHosts() []string {

	var hosts []string
	var configFile string = "/home/steve/.ssh/config"

	config, err := os.Open(configFile)
	if err != nil {
		fmt.Printf("Could not open config file [%s]: %s", configFile, err)
	}
	defer config.Close()
	s := bufio.NewScanner(config)
	for s.Scan() {
		if strings.HasPrefix(s.Text(), "Host ") {
			if len(s.Text()) > 5 {
				h := strings.TrimSpace(s.Text()[5:])
				if h != "*" && len(h) > 0 {
					hosts = append(hosts, h)
					fmt.Printf("Found host: %s\n", h)
				}
			}
		}
	}

	return hosts
}

func showError(e error) {

	w := a.NewWindow("Error")
	lbl := widget.NewLabel(fmt.Sprintf("An error has occurred - see below for details:\n\n%s", e))
	// Build window
	w.SetContent(
		container.NewMax(
			lbl,
		),
	)
	w.Resize(fyne.NewSize(200, 100))

	// Start GUI
	w.Show()

}
