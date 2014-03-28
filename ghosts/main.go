package main

import (
    "flag"

    "github.com/ddliu/ghosts"
)

var list bool
var config string
var target string
var raw bool
var printOut bool

func main() {
    flag.BoolVar(&list, "list", false, "Show group list")
    flag.StringVar(&config, "config", "", "Path to the config file (optional)")
    flag.StringVar(&target, "target", "", "Path to the host file (optional)")
    flag.BoolVar(&raw, "raw", false, "Keep comments and empty lines?")
    flag.BoolVar(&printOut, "print", false, "Print out results instead of write to the host file")

    flag.Parse()

    if config != "" {
        ghosts.Environment.ConfigPath = config
    }

    if target != "" {
        ghosts.Environment.TargetPath = target
    }

    app := &ghosts.App{}

    err := app.ParseConfig()
    if err != nil {
        fmt.Println("Parse config file error:", err)
        return
    }

    if list {
        for _, v := range app.GetNames() {
            fmt.Println(v)
        }
        return
    }

    if printOut {
        if raw {
            content, err := app.GenerateRaw()
            if err != nil
        }
        return
    }

    if raw {

    }
    flag.Args()
}