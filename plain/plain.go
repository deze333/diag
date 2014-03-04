// Plain logging produces output to plain text log file
package plain

import (
	"fmt"
	"strings"
    "time"
)

const (
    sep = "------------------------------------------------------------"
    sepSos = "============================================================"
)

func DEBUG(t time.Time, name, title string, args ...interface{}) string {
	out := []string{
        sep,
        t.Format(time.ANSIC),
    }

	out = append(out, fmt.Sprintf("\"%s\"\n%s", name, title))

    if len(args) == 1 {
        out = append(out, fmt.Sprintf(" %v", args[0]))
    } else {
        for i := 0; i + 1 < len(args); i += 2 {
            out = append(out, fmt.Sprintf("    * %s = %s", args[i], args[i+1]))
        }
    }

	return strings.Join(out, "\n")
}

func NOTE(t time.Time, msg string, args ...interface{}) string {
	out := []string{
        sep,
        t.Format(time.ANSIC),
    }
    out = append(out, "\n")
    out = append(out, fmt.Sprintf(">>> %s:\n", msg))

    if len(args) == 1 {
        out = append(out, fmt.Sprintf(" %v", args[0]))
    } else {
        for i := 0; i + 1 < len(args); i += 2 {
            out = append(out, fmt.Sprintf("* %s = %s\n", args[i], args[i+1]))
        }
    }

	return strings.Join(out, "")
}

func WARNING(t time.Time, name, title string, args ...interface{}) string {
	out := []string{
        sepSos,
        t.Format(time.ANSIC),
    }

    out = append(out, fmt.Sprintf("\"%s\"\n!!! WARNING: %s", name, title))

    if len(args) == 1 {
        out = append(out, fmt.Sprintf(" %v", args[0]))
    } else {
        for i := 0; i + 1 < len(args); i += 2 {
            out = append(out, fmt.Sprintf("    * %s = %s", args[i], args[i+1]))
        }
    }

	return strings.Join(out, "\n")
}
func ERROR(t time.Time, name, title string, args ...interface{}) string {
	out := []string{
        sepSos,
        t.Format(time.ANSIC),
    }

    out = append(out, fmt.Sprintf("\"%s\"\n!!! ERROR: %s", name, title))

    if len(args) == 1 {
        out = append(out, fmt.Sprintf(" %v", args[0]))
    } else {
        for i := 0; i + 1 < len(args); i += 2 {
            out = append(out, fmt.Sprintf("    * %s = %s", args[i], args[i+1]))
        }
    }

	return strings.Join(out, "\n")
}
