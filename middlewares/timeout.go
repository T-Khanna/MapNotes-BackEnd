package middlewares

import (
	"net/http"
	"time"
)

const _TIMEOUT_SECS time.Duration = 5

/*
 Middleware that causes request to timeout after `_TIMEOUT_SECS` seconds.
*/
func Timeout(h http.Handler) http.Handler {
	return http.TimeoutHandler(h, _TIMEOUT_SECS*time.Second, "Request timed out.")
}
