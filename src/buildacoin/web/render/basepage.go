package render

import (
    "buildacoin/data"
    "io"
    "net/http"
    "text/template"
)

const (
    DefaultPageTitle = "Build-a-Coin Cryptocurrency Creator"
)

type basePage struct {
    headerMarkup string
    style string
    markupTemplate *template.Template
}

func SimplePage(conf *data.Conf, markupAsset, styleAsset string) (*basePage,
        error) {
    var header string
    var style string

    // build template from markup assets
    markup := template.New("base")
    markupStr, err := data.StringAsset(conf, "markup/base.html")
    if err != nil {
        return nil, err
    }
    _, err = markup.Parse(markupStr)
    if err != nil {
        return nil, err
    }

    body := markup.New("body")
    if len(markupAsset) > 0 {
        markupStr, err = data.StringAsset(conf, markupAsset)
        if err != nil {
            return nil, err
        }
    } else {
        markupStr = ""
    }
    _, err = body.Parse(markupStr)
    if err != nil {
        return nil, err
    }

    // build style from css assets
    for _, styleName := range []string { "style/base.css", styleAsset } {
        if len(styleName) <= 0 {
            continue
        }
        styleStr, err := data.StringAsset(conf, styleName)
        if err != nil {
            return nil, err
        }
        style += styleStr
    }

    return &basePage { header, style, markup, }, nil
}

func (tt *basePage) Execute(out io.Writer, arg interface{},
        errs []error) error {
    return tt.markupTemplate.Execute(out, map[string]interface{} {
        "title": DefaultPageTitle,
        "style": tt.style,
        "content": arg,
        "errors": errs,
    })
}

func (tt *basePage) ServeHTTP(out http.ResponseWriter, req *http.Request) {
    tt.Execute(out, nil, nil)
}
