package logger

import "github.com/sirupsen/logrus"

func init() {
	formatter := &logrus.TextFormatter{
		DisableTimestamp:       false,
		TimestampFormat:        "2006-01-02 15:04:05",
		DisableColors:          false,
		QuoteEmptyFields:       true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
		FullTimestamp:          false,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "@timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyFunc:  "caller",
		},
	}

	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(formatter)
}
