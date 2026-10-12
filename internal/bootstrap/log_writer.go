package bootstrap

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	upaths "github.com/neosy/elengrab/internal/app/utils/paths"
	iconfig "github.com/neosy/elengrab/internal/config"
	"github.com/neosy/elengrab/internal/pkg/filex"
)

var (
	startupFileName = strings.ToLower(iconfig.AppName) + "-startup.log"
	runtimeFileName = strings.ToLower(iconfig.AppName) + ".log"
)

type LogWriter struct {
	mu          sync.Mutex
	console     io.Writer
	startupFile *os.File
	runtimeFile *os.File
}

func NewLogWriter(appRoot string, config iconfig.ElengrabLoggingConfig) (*LogWriter, error) {
	writer := &LogWriter{
		console:     os.Stdout,
		startupFile: nil,
		runtimeFile: nil,
	}

	if !config.StartupEnabled && !config.RuntimeEnabled {
		return writer, nil
	}

	var path string

	if config.Dir != "" {
		var err error
		path, err = filex.AbsPath(appRoot, config.Dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to resolve log directory: %v\n", err)
			path = ""
		}
	}

	if path != "" {
		err := upaths.EnsureDir(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to ensure log directory: %v\n", err)
			path = ""
		}
	}

	if path == "" {
		fmt.Fprintln(os.Stderr, "File logging is disabled: log directory is unavailable.")
		return writer, nil
	}

	if config.StartupEnabled {
		filePath := filepath.Join(path, startupFileName)

		logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to open %s: %v\n", startupFileName, err)
		} else {
			writer.startupFile = logFile
		}
	}

	if config.RuntimeEnabled {
		filePath := filepath.Join(path, runtimeFileName)

		logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to open %s: %v\n", runtimeFileName, err)
		} else {
			writer.runtimeFile = logFile
		}
	}

	return writer, nil
}

func (w *LogWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	n, err := w.console.Write(p)
	if err != nil {
		return n, err
	}

	if w.startupFile != nil {
		if _, err := w.startupFile.Write(p); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to write to %s: %v\n", startupFileName, err)
			w.startupFile.Close()
			w.startupFile = nil
		}
	}

	if w.runtimeFile != nil {
		if _, err := w.runtimeFile.Write(p); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to write to %s: %v\n", runtimeFileName, err)
			w.runtimeFile.Close()
			w.runtimeFile = nil
		}
	}

	return n, nil
}

func (w *LogWriter) StopStartupFileLogging() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.startupFile == nil {
		return nil
	}

	err := w.startupFile.Close()
	w.startupFile = nil

	return err
}

func (w *LogWriter) StopRuntimeFileLogging() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.runtimeFile == nil {
		return nil
	}

	err := w.runtimeFile.Close()
	w.runtimeFile = nil

	return err
}

func (w *LogWriter) Close() error {
	err := w.StopStartupFileLogging()
	err = errors.Join(err, w.StopRuntimeFileLogging())

	return err
}
