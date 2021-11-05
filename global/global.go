// Package global automatically sets http.DefaultTransport to loghttp.DefaultTransport when loaded.
package global

import (
	"net/http"

	"github.com/yuta-hidaka/go-loghttp"
)

func init() {
	http.DefaultTransport = loghttp.DefaultTransport
}
