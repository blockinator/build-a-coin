package render

import (
    "buildacoin/data"
    "errors"
    "io"
    "net/http"
)

// show a detailed error page when debugging, and a nice safe generic error
// otherwise
type ErrorPage struct {
    conf *data.Conf
    reason string
    fancyPage *basePage
}

func NewErrorPage(conf *data.Conf, reason string) ErrorPage {
    fancyPage, _ := SimplePage(conf, "", "")
    return ErrorPage { conf, reason, fancyPage }
}

func (tt ErrorPage) ServeHTTP(out http.ResponseWriter, req *http.Request) {
    out.WriteHeader(http.StatusInternalServerError)

    var errorString string
    if tt.conf.Debug() {
        errorString = tt.reason
    } else {
        errorString = "something went horribly wrong!"
    }

    if tt.fancyPage != nil {
        tt.fancyPage.Execute(out, nil, []error { errors.New(errorString) })
    } else {
        _, err := io.WriteString(out, errorString + "\n")
        if err != nil {
            // TODO log
        }
    }
}
