package common

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var (
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
)

func SetupServerLogger() func() {
	logFile, latest, cleanup, err := SetupLoggerFiles()
	if err != nil {
		panic(fmt.Sprintf("error setting up server logger: %s", err.Error()))
	}

	writer := io.MultiWriter(os.Stdout, logFile, latest)

	InfoLogger = log.New(writer, "INFO: ", log.LstdFlags)
	WarningLogger = log.New(writer, "WARNING: ", log.LstdFlags)
	ErrorLogger = log.New(writer, "ERROR: ", log.LstdFlags)

	return cleanup
}

const logDir = "logs/serverLogs"

func SetupLoggerFiles() (*os.File, *os.File, func(), error) {
	if _, err := os.Stat(logDir); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(logDir, os.ModePerm); err != nil {
			return nil, nil, nil, fmt.Errorf("error creating logger directory: %v", err.Error())
		}
	}

	fileNameDate := time.Now().Format("2006-01-02-15:04:05")
	duplicateTestDate := fileNameDate
	i := 0
	for {
		if _, err := os.Stat(fmt.Sprintf("%s/%s.log", logDir, duplicateTestDate)); errors.Is(err, os.ErrNotExist) {
			fileNameDate = duplicateTestDate
			break
		} else {
			i++
			duplicateTestDate = fmt.Sprintf("%s-%d", fileNameDate, i)
		}
	}

	logFile, latestFile, err := createLog(fileNameDate)
	if err != nil {
		return nil, nil, nil, err
	}

	return logFile, latestFile, func() {
		logFile.Close()
		latestFile.Close()
	}, nil
}

func createLog(time string) (log, latest *os.File, err error) {
	log, err = os.Create(fmt.Sprintf("%s/%s.log", logDir, time))
	if err != nil {
		return nil, nil, err
	}
	latest, err = os.Create(fmt.Sprintf("%s/latest.log", logDir))
	if err != nil {
		return log, nil, err
	}

	return
}
