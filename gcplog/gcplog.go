package gcplog

import (
	"time"

	logger "github.com/sirupsen/logrus"
)

const (
	SeverityDEFAULT   = "DEFAULT"
	SeverityDEBUG     = "DEBUG"
	SeverityINFO      = "INFO"
	SeverityNOTICE    = "NOTICE"
	SeverityWARNING   = "WARNING"
	SeverityERROR     = "ERROR"
	SeverityCRITICAL  = "CRITICAL"
	SeverityALERT     = "ALERT"
	SeverityEMERGENCY = "EMERGENCY"
)

// LogLevelMap NOTICE and EMERGENCY is not mapped, map if necessary
var LogLevelMap = map[logger.Level]string{
	logger.PanicLevel: SeverityALERT,
	logger.FatalLevel: SeverityCRITICAL,
	logger.ErrorLevel: SeverityERROR,
	logger.WarnLevel:  SeverityWARNING,
	logger.InfoLevel:  SeverityINFO,
	logger.DebugLevel: SeverityDEBUG,
	logger.TraceLevel: SeverityDEBUG,
}

type gcpFormatter struct {
	jsonFormatter logger.JSONFormatter
}

func NewGCPFormatter() *gcpFormatter {
	return &gcpFormatter{
		jsonFormatter: logger.JSONFormatter{
			FieldMap: logger.FieldMap{
				logger.FieldKeyTime:  "timestampLocal",
				logger.FieldKeyLevel: "logrusLevel",
				logger.FieldKeyMsg:   "message",
			},
		},
	}
}

func (f *gcpFormatter) Format(entry *logger.Entry) ([]byte, error) {
	gcpSeverity, ok := LogLevelMap[entry.Level]
	if !ok {
		gcpSeverity = SeverityDEFAULT
	}
	entry.Data["severity"] = gcpSeverity
	entry.Data["timestamp"] = entry.Time.UTC().Format(time.RFC3339Nano)
	return f.jsonFormatter.Format(entry)
}
