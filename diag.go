// Outputs colorful log diagnostics
package diag

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/deze333/diag/plain"
	"github.com/deze333/diag/util"
	"github.com/deze333/diag/xterm"
)

//------------------------------------------------------------
// Variables
//------------------------------------------------------------

type loggers struct {
	fnametpl string
	tstamp   time.Time
	timer    *time.Timer

	xtermLog *log.Logger

	plainFileDir string
	plainFile    *os.File
	plainLog     *log.Logger

	htmlFileDir string
	htmlFile    *os.File
	htmlLog     *log.Logger

	historySize int
}

var _logger *loggers

//------------------------------------------------------------
// Init
//------------------------------------------------------------

func minStart() {
	if _logger != nil && _logger.xtermLog != nil {
		_logger.xtermLog.Print(xterm.WARNING(
			time.Now(),
			"diag",
			"package diag config not provided, assuming screen only output"))
	} else {
		fmt.Println("[diag] package diag config not provided, assuming screen only output")
	}
	_logger = &loggers{}
	_logger.xtermLog = log.New(os.Stdout, "", 0)
}

//------------------------------------------------------------
// API
//------------------------------------------------------------

func SetHistory(size int) {
	_logger.historySize = size
}

func Start(directory string, filename string, xterm, plain, html bool) (err error) {
	_logger = &loggers{}
	_logger.historySize = 3

	// Default screen output
	if xterm {
		_logger.xtermLog = log.New(os.Stdout, "", 0)
	}

	if filename == "" || directory == "" {
		return
	}
	_logger.fnametpl = filename

	// Mark start time
	_logger.tstamp = time.Now()

	// Add timer to rotate logs at the end of day
	if plain || html {
		_logger.timer = time.AfterFunc(rotationDelta(_logger.tstamp), rotateLogs)
	}

	// Create time stamped filename: RFC3339 = "2006-01-02T15:04:05Z07:00"
	//filename = strings.Replace(filename, "{}", _logger.tstamp.Format(time.RFC3339), 1)
	filename = strings.Replace(filename, "{}", "", 1)

	// Log file for plain
	if plain {
		_logger.plainFileDir = path.Join(directory, "plain")
		err := os.MkdirAll(_logger.plainFileDir, 0775)
		if err != nil {
			return err
		}
		f, err := os.Create(path.Join(_logger.plainFileDir, filename))
		if err != nil {
			return err
		}
		_logger.plainFile = f
		_logger.plainLog = log.New(f, "", 0)
	}

	// Log file for HTML
	if html {
		_logger.htmlFileDir = path.Join(directory, "html")
		err := os.MkdirAll(_logger.htmlFileDir, 0775)
		if err != nil {
			return err
		}
		f, err := os.Create(path.Join(_logger.htmlFileDir, filename))
		if err != nil {
			return err
		}
		_logger.htmlFile = f
		_logger.htmlLog = log.New(f, "", 0)
	}

	return
}

// Rotates logs
func rotateLogs() {
	DEBUG("diag", "Rotating logs", "closing time stamp", _logger.tstamp.Format(time.ANSIC))
	// Close current logs
	// Rename defaut logs that are about to be closed timestamped
	tstamp := "_" + _logger.tstamp.Format(time.Stamp)
	filename := strings.Replace(_logger.fnametpl, "{}", "", 1)

	// Plain log
	if _logger.plainFile != nil {
		// Stop logging and close file
		f := _logger.plainFile
		_logger.plainFile = nil
		f.Close()
		// Rename file
		err := os.Rename(
			f.Name(),
			path.Join(_logger.plainFileDir, strings.Replace(_logger.fnametpl, "{}", tstamp, 1)))
		if err != nil {
			SOS("diag", "Error renaming plain log file. Plain logging stopped.", "msg", err)
		} else {
			// Create new logging file with default name (ie, webapp.log)
			f, err := os.Create(path.Join(_logger.plainFileDir, filename))
			if err != nil {
				SOS("diag", "Error creating plain log file. Plain logging stopped.", "msg", err)
			} else {
				// Start logging
				_logger.plainFile = f
				_logger.plainLog = log.New(f, "", 0)
			}
		}
	}

	// Html log
	if _logger.htmlFile != nil {
		f := _logger.htmlFile
		_logger.htmlFile = nil
		f.Close()
	}

	// Set new logging start time
	_logger.tstamp = time.Now()
	_logger.timer = time.AfterFunc(rotationDelta(_logger.tstamp), rotateLogs)

	// Add first log record
	DEBUG("diag", "New log started", "opening time stamp", _logger.tstamp.Format(time.ANSIC))

	// Clean up old logs
	cleanLogs(_logger.plainFileDir, _logger.historySize)
	cleanLogs(_logger.htmlFileDir, _logger.historySize)
}

// Calculates delta time from given time
// to the end of cycle
func rotationDelta(t time.Time) time.Duration {
	// Advance to next day
	t2 := t.Add(time.Hour * 24)
	// Event will take place next day 23:59:00
	t2 = time.Date(
		t2.Year(),
		t2.Month(),
		t2.Day(),
		//23, 59, 0, 0,
		12, 00, 0, 0,
		t2.Location())

	// TEMPORARY TEST PLUG: Set to very short period
	//t2 = t.Add(time.Second * 60 * 10)

	//fmt.Printf("!!!!!! T1 = %v\n", t)
	//fmt.Printf(">>>>>> T2 = %v\n", t2)
	//fmt.Printf("###### ROTATION DELTA = %v\n", t2.Sub(t))

	return t2.Sub(t)
}

// Cleans logs directory by removing
// all log files that are older than last N logs.
// If historySize < 0 then no logs deleted.
func cleanLogs(dir string, historySize int) {
	if dir == "" || historySize < 0 {
		return
	}

	// Read directory and sort files with most recent on top
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		SOS("diag", "Error cleaning log directory", "err", err, "dir", dir)
		return
	}
	sort.Sort(FilesByDate(fis))

	// Delete files that exceed given history size
	for i := historySize + 1; i < len(fis); i++ {
		if err := os.Remove(path.Join(dir, fis[i].Name())); err != nil {
			SOS("diag", "Error deleting old log file", "err", err, "dir", dir, "file", fis[i].Name())
		}
	}
}

// Sorting of files:
// This type allows sorting of a slice of FileInfo
// by modification date, most recent on top.
type FilesByDate []os.FileInfo

func (f FilesByDate) Len() int {
	return len(f)
}
func (f FilesByDate) Less(i, j int) bool {
	return f[i].ModTime().Unix() > f[j].ModTime().Unix()
}
func (f FilesByDate) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

// Close all file based log output.
// No further log file writes will happen.
// Screen output will still work.
func Close() {
	DEBUG("diag", "Closing log file output")

	if _logger.plainFile != nil {
		_logger.plainFile.Close()
		_logger.plainFile = nil
	}
	if _logger.htmlFile != nil {
		_logger.htmlFile.Close()
		_logger.htmlFile = nil
	}
}

// Print log
func Print(v ...interface{}) {
	// Xterm screen log
	if _logger.xtermLog != nil {
		_logger.xtermLog.Print(v...)
	}

	// Plain file output
	if _logger.plainLog != nil {
		_logger.plainLog.Print(v...)
	}

	// HTML file output
	if _logger.htmlLog != nil {
		_logger.htmlLog.Print(v...)
	}
}

// Prinft log
func Printf(format string, v ...interface{}) {
	// Xterm screen log
	if _logger.xtermLog != nil {
		_logger.xtermLog.Printf(format, v...)
	}

	// Plain file output
	if _logger.plainLog != nil {
		_logger.plainLog.Printf(format, v...)
	}

	// HTML file output
	if _logger.htmlLog != nil {
		_logger.htmlLog.Printf(format, v...)
	}
}

// Outputs debug message to at least screen logger.
// If file based loggers were configured then
// they will record that message too.
func DEBUG(name, title string, v ...interface{}) {
	if _logger == nil {
		minStart()
	}

	t := time.Now()

	// Xterm screen log
	if _logger.xtermLog != nil {
		_logger.xtermLog.Print(xterm.DEBUG(t, name, title, v...))
	}

	// Plain file output
	if _logger.plainLog != nil {
		_logger.plainLog.Print(plain.DEBUG(t, name, title, v...))
	}

	// HTML file output
	if _logger.htmlLog != nil {
		//_logger.htmlLog.Printf(format, v...)
	}
}

// Simple NOTE
func NOTE(msg string, v ...interface{}) {
	if _logger == nil {
		minStart()
	}

	t := time.Now()

	// Xterm screen log
	if _logger.xtermLog != nil {
		_logger.xtermLog.Print(xterm.NOTE(t, msg, v...))
	}

	// Plain file output
	if _logger.plainLog != nil {
		_logger.plainLog.Print(plain.NOTE(t, msg, v...))
	}

	// HTML file output
	if _logger.htmlLog != nil {
		//_logger.htmlLog.Printf(format, v...)
	}
}

// Simple NOTE 2 (Inverse color)
func NOTE2(msg string, v ...interface{}) {
	if _logger == nil {
		minStart()
	}

	t := time.Now()

	// Xterm screen log
	if _logger.xtermLog != nil {
		_logger.xtermLog.Print(xterm.NOTE2(t, msg, v...))
	}

	// Plain file output
	if _logger.plainLog != nil {
		_logger.plainLog.Print(plain.NOTE(t, msg, v...))
	}

	// HTML file output
	if _logger.htmlLog != nil {
		//_logger.htmlLog.Printf(format, v...)
	}
}

// Outputs WARNING message
func WARNING(name, title string, v ...interface{}) {
	if _logger == nil {
		minStart()
	}

	t := time.Now()

	// Xterm screen log
	if _logger.xtermLog != nil {
		_logger.xtermLog.Print(xterm.WARNING(t, name, title, v...))
	}

	// Plain file output
	if _logger.plainLog != nil {
		_logger.plainLog.Print(plain.WARNING(t, name, title, v...))
	}

	// HTML file output
	if _logger.htmlLog != nil {
		//_logger.htmlLog.Printf(format, v...)
	}
}

// Outputs ERROR message
func ERROR(name, title string, v ...interface{}) {
	if _logger == nil {
		minStart()
	}

	t := time.Now()

	// Xterm screen log
	if _logger.xtermLog != nil {
		_logger.xtermLog.Print(xterm.ERROR(t, name, title, v...))
	}

	// Plain file output
	if _logger.plainLog != nil {
		_logger.plainLog.Print(plain.ERROR(t, name, title, v...))
	}

	// HTML file output
	if _logger.htmlLog != nil {
		//_logger.htmlLog.Printf(format, v...)
	}
}

// Outputs SOS message to at least screen logger.
// And attempts to immediately contact a human.
// If file based loggers were configured then
// they will record that message too.
// NEW: Add "stack" as the last of v and stack trace will be appended.
func SOS(name, title string, v ...interface{}) {
	if len(v) != 0 && fmt.Sprint(v[len(v)-1]) == "stack" {
		v = append(v, util.Stack())
	}
	notifyEmail(name, title, v...)
	ERROR(name, title, v...)
}

func SOS_Stack(name, title string, v ...interface{}) {
	v = append(v, "stack")
	v = append(v, util.Stack())
	notifyEmail(name, title, v...)
	ERROR(name, title, v...)
}
