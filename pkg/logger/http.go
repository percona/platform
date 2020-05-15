package logger

import (
	"net/http"
	"net/http/httputil"
)

// LoggingRoundTripper returns http.RoundTripper with request/response logger.
func LoggingRoundTripper(rt http.RoundTripper, printf func(format string, args ...interface{})) http.RoundTripper {
	return &loggingRoundTripper{
		rt:     rt,
		printf: printf,
	}
}

type loggingRoundTripper struct {
	rt     http.RoundTripper
	printf func(format string, v ...interface{})
}

func (lrt *loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	b, _ := httputil.DumpRequestOut(req, true)
	if len(b) != 0 {
		lrt.printf("Request:\n%s", b)
	}

	resp, err := lrt.rt.RoundTrip(req)

	if resp != nil {
		b, _ = httputil.DumpResponse(resp, true)
		if len(b) != 0 {
			lrt.printf("Response:\n%s", b)
		}
	}

	return resp, err
}
