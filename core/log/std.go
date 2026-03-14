package log

import "gServ/core/config"

const (
	std_error_prefix_text = "ERROR"
	std_info_prefix_text  = "INFO"
	std_warn_prefix_text  = "WARN"
	std_debug_prefix_text = "DEBUG"
)

func StdErrorf(format string, v ...any) {
	std_error_logger.Printf(format, v...)
}

func StdInfof(format string, v ...any) {
	std_info_logger.Printf(format, v...)
}

func StdWarnf(format string, v ...any) {
	std_warn_logger.Printf(format, v...)
}

func StdDebugf(format string, v ...any) {
	if config.GetConfig().Server.Mode == config.SERVER_MODE_DEV {
		std_debug_logger.Printf(format, v...)
	}
}
