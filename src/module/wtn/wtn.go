package wtn

import (
	"WTN-Sniper/src/packages/file_manager"
	console_logger "WTN-Sniper/src/packages/logger"
	"WTN-Sniper/src/packages/notifier"
	"WTN-Sniper/src/packages/proxy_manager"
	"WTN-Sniper/src/packages/session"
	"fmt"
	"os"
	"time"

	"github.com/bogdanfinn/tls-client/profiles"
)

var l = console_logger.New("WTN Sniper", true)

func InitWtn(proxyFile string, settingsData file_manager.SettingsData) {
	proxySession := proxy_manager.New(proxyFile)
	s, err := session.New(30, profiles.Chrome_120, proxySession)
	if err != nil {
		l.Error(err.Error())
		time.Sleep(20 * time.Second)
		os.Exit(1)
	}

	notifierSession := notifier.New(20)
	WTNSession := &WTNSession{
		Session:         s,
		notifierSession: notifierSession,
		settingsData:    settingsData,
	}
	go notifierSession.StartQueue(2 * time.Second)
	go WTNSession.updateSettings()

	// WTN module has been initialised
	loginStep := true
	for loginStep {
		loginSession, err := WTNSession.Login()
		if err != nil {
			l.Error(err.Error())
			time.Sleep(time.Duration(WTNSession.settingsData.Delay) * time.Second)
		}
		WTNSession.loginData = loginSession
		loginStep = false
	}

	// Save login data
	err = WTNSession.saveLoginData()
	if err != nil {
		l.Error("An error occured while saving the login data! Skipping...")
	}

	l.Success("Logged in!")

	// Now we are logged in, we can proceed

	for {
		l.Clear()
		console_logger.PrintLogo()
		fmt.Println("")

		l.Notification("Logged into the account. WTN Sniper ready")
		l.General(fmt.Sprintf("Using %s as proxy file", proxyFile))
		fmt.Println("")

		choice := actions()
		switch choice {
		case "Export":
			err = WTNSession.exportListings()
			if err != nil {
				l.Error(err.Error())
			} else {
				l.Success("Listings exported in ./Settings/listings.csv")
			}
		case "Extend":
			err = WTNSession.extendListings()
			if err != nil {
				l.Error(err.Error())
			} else {
				l.Success("Extend time task completed!")
			}
		case "Download":
			fmt.Println("Download labels")
			err = WTNSession.label()
			if err != nil {
				l.Error(err.Error())
			} else {
				l.Success("Labels downloaded")
			}
		case "Sniper":
			// Wrap this into another function to do not have anything in the
			err = WTNSession.sniper()
			if err != nil {
				l.Error(err.Error())
			} else {
				l.Success("Sniper task completed")
			}
		default:
			l.Error("Invalid selection, please retry!")
			time.Sleep(5 * time.Second)
		}
		time.Sleep(15 * time.Second)
	}
}

// We use this function to update the delay settings
func (s *WTNSession) updateSettings() {
	for {
		tempData, err := file_manager.ReadSettings()
		if err != nil {
			l.Error(err.Error())
			time.Sleep(20 * time.Second)
			os.Exit(1)
		}

		if s.settingsData.Delay != tempData.Delay {
			l.Notification(fmt.Sprintf("Delay changed to: %v", tempData.Delay))
		}

		s.settingsData = tempData
		time.Sleep(2 * time.Second)
	}
}
