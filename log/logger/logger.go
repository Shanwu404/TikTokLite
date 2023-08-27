package logger

import (
	"errors"

	"go.uber.org/zap"
)

var (
	loggerInfo  *zap.SugaredLogger
	loggerDebug *zap.SugaredLogger
	loggerError *zap.SugaredLogger

	Infoln  = loggerInfo.Infoln
	Debugln = loggerDebug.Debugln
	Errorln = loggerError.Errorln

	Infof  = loggerInfo.Infof
	Debugf = loggerDebug.Debugf
	Errorf = loggerError.Errorf
)

func init() {
	var err error
	loggerInfoConfig := zap.NewDevelopmentConfig()
	loggerInfoConfig.Encoding = "console"
	loggerInfoConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	loggerInfoConfig.OutputPaths = []string{"log/info.log"}
	loggerInfo_, err := loggerInfoConfig.Build()
	errHandler(err)
	loggerInfo = loggerInfo_.Sugar()

	loggerDebugConfig := zap.NewDevelopmentConfig()
	loggerDebugConfig.Encoding = "console"
	loggerDebugConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	loggerDebugConfig.OutputPaths = []string{"log/debug.log"}
	loggerDebug_, err := loggerDebugConfig.Build()
	errHandler(err)
	loggerDebug = loggerDebug_.Sugar()

	loggerErrorConfig := zap.NewDevelopmentConfig()
	loggerErrorConfig.Encoding = "console"
	loggerErrorConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	// debug 或者 release一样
	loggerErrorConfig.OutputPaths = []string{"log/error.log", "stderr"}
	loggerError_, err := loggerErrorConfig.Build()
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
