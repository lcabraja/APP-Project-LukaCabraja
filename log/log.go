package log

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

const defaultJoin = ", "
const verbose = true

var logLevelWidth = 5

func Dev(data ...interface{}) {
	var caller *string = nil
	if verbose {
		pc, _, _, _ := runtime.Caller(1)
		callingFunc := runtime.FuncForPC(pc).Name()
		caller = &callingFunc
	}
	log(DEV, caller, data...)
}

func Devf(format string, data ...interface{}) {
	var caller *string = nil
	if verbose {
		pc, _, _, _ := runtime.Caller(1)
		callingFunc := runtime.FuncForPC(pc).Name()
		caller = &callingFunc
	}
	logf(DEV, caller, format, data...)
}

func D(data ...interface{}) {
	var caller *string = nil
	if verbose {
		pc, _, _, _ := runtime.Caller(1)
		callingFunc := runtime.FuncForPC(pc).Name()
		caller = &callingFunc
	}
	log(DEBUG, caller, data...)
}

func Df(format string, data ...interface{}) {
	var caller *string = nil
	if verbose {
		pc, _, _, _ := runtime.Caller(1)
		callingFunc := runtime.FuncForPC(pc).Name()
		caller = &callingFunc
	}
	logf(DEBUG, caller, format, data...)
}

func I(data ...interface{}) {
	var caller *string = nil
	if verbose {
		pc, _, _, _ := runtime.Caller(1)
		callingFunc := runtime.FuncForPC(pc).Name()
		caller = &callingFunc
	}
	log(INFO, caller, data...)
}

func If(format string, data ...interface{}) {
	var caller *string = nil
	if verbose {
		pc, _, _, _ := runtime.Caller(1)
		callingFunc := runtime.FuncForPC(pc).Name()
		caller = &callingFunc
	}
	logf(INFO, caller, format, data...)
}

func E(data ...interface{}) {
	var caller *string = nil
	if verbose {
		pc, _, _, _ := runtime.Caller(1)
		callingFunc := runtime.FuncForPC(pc).Name()
		caller = &callingFunc
	}
	log(ERROR, caller, data...)
}

func Ef(format string, data ...interface{}) {
	var caller *string = nil
	if verbose {
		pc, _, _, _ := runtime.Caller(1)
		callingFunc := runtime.FuncForPC(pc).Name()
		caller = &callingFunc
	}
	logf(ERROR, caller, format, data...)
}

func F(data ...interface{}) {
	pc, _, _, _ := runtime.Caller(1)
	callingFunc := runtime.FuncForPC(pc).Name()
	caller := strings.Split(callingFunc, ".")[1]
	log(FATAL, &caller, data...)
	os.Exit(1)
}

func Ff(format string, data ...interface{}) {
	pc, _, _, _ := runtime.Caller(1)
	caller := runtime.FuncForPC(pc).Name()
	logf(FATAL, &caller, format, data...)
	os.Exit(1)
}

func log(level LogLevel, caller *string, data ...interface{}) {
	ds := stringify(defaultJoin, data...)
	logf(level, caller, "%s", ds)
}

func logf(level LogLevel, caller *string, format string, data ...interface{}) {
	var final string

	time := dye(getTime(), Green, BackgroundReset, true)
	ds := fmt.Sprintf(format, data...)
	dss := strings.TrimSpace(ds)
	if caller != nil {
		sc := strings.LastIndex(*caller, "/")
		*caller = dye((*caller)[sc+1:], Magenta, BackgroundReset, true)
		final = fmt.Sprintf("%s [%s] {%s}: %s", time, l(level), *caller, dss)
	} else {
		final = fmt.Sprintf("%s [%s] %s", time, l(level), dss)
	}

	output(final)
}

func l(ll LogLevel) string {
	l := ll.String()
	if len(l) > logLevelWidth {
		logLevelWidth = len(l)
	}
	return dye(fmt.Sprintf("%-*s", logLevelWidth, l), ll.Color(), ll.Background(), true)
}

func getTime() string {
	t := time.Now()
	return fmt.Sprint(t.Format("2006-01-02_15:04:05"))
}

func stringify(join string, v ...interface{}) string {
	var s []string
	for _, val := range v {
		s = append(s, fmt.Sprintf("%v", val))
	}
	j := strings.Join(s, join)
	return strings.ReplaceAll(j, "\n", "_")
}

func output(s string) {
	var (
		final string
		trunc string
	)
	final = strings.ReplaceAll(s, "\n", "_")

	width, _, _ := -1, 0, 0

	if width > 0 && len(final) > width {
		trunc = final[:width]
	} else {
		trunc = final
	}

	fmt.Println(trunc)
}
