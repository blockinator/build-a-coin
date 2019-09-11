package web

import (
    "net/http"
    "buildacoin/data"
)

// Get the address of the requester; preferring X-Forwarded-For if it is
// present and it is expected to be behind a proxy.
func RequestOrigin(conf *data.Conf, req *http.Request) string {
    forwardedFor := req.Header.Get("X-Forwarded-For")
    if forwardedFor != "" && conf.Proxied() {
        return forwardedFor
    }
    return req.RemoteAddr
}
