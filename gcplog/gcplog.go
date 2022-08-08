package gcplog

import (
	"os"
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

func InitLogger(debug bool) {
	if _, ok := os.LookupEnv("K_SERVICE"); ok {
		logger.Info("[initLogger] K_SERVICE set, GCP environment detected, hooking logging logic...")
		logger.SetFormatter(NewGCPFormatter())
		logger.Info("[initLogger] Logging logic redirected.")
	} else {
		logger.Info("[initLogger] K_SERVICE not set, not in GCP environment, using normal logging.")
	}
	if debug {
		logger.SetLevel(logger.DebugLevel)
		logger.Debug("[initLogger] Enabled debug logging.")
	}
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
