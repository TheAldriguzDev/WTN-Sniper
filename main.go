package main

import (
	wtn "WTN-Sniper/src/module/WTN"
	"WTN-Sniper/src/packages/file_manager"
	console_logger "WTN-Sniper/src/packages/logger"
	"WTN-Sniper/src/packages/menu"
	"fmt"
	"os"
	"time"
)

var l = console_logger.New("Dashboard", true)

func main() {
	if !file_manager.CheckSetup() {
		l.Error("Error while checking the bot setup!")
		time.Sleep(20 * time.Second)
		os.Exit(1)
	}

	settingsData, err := file_manager.ReadSettings()
	if err != nil {
		l.Error(err.Error())
		time.Sleep(20 * time.Second)
		os.Exit(1)
	}

	fmt.Println("")
	// We can now boot up the bot
	l.Clear()
	console_logger.PrintLogo()

	for {
		choice := menu.DashboardMenu("YOUR_NAME")
		switch choice {
		case "WTN":
			proxyFile, err := file_manager.ReadProxies("Proxies")
			if err != nil {
				l.Error(err.Error())
				time.Sleep(20 * time.Second)
				os.Exit(1)
			}
			wtn.InitWtn(proxyFile, settingsData)
			time.Sleep(5 * time.Second)
		case "TEST":
			testResults := wtn.TestWebhook(settingsData)
			if testResults {
				l.General("Webhooks tested!")
			} else {
				l.Warning("An error occured while testing the webhooks!")
				time.Sleep(5 * time.Second)
			}
			time.Sleep(5 * time.Second)
		default:
			// Uscirà solo in caso di errore nell'input
			fmt.Println("Invalid selection! Exiting...")
			os.Exit(1)
		}
	}
}
