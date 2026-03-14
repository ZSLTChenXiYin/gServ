package log

import (
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"
)

var (
	std_error_logger = newErrorLogger()
	std_info_logger  = newInfoLogger()
	std_warn_logger  = newWarnLogger()
	std_debug_logger = newDebugLogger()

	zap_logger *zap.Logger
)

func Init(path string) error {
	var err error
	// 创建 zap logger
	zap_logger, err = loadZapLogger(path)
	if err != nil {
		return fmt.Errorf("zap日志初始化错误: %v", err)
	}

	return nil
}

func newErrorLogger() *log.Logger {
	return log.New(os.Stderr, fmt.Sprintf("[%s] ", StdRedString(std_error_prefix_text)), log.LstdFlags)
}

func newInfoLogger() *log.Logger {
	return log.New(os.Stdout, fmt.Sprintf("[%s] ", StdGreenString(std_info_prefix_text)), log.LstdFlags)
}

func newWarnLogger() *log.Logger {
	return log.New(os.Stdout, fmt.Sprintf("[%s] ", StdYellowString(std_warn_prefix_text)), log.LstdFlags)
}

func newDebugLogger() *log.Logger {
	return log.New(os.Stdout, fmt.Sprintf("[%s] ", StdBlueString(std_debug_prefix_text)), log.LstdFlags)
}

func GetStdErrorLogger() *log.Logger {
	return std_error_logger
}

func GetStdInfoLogger() *log.Logger {
	return std_info_logger
}

func GetStdWarnLogger() *log.Logger {
	return std_warn_logger
}

func GetStdDebugLogger() *log.Logger {
	return std_debug_logger
}

func GetZapLogger() *zap.Logger {
	return zap_logger
}

func GetZapGormLogger() *ZapGormLogger {
	return NewZapGormLogger(zap_logger)
}
