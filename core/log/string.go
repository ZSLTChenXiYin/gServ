package log

const (
	color_red    = "\033[91m"
	color_green  = "\033[92m"
	color_yellow = "\033[93m"
	color_blue   = "\033[94m"
	color_white  = "\033[0m"
)

func StdRedString(str string) string {
	return color_red + str + color_white
}

func StdGreenString(str string) string {
	return color_green + str + color_white
}

func StdYellowString(str string) string {
	return color_yellow + str + color_white
}

func StdBlueString(str string) string {
	return color_blue + str + color_white
}
