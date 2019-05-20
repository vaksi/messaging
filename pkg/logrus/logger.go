package logrus

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)
func Init(filename string, debug bool, paths ...string) {
	SetFormatter(&log.JSONFormatter{})
	if debug {
		var dir []string
		// dw, _ := os.Getwd()
		for _, path := range paths {
			if path != "" {
				absin := absPathify(path)
				if !stringInSlice(absin, dir) {
					FileHandler(absin, filename) // log file handler
					dir = append(dir, absin)
				}
			}
		}
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func absPathify(inPath string) string {
	if strings.HasPrefix(inPath, "$") {
		end := strings.Index(inPath, string(os.PathSeparator))
		inPath = os.Getenv(inPath[1:end]) + inPath[end:]
	}

	if filepath.IsAbs(inPath) {
		return filepath.Clean(inPath)
	}

	p, err := filepath.Abs(inPath)
	if err == nil {
		return filepath.Clean(p)
	}
	return ""
}

// SetFormatter sets the standard logger formatter.
func SetFormatter(formatter log.Formatter) {
	log.SetFormatter(formatter)
}

// FileHandler handles log to file
func FileHandler(dir, filename string) {
	path := strings.Join([]string{dir, filename}, "/")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		Info("Create Dir", log.Fields{
			"path":  path,
			"error": err,
		})
		return
		// err = os.MkdirAll(dir, 0755)
		// if err != nil {
		//     panic(err)
		// }
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v", err))
	}
	SetOutput(file)
}

func getCallerInfo() (file string, line int, ok bool) {
	_, file, line, ok = runtime.Caller(3)
	return
}

// SetLevel sets the standard entry level
func SetLevel(levelStr string) error {
	level, err := log.ParseLevel(levelStr)
	if err != nil {
		return err
	}

	log.SetLevel(level)
	return nil
}

// SetOutput sets the standard logger output.
func SetOutput(output io.Writer) {
	log.SetOutput(output)
}

// An entry is the final or intermediate logging entry. It contains all
// the fields passed with WithField{,s}. It's finally logged when Debug, Info,
// Warn, Error, Fatal or Panic is called on it. These objects can be reused and
// passed around as much as you wish to avoid field duplication.
func entry() *log.Entry {
	file, line, _ := getCallerInfo()
	return log.WithFields(log.Fields{
		"on": fmt.Sprintf("%s:%d", file, line),
	})
}

func Debug(msg ...interface{}) {
	entry().Debug(msg...)
}

func Info(msg ...interface{}) {
	entry().Info(msg...)
}

func Warn(msg ...interface{}) {
	entry().Warn(msg...)
}

func Error(msg ...interface{}) {
	entry().Error(msg...)
}

func Fatal(msg ...interface{}) {
	entry().Fatal(msg...)
}