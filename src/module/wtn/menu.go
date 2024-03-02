package wtn

import (
	console_logger "WTN-Sniper/src/packages/logger"
	"fmt"
	"os"
	"time"
)

var lm = console_logger.New("WTN SNIPER", true)

func actions() string {
	var selection int
	options := 4
	lm.General(">>> Please select an action")
	lm.General(">>> 1 - Export listings")
	lm.General(">>> 2 - Extend listings")
	lm.General(">>> 3 - Download labels")
	lm.General(">>> 4 - Start sniper")
	lm.General(">>> 0 - Exit")

	fmt.Printf(console_logger.TextColors.White+"\r[%s] [%s] >>> "+console_logger.TextColors.Reset, "WTN SNIPER", lm.GetTime())

	if _, err := fmt.Scan(&selection); err != nil || selection <= -1 || selection > options {
		lm.Error("Wrong input! Please retry!")
		time.Sleep(3 * time.Second)
	}

	switch selection {
	case 0:
		os.Exit(0)
		return ""
	case 1:
		return "Export"
	case 2:
		return "Extend"
	case 3:
		return "Download"
	case 4:
		return "Sniper"
	default:
		return ""
	}
}
