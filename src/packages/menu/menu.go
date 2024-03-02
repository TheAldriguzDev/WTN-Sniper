package menu

import (
	console_logger "WTN-Sniper/src/packages/logger"
	"fmt"
	"os"
	"time"
)

var ld = console_logger.New("Dashboard", true)

func DashboardMenu(userName string) string {
	ld.Clear()
	console_logger.PrintLogo()
	fmt.Println("")
	ld.Notification("Welcome back, " + userName + "!")
	fmt.Println("")

	var selection int
	options := 2
	ld.General(">>> Please select an action")
	ld.General(">>> 1 - WTN Sniper")
	ld.General(">>> 2 - Test webhook")
	ld.General(">>> 0 - Exit")

	fmt.Printf(console_logger.TextColors.White+"\r[%s] [%s] >>> "+console_logger.TextColors.Reset, "Dashboard", ld.GetTime())

	if _, err := fmt.Scan(&selection); err != nil || selection <= -1 || selection > options {
		ld.Error("Wrong input! Please retry!")
		time.Sleep(3 * time.Second)
	}

	switch selection {
	case 0:
		os.Exit(0)
		return ""
	case 1:
		return "WTN"
	case 2:
		return "TEST"
	default:
		os.Exit(0)
		return ""
	}
}
