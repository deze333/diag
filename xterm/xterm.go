// Xterm color logging produces output to terminal
package xterm

import (
	"fmt"
	"strings"
	"time"
)

const (
	WHITE  = "\033[37m"
	CYAN   = "\033[36m"
	PURPLE = "\033[95m"
	BLUE   = "\033[94m"
	GREEN  = "\033[92m"
	YELLOW = "\033[93m"
	RED    = "\033[91m"
	CLEAR  = "\033[0m"

	UNDERLINE            = "\033[4m"
	UNDERLINE_WHITE      = "\033[4m\033[37m"
	UNDERLINE_BOLD_WHITE = "\033[4m\033[1m\033[37m"

	BOLD        = "\033[1m"
	BOLD_CYAN   = "\033[1m\033[36m"
	BOLD_PURPLE = "\033[1m\033[95m"
	BOLD_YELLOW = "\033[1m\033[93m"
	BOLD_RED    = "\033[1m\033[91m"

	INVERSE_WHITE  = "\033[7m\033[37m"
	INVERSE_PURPLE = "\033[7m\033[95m"
	INVERSE_BLUE   = "\033[7m\033[94m"
	INVERSE_GREEN  = "\033[7m\033[92m"
	INVERSE_YELLOW = "\033[7m\033[93m"
	INVERSE_RED    = "\033[7m\033[91m"
)

// DEBUG output
func DEBUG(time time.Time, name, title string, args ...interface{}) string {
	out := []string{}
	out = append(out, fmt.Sprintf("\n%s%s\n%s%s%s", CYAN, name, YELLOW, title, CLEAR))

	if len(args) == 1 {
		out = append(out, fmt.Sprintf(" %s%v%s", WHITE, args[0], CLEAR))
	} else {
		for i := 0; i+1 < len(args); i += 2 {
			k := args[i]
			v := args[i+1]
			if k == "" && v == "" {
				out = append(out, fmt.Sprintf("    %s*", WHITE))
				continue
			}
			out = append(out, fmt.Sprintf("    %s* %s%v = %s%v", WHITE, BLUE, k, WHITE, v))
			if i == len(args)-2 {
				out = append(out, CLEAR)
			}
		}
	}

	return strings.Join(out, "\n")
}

// NOTE output
func NOTE(time time.Time, msg string, args ...interface{}) string {
	out := []string{}
	out = append(out, fmt.Sprintf("%s%v%s", INVERSE_WHITE, msg, CLEAR))

	if len(args) == 1 {
		out = append(out, fmt.Sprintf(" %s%v%s", WHITE, args[0], CLEAR))
	} else {
		for i := 0; i+1 < len(args); i += 2 {
			k := args[i]
			v := args[i+1]
			if k == "" && v == "" {
				out = append(out, fmt.Sprintf("    %s*", WHITE))
				continue
			}
			out = append(out, fmt.Sprintf(" %s%v = %s%v", BLUE, k, WHITE, v))
		}
		out = append(out, CLEAR)
	}

	return strings.Join(out, "")
}

// NOTE2 output
func NOTE2(time time.Time, msg string, args ...interface{}) string {
	out := []string{}
	out = append(out, fmt.Sprintf("%s%v%s:", INVERSE_BLUE, msg, CLEAR))

	if len(args) == 1 {
		out = append(out, fmt.Sprintf(" %s%v%s", WHITE, args[0], CLEAR))
	} else {
		for i := 0; i+1 < len(args); i += 2 {
			out = append(out, fmt.Sprintf(" %s%v = %s%v", BLUE, args[i], WHITE, args[i+1]))
		}
		out = append(out, CLEAR)
	}

	return strings.Join(out, "")
}

// WARNING output
func WARNING(time time.Time, name, title string, args ...interface{}) string {
	out := []string{}
	out = append(out, fmt.Sprintf("\n%s%s\n%s%s%s", CYAN, name, INVERSE_YELLOW, title, CLEAR))

	if len(args) == 1 {
		out = append(out, fmt.Sprintf(" %s%v%s", WHITE, args[0], CLEAR))
	} else {
		for i := 0; i+1 < len(args); i += 2 {
			k := args[i]
			v := args[i+1]
			if k == "" && v == "" {
				out = append(out, fmt.Sprintf("    %s*", WHITE))
				continue
			}
			out = append(out, fmt.Sprintf("    %s* %s%v = %s%v", WHITE, BLUE, k, WHITE, v))
			if i == len(args)-2 {
				out = append(out, CLEAR)
			}
		}
	}

	return strings.Join(out, "\n")
}

// ERROR output
func ERROR(time time.Time, name, title string, args ...interface{}) string {
	out := []string{}
	out = append(out, fmt.Sprintf("\n%s%s\n%s%s%s", CYAN, name, INVERSE_RED, title, CLEAR))

	if len(args) == 1 {
		out = append(out, fmt.Sprintf(" %s%v%s", WHITE, args[0], CLEAR))
	} else {
		for i := 0; i+1 < len(args); i += 2 {
			k := args[i]
			v := args[i+1]
			if k == "" && v == "" {
				out = append(out, fmt.Sprintf("    %s*", WHITE))
				continue
			}
			out = append(out, fmt.Sprintf("    %s* %s%v = %s%v", WHITE, BLUE, k, WHITE, v))
			if i == len(args)-2 {
				out = append(out, CLEAR)
			}
		}
	}

	return strings.Join(out, "\n")
}
