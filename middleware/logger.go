package middleware

import (
	"net/http"
	"os"

	apachelog "github.com/lestrrat-go/apache-logformat"
)

const (
	// LogFormatCommon apache common log format
	LogFormatCommon string = `%h %l %u %t "%r" %>s %b`
	// LogFormatCombined apache combined log format
	LogFormatCombined string = `%h %l %u %t "%r" %>s %b "%{Referer}i" "%{User-agent}i"`
)

// Logger Log request with apache format
func Logger(format string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log, err := apachelog.New(format)
		if err != nil {
			panic("format log not valid")
		}

		return log.Wrap(next, os.Stderr)
	}
}
