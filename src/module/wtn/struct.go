package wtn

import (
	"WTN-Sniper/src/packages/file_manager"
	"WTN-Sniper/src/packages/notifier"
	"WTN-Sniper/src/packages/session"
)

type WTNSession struct {
	*session.Session
	notifierSession *notifier.Notifier
	settingsData    file_manager.SettingsData
	loginData       LoginSession
}
