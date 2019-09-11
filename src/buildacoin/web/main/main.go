package main

import (
    "buildacoin/data"
    "buildacoin/web"
    "flag"
)

func main() {
    var confPath string

    flag.StringVar(&confPath, "conf", "", "load conf file instead of default")
    flag.Parse()

    conf, err := data.LoadConfFromArg(confPath);
    if err != nil {
        panic(err.Error())
    }

    web.Serve(conf)
}
