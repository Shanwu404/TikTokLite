package logger

import (
	"errors"

	"go.uber.org/zap"
)

var (
	loggerInfo  *zap.SugaredLogger
	loggerDebug *zap.SugaredLogger
	loggerError *zap.SugaredLogger
)

func init() {
	var err error
	loggerInfoConfig := zap.NewDevelopmentConfig()
	loggerInfoConfig.Encoding = "console"
	loggerInfoConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	loggerInfoConfig.OutputPaths = []string{"log/info.log"}
	loggerInfo_, err := loggerInfoConfig.Build(zap.AddCallerSkip(1))
	errHandler(err)
	loggerInfo = loggerInfo_.Sugar()

	loggerDebugConfig := zap.NewDevelopmentConfig()
	loggerDebugConfig.Encoding = "console"
	loggerDebugConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	loggerDebugConfig.OutputPaths = []string{"log/debug.log"}
	loggerDebug_, err := loggerDebugConfig.Build(zap.AddCallerSkip(1))
	errHandler(err)
	loggerDebug = loggerDebug_.Sugar()

	loggerErrorConfig := zap.NewDevelopmentConfig()
	loggerErrorConfig.Encoding = "console"
	loggerErrorConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	loggerErrorConfig.DisableStacktrace = true
	// debug 或者 release一样
	loggerErrorConfig.OutputPaths = []string{"log/error.log", "stderr"}
	loggerError_, err := loggerErrorConfig.Build(zap.AddCallerSkip(1))
	errHandler(err)
	loggerError = loggerError_.Sugar()
}

func Sync() []error {
	errs := make([]error, 0, 3)
	errs = append(errs, loggerInfo.Sync())
	errs = append(errs, loggerDebug.Sync())
	errs = append(errs, loggerError.Sync())
	return errs
}

func errHandler(err error) {
	if err != nil {
		panic(errors.Join(errors.New("创建logger失败"), err))
	}
}
func Infoln(args ...any) {
	loggerInfo.Infoln(args...)
}
func Debugln(args ...any) {
	loggerDebug.Debugln(args...)
}
func Errorln(args ...any) {
	loggerError.Errorln(args...)
}
func Infof(template string, args ...any) {
	loggerInfo.Infof(template, args...)
}
func Debugf(template string, args ...any) {
	loggerDebug.Debugf(template, args...)
}
func Errorf(template string, args ...any) {
	loggerError.Errorf(template, args...)
}
