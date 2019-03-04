package philifence

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strings"
)

var slugger = regexp.MustCompile("[^a-z0-9]+")

func info(format string, vals ...interface{}) {
	log.Printf("INFO: "+format, vals...)
}

func warn(err error, extra string) {
	if err != nil {
		log.Printf("WARN: %s - %s", err, extra)
	}
}

func fatal(format string, vals ...interface{}) {
	msg := sprintf("FATAL: "+format, vals...)
	log.Fatal(msg)
}

func check(err error) {
	if err != nil {
		panic(err)
		fatal("%s", err)
	}
}

func assert(ok bool) {
	assert2(ok, "assertion failed!")
}

func assert2(ok bool, msg string, args ...interface{}) {
	if !ok {
		panic(fmt.Sprintf(msg, args...))
	}
}

func sprintf(format string, vals ...interface{}) string {
	return fmt.Sprintf(format, vals...)
}

func errorf(format string, vals ...interface{}) error {
	return fmt.Errorf(format, vals...)
}

func sluggify(path string) string {
	s := filepath.Base(path)
	s = strings.TrimSuffix(s, filepath.Ext(s))
	return strings.Trim(slugger.ReplaceAllString(strings.ToLower(s), "-"), "-")
}
