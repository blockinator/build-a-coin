package render

import (
    "buildacoin/data"
    "errors"
    "io"
    "net/http"
)

type NotFoundPage struct {
    conf *data.Conf
    fancyPage *basePage
}

func NewNotFoundPage(conf *data.Conf) NotFoundPage {
    fancyPage, _ := SimplePage(conf, "", "")
    return NotFoundPage { conf, fancyPage }
}

func (tt NotFoundPage) ServeHTTP(out http.ResponseWriter, req *http.Request) {
    out.WriteHeader(http.StatusNotFound)

    errorString := "not found!"

    if tt.fancyPage != nil {
        tt.fancyPage.Execute(out, nil, []error { errors.New(errorString) })
    } else {
        _, err := io.WriteString(out, errorString + "\n")
        if err != nil {
            // TODO log
        }
    }
}
