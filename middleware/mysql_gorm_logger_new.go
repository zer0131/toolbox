package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/zer0131/toolbox/log"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

func NewLoggerMe(config logger.Config) *loggerMe {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	if config.Colorful {
		infoStr = logger.Green + "%s\n" + logger.Reset + logger.Green + "[info] " + logger.Reset
		warnStr = logger.BlueBold + "%s\n" + logger.Reset + logger.Magenta + "[warn] " + logger.Reset
		errStr = logger.Magenta + "%s\n" + logger.Reset + logger.Red + "[error] " + logger.Reset
		traceStr = logger.Green + "%s\n" + logger.Reset + logger.Yellow + "[%.3fms] " + logger.BlueBold + "[rows:%v]" + logger.Reset + " %s"
		traceWarnStr = logger.Green + "%s " + logger.Yellow + "%s\n" + logger.Reset + logger.RedBold + "[%.3fms] " + logger.Yellow + "[rows:%v]" + logger.Magenta + " %s" + logger.Reset
		traceErrStr = logger.RedBold + "%s " + logger.MagentaBold + "%s\n" + logger.Reset + logger.Yellow + "[%.3fms] " + logger.BlueBold + "[rows:%v]" + logger.Reset + " %s"
	}

	return &loggerMe{
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

type loggerMe struct {
	logger.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

// LogMode log mode
func (l *loggerMe) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l loggerMe) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		log.Infof(ctx, l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l loggerMe) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		log.Warnf(ctx, l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l loggerMe) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		log.Errorf(ctx, l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
func (l loggerMe) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {

	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			log.Errorf(ctx, l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Errorf(ctx, l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			log.Warnf(ctx, l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Warnf(ctx, l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.LogLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			log.Infof(ctx, l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Infof(ctx, l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
