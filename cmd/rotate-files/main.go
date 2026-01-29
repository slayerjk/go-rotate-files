package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	vafswork "github.com/slayerjk/go-vafswork"
)

const (
	appName = "rotate-files"
)

func main() {
	// defining default values
	var (
		workDir         string    = vafswork.GetExePath()
		logsPathDefault string    = workDir + "/logs" + "_" + appName
		startTime       time.Time = time.Now()
	)

	// flags
	logsDir := flag.String("log-dir", logsPathDefault, "set custom log dir")
	logsToKeep := flag.Int("keep-logs", 7, "set number of logs to keep after rotation")
	dirsToRotate := flag.String("d", "NONE", "REQUIRED, abs path of dir or dirs, separeted by coma")
	filesToKeep := flag.Int("r", -1, "REQUIRED, most recent files to keep")

	flag.Usage = func() {
		fmt.Println("Go: Rotate Files")
		fmt.Println("Version = 0.0.1")
		fmt.Println("Usage: <app> [-opt] ...")
		fmt.Println("Flags:")
		flag.PrintDefaults()
	}

	flag.Parse()

	// logging
	// create log dir
	if err := os.MkdirAll(*logsDir, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stdout, "failed to create log dir %s:\n\t%v", *logsDir, err)
		os.Exit(1)
	}
	// set current date
	dateNow := time.Now().Format("02.01.2006")
	// create log file
	logFilePath := fmt.Sprintf("%s/%s_%s.log", *logsDir, appName, dateNow)
	// open log file in append mode
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stdout, "failed to open created log file %s:\n\t%v", logFilePath, err)
		os.Exit(1)
	}
	defer logFile.Close()
	// set logger
	logger := slog.New(slog.NewTextHandler(logFile, nil))

	// starting programm notification
	logger.Info("Program Started", "app name", appName)

	// rotate logs
	logger.Info("Log rotation first", "logsDir", *logsDir, "logs to keep", *logsToKeep)
	if _, err := vafswork.RotateFilesByMtime(*logsDir, *logsToKeep); err != nil {
		fmt.Fprintf(os.Stdout, "failed to rotate logs:\n\t%v", err)
	}

	// main code here
	logger.Info("checking all REQUIRED flags are set")

	// check if dirsToRotate flag(-d) is set
	if *dirsToRotate == "NONE" {
		logger.Error("flag '-d' is not set, exiting")
		os.Exit(1)
	}

	// check if filesToRotate is positive int
	if *filesToKeep == -1 {
		logger.Error("flag '-r' is not set, exiting")
		os.Exit(1)
	}

	logger.Info("all REQUIRED flags are set, moving on")

	// getting slice of dirs
	dirsToRotateList := strings.Split(*dirsToRotate, ",")

	// iterating over dirsToRotateList for files rotation
	logger.Info("starting file rotation process")

	for _, dir := range dirsToRotateList {
		dirToRotate := strings.Trim(dir, " ")

		logger.Info("now rotating files", "dir", dir, "filesToKeeep", *filesToKeep)

		deletedFiles, err := vafswork.RotateFilesByMtime(dirToRotate, *filesToKeep)
		if err != nil {
			logger.Warn("failed rotating files, skipping dir", "dir", dir, "filesToKeeep", *filesToKeep, "err", err)
		}

		logger.Info("done rotating files", "dir", dir)
		for _, file := range deletedFiles {
			logger.Info("DELETED", "file", file)
		}
	}

	// count & print estimated time
	logger.Info("Program Done", slog.Any("estimated time(sec)", time.Since(startTime).Seconds()))
}
