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
    if len(list) != 2 || list[0].IP != "127.0.0.1" || list[0].Host != "localhost" {
        t.Error("Invalid result")
    }
}

func TestMain(t *testing.T) {
    // g := &HostGroup {}

}