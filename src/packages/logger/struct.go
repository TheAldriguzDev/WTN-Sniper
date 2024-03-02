package console_logger

type txt_color struct {
	Reset   string
	Black   string
	Red     string
	Green   string
	Yellow  string
	Blue    string
	Magenta string
	Cyan    string
	White   string
}

type logger struct {
	prefix            string
	enable_time_stamp bool
}
