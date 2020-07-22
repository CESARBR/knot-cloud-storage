package logging

import (
	"log/syslog"
	"os"

	"github.com/sirupsen/logrus"
	logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"
)

// Logrus represents the logrus logger
type Logrus struct {
	level  string
	syslog bool
}

// NewLogrus creates a new logrus instance
func NewLogrus(level string, syslog bool) *Logrus {
	return &Logrus{level, syslog}
}

// Get returns a logrus instance based on the specific context
func (l *Logrus) Get(context string) *logrus.Entry {
	log := logrus.New()
	log.Out = os.Stderr
	level, _ := logrus.ParseLevel(l.level)
	log.SetLevel(level)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	if l.syslog {
		setupSyslog(log)
	}

	logger := log.WithFields(logrus.Fields{
		"Context": context,
	})

	return logger
}

func setupSyslog(log *logrus.Logger) {
	hook, err := logrus_syslog.NewSyslogHook("", "", syslog.LOG_INFO, "storage")
	if err != nil {
		log.Error("Unable to connect to local syslog daemon")
		return
	}

	log.AddHook(hook)
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: false,
	})
}
