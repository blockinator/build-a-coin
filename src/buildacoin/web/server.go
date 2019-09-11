package web

import (
    "buildacoin/data"
    "buildacoin/web/render"
    "net/http"
    "strconv"
)

const (
    // URL base inside which all the explicit coin pages are kept
    CoinPagesRoot = "/coin/"
)

// Serve the build-a-coin web interface
func Serve(conf *data.Conf) {
    mux := http.NewServeMux()

    listenAt := conf.ListenAddr4() + ":" + strconv.Itoa(int(conf.ListenPort()))

    // static pages
    for _, pageName := range []string { "faq", "about", "help" } {
        handler, err := render.SimplePage(conf, "markup/" + pageName + ".html",
            "style/" + pageName + ".css")
        if err == nil {
            mux.Handle("/" + pageName, handler)
        } else {
            // TODO log
        }
    }

    baseCoinCount := conf.BaseCoinCount()
    if baseCoinCount < 1 {
        handleRoot(conf, mux, render.NewErrorPage(conf, "no coins!?"))
    } else {
        // wire up a page for each base coin
        for ii := 0; ii < conf.BaseCoinCount(); ii++ {
            name := conf.BaseCoin(ii)

            base, err := data.LoadMeta(conf, name)
            if err != nil {
                http.ListenAndServe(listenAt, render.NewErrorPage(conf,
                    "error loading base coin '" + name + "': " + err.Error()))
                return
            }

            baseHandler, err := render.NewCoinPage(conf, base)
            if err != nil {
                http.ListenAndServe(listenAt, render.NewErrorPage(conf,
                    "error loading base coin handler '" + name + "': " +
                    err.Error()))
                return
            }

            // the first base coin is the default, route the root page to it
            if ii == 0 {
                handleRoot(conf, mux, baseHandler)
            }
            mux.Handle(CoinPagesRoot + base.Id(), baseHandler)
        }
    }

    http.ListenAndServe(listenAt, mux)
}

// assigning a handler to / makes all pages 'exist' without a hack?  For shame,
// Go.
type hackyRootHandler struct {
    child http.Handler
    notFound http.Handler
}
func (tt hackyRootHandler) ServeHTTP(out http.ResponseWriter,
        req *http.Request) {
    if req.URL.Path != "/" {
        tt.notFound.ServeHTTP(out, req)
        return
    }
    tt.child.ServeHTTP(out, req)
}
func handleRoot(conf *data.Conf, mux *http.ServeMux, handler http.Handler) {
    mux.Handle("/", hackyRootHandler { handler, render.NewNotFoundPage(conf) })
}
