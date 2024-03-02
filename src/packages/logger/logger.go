package console_logger

import (
	"fmt"
	"time"
)

var TextColors = txt_color{
	Reset:   "\x1b[0m",
	Black:   "\x1b[30m",
	Red:     "\x1b[31m",
	Green:   "\x1b[32m",
	Yellow:  "\x1b[33m",
	Blue:    "\x1b[34m",
	Magenta: "\x1b[35m",
	Cyan:    "\x1b[36m",
	White:   "\x1b[37m",
}

func New(prefix string, enable_time_stamp bool) *logger {
	return &logger{
		prefix:            prefix,
		enable_time_stamp: enable_time_stamp,
	}
}

func PrintLogo() {
	fmt.Print(TextColors.Red + `
	██╗    ██╗████████╗███╗   ██╗    ███████╗███╗   ██╗██╗██████╗ ███████╗██████╗ 
	██║    ██║╚══██╔══╝████╗  ██║    ██╔════╝████╗  ██║██║██╔══██╗██╔════╝██╔══██╗
	██║ █╗ ██║   ██║   ██╔██╗ ██║    ███████╗██╔██╗ ██║██║██████╔╝█████╗  ██████╔╝
	██║███╗██║   ██║   ██║╚██╗██║    ╚════██║██║╚██╗██║██║██╔═══╝ ██╔══╝  ██╔══██╗
	╚███╔███╔╝   ██║   ██║ ╚████║    ███████║██║ ╚████║██║██║     ███████╗██║  ██║
	 ╚══╝╚══╝    ╚═╝   ╚═╝  ╚═══╝    ╚══════╝╚═╝  ╚═══╝╚═╝╚═╝     ╚══════╝╚═╝  ╚═╝
` + TextColors.Reset)
}

func (l *logger) Clear() {
	fmt.Print("\033[H\033[2J")
}

func (l *logger) GetTime() string {
	now := time.Now()
	return now.Format("15:04:05")
}

func (l *logger) Warning(message string) {
	if l.enable_time_stamp {
		fmt.Printf("%s\r[%s] [%s] %s\n%s", TextColors.Yellow, l.prefix, l.GetTime(), message, TextColors.Reset)
	} else {
		fmt.Printf("%s\r[%s] %s\n%s", TextColors.Yellow, l.prefix, message, TextColors.Reset)
	}
}

func (l *logger) Error(message string) {
	if l.enable_time_stamp {
		fmt.Printf("%s\r[%s] [%s] %s\n%s", TextColors.Red, l.prefix, l.GetTime(), message, TextColors.Reset)
	} else {
		fmt.Printf("%s\r[%s] %s\n%s", TextColors.Red, l.prefix, message, TextColors.Reset)
	}
}

func (l *logger) Success(message string) {
	if l.enable_time_stamp {
		fmt.Printf("%s\r[%s] [%s] %s\n%s", TextColors.Green, l.prefix, l.GetTime(), message, TextColors.Reset)
	} else {
		fmt.Printf("%s\r[%s] %s\n%s", TextColors.Green, l.prefix, message, TextColors.Reset)
	}
}

func (l *logger) General(message string) {
	if l.enable_time_stamp {
		fmt.Printf("%s\r[%s] [%s] %s\n%s", TextColors.White, l.prefix, l.GetTime(), message, TextColors.Reset)
	} else {
		fmt.Printf("%s\r[%s] %s\n%s", TextColors.White, l.prefix, message, TextColors.Reset)
	}
}

func (l *logger) Notification(message string) {
	if l.enable_time_stamp {
		fmt.Printf("%s\r[%s] [%s] %s\n%s", TextColors.Magenta, l.prefix, l.GetTime(), message, TextColors.Reset)
	} else {
		fmt.Printf("%s\r[%s] %s\n%s", TextColors.Magenta, l.prefix, message, TextColors.Reset)
	}
}
