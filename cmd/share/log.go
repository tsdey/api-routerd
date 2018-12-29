// SPDX-License-Identifier: Apache-2.0

package share

import (
	"github.com/sirupsen/logrus"
	"os"
	"path"
)

const (
	defaultLogDir = "/var/log/api-router"
	defaultLogFile = "api-router.log"
)

func InitLog() {
	log := logrus.New()

	log.Out = os.Stderr
	log.Level = logrus.DebugLevel

	logDir := defaultLogDir

	err := CreateDirectory(logDir, 0644)
	if (err != nil) {
		log.Errorf("Failed to create log directory. path: %s, err: %s", logDir, err)
		return
	}

	logFile := path.Join(logDir, defaultLogFile)
	f, err := os.OpenFile(logFile, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0644)
	if err != nil {
		log.Errorf("Failed to create log file. path: %s, err: %s", logFile, err)
	}

	log.SetOutput(f)
	log.SetReportCaller(true)
	log.Info("Starting API Router")
}
