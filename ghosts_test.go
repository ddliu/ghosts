package ghosts

import (
    "testing"
)

func TestEntry(t *testing.T) {
    n := &HostEntry {
        Host: "localhost",
        IP: "127.0.0.1",
    }

    err := n.Validate()
    if err != nil {
        t.Error(err)
    }

    t.Log(n.GetRaw())
}

func TestScope(t *testing.T) {
    n := &HostScope {
        Scope: `
        # comment
        127.0.0.1   localhost # localhost

        192.168.0.111   google.com

        `,
    }

    err := n.Validate()
    if err != nil {
        t.Error(err)
    }

    list := n.GetList()
    if len(list) != 2 || list[0].IP != "127.0.0.1" || list[0].Host != "localhost" {
        t.Error("Invalid result")
    }
}

func TestFile(t *testing.T) {
    n := &HostFile {
        FilePath: "tests/hosts",
    }

    err := n.Validate()
    if err != nil {
        t.Error(err)
    }

    list := n.GetList()
    if len(list) != 3 || list[0].IP != "192.168.1.200" || list[0].Host != "ns.example.com" {
        t.Error("Invalid result")
    }
}

func TestRemoteFile(t *testing.T) {
    n := &HostFileRemote {
        FileUrl: "https://raw.githubusercontent.com/ddliu/hosts/master/ghosts.example.hosts",
    }

    err := n.Validate()
    if err != nil {
        t.Error(err)
    }

    list := n.GetList()
    if len(list) != 2 || list[0].IP != "127.0.0.1" || list[0].Host != "localhost" {
        t.Error("Invalid result")
    }
}

func TestGroup(t *testing.T) {
    g := &HostGroup {
    }
    g.Add(&HostEntry {
        Host: "a",
        IP: "192.168.0.1",
    }, &HostEntry {
        Host: "b",
        IP: "192.168.0.2",
    }, &HostEntry {
        Host: "c",
        IP: "192.168.0.3",
    })

    err := g.Validate()
    if err != nil {
        t.Error(err)
    }

    list := g.GetList()
    if len(list) != 3 || list[2].IP != "192.168.0.3" || list[2].Host != "c" {
        t.Error("Invalid result")
    }
}

// TODO: Config path should not be system wide?
// func TestConfig(t *testing.T) {
//     c := &Config{Path: "tests/ghosts.yml"}
//     _, err := c.Parse()
//     if err != nil {
//         t.Error(err)
//     }
// }

func TestApp(t *testing.T) {
    Environment.ConfigPath = "tests/ghosts.yml"
    Environment.TargetPath = "tests/_hosts"

    app := &App{}
    err := app.ParseConfig()
    if err != nil {
        t.Error(err)
    }
    err = app.Switch("base", "local")

    if err != nil {
        t.Error(err)
    }
}