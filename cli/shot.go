package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var Log *Logger

type LogLevel int

const (
	ERROR LogLevel = iota
	WARN
	INFO
	DEBUG
	TRACE
)

var levelNames = map[LogLevel]string{
	ERROR: "ERROR",
	WARN:  "WARN ",
	INFO:  "INFO ",
	DEBUG: "DEBUG",
	TRACE: "TRACE",
}

type Logger struct {
	level  LogLevel
	logger *log.Logger
}

func NewLogger(level LogLevel) *Logger {
	return &Logger{
		level:  level,
		logger: log.New(os.Stderr, "", 0),
	}
}

func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

func (l *Logger) log(level LogLevel, format string, args ...any) {
	if level <= l.level {
		timestamp := time.Now().Format("15:04:05")
		prefix := fmt.Sprintf("[%s] %s ", timestamp, levelNames[level])
		message := fmt.Sprintf(format, args...)
		l.logger.Printf("%s%s", prefix, message)
	}
}

	Log = NewLogger(level)
}
func init() {
	Log = NewLogger(ERROR)
}

var ErrInvalidShebang = errors.New("not a valid shebang line")

// parseShebang parses a shebang line and returns the interpreter name followed by any arguments
//
// The line must start with #! (leading/trailing whitespace).
// If it does not, ErrInvalidShebang is returned.
//
// Returns nil for empty lines, invalid shebangs, or lines that don't start with an interpreter path.
func parseShebang(line string) ([]string, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, ErrInvalidShebang
	}

	if !strings.HasPrefix(line, "#!") {
		return nil, ErrInvalidShebang
	}

	lineWithoutShebang := line[2:]
	if lineWithoutShebang == "" {
		return nil, ErrInvalidShebang
	}

	parts := strings.Fields(lineWithoutShebang)
	if len(parts) == 0 {
		return nil, ErrInvalidShebang
	}

	prog := parts[0]

	result := make([]string, 0, len(parts))
	result = append(result, prog)

	if len(parts) > 1 {
		result = append(result, parts[1:]...)
	}

	return result, nil
}

func main() {
	flag.Parse()
	ParseVerbosity()

	// Example log messages at different levels
	Log.Error("This is an error message (always shown)")
	Log.Warn("This is a warning message (shown with -v)")
	Log.Info("This is an info message (shown with -v -v)")
	Log.Debug("This is a debug message (shown with -v -v -v)")
	Log.Trace("This is a trace message (shown with -v -v -v -v)")

	// Example of using the logger in your application
	Log.Info("Application started")
	Log.Debug("Connecting to database...")
	Log.Trace("SQL: SELECT * FROM users")
	Log.Info("Application finished")

	// input, err := io.ReadAll(os.Stdin)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "shot: cannot read script: %v\n", err)
	// 	os.Exit(1)
	// }
	// log.Printf("shot: read %v bytes\n", len(input))

	// sheBangLine := strings.SplitN(string(input), "\n", 2)[0]
	// parsedShebang, err := parseShebang(sheBangLine)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "shot: cannot parse shebang: %q: %v\n", sheBangLine, err)
	// 	os.Exit(1)
	// }

	// cmd := exec.Command(parsedShebang[0], parsedShebang[1:]...)
	// cmd.Args = append(cmd.Args, os.Args[1:]...)
	// cmd.Stdin = bytes.NewReader(input)
	// cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	// err = cmd.Run()
	// if err != nil {
	// 	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() >= 0 {
	// 		os.Exit(exitErr.ExitCode())
	// 	}
	// 	fmt.Fprintf(os.Stderr, "shot: %v\n", err)
	// 	os.Exit(1)
	// }
}

// func extractShotConfig(r io.Reader) map[string]string {
// 	config := make(map[string]string)
// 	scanner := bufio.NewScanner(r)
// 	inBlock := false

// 	for scanner.Scan() {
// 		line := strings.TrimSpace(scanner.Text())

// 		if line == "# shot: configuration" || line == "# shot: config" {
// 			inBlock = true
// 			continue
// 		}
// 		if inBlock && (line == "# shot: end" || line == "# shot: configuration-end") {
// 			break
// 		}

// 		if !inBlock {
// 			continue
// 		}

// 		if !strings.HasPrefix(line, "#") {
// 			continue
// 		}
// 		content := strings.TrimSpace(strings.TrimPrefix(line, "#"))
// 		if !strings.HasPrefix(content, "SHOT_") {
// 			continue
// 		}

// 		parts := strings.SplitN(content, "=", 2)
// 		if len(parts) != 2 {
// 			continue
// 		}

// 		key := strings.TrimSpace(parts[0])
// 		val := strings.TrimSpace(parts[1])
// 		val = strings.Trim(val, `"'`)

// 		config[key] = val
// 	}

// 	return config
// }
