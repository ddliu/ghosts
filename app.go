package ghosts

import (
    "sort"
    "errors"
    "io/ioutil"
)

type App struct {
    names map[string]HostNode
}

func (this *App) ParseConfig() error {
    c := &Config {
        Path: GetConfigPath(),
    }

    var err error
    this.names, err = c.Parse()
    if err != nil {
        return err
    }

    return nil
}

func (this *App) GetNames() []string {
    var names []string
    for k, _ := range this.names {
        names = append(names, k)
    }
    
    sort.Strings(names)

    return names
}

func (this *App) Generate(names ...string) (string, error) {
    g := &HostGroup{}
    for _, name := range names {
        node, ok := this.names[name]
        if !ok {
            return "", errors.New("Group does not exist: " + name)
        }
        g.Add(node)
    }

    if err := g.Validate(); err != nil {
        return "", err
    }

    return MergeHosts(g.GetList()...), nil
}

func (this *App) GenerateRaw(names ...string) (string, error) {
    g := &HostGroup{}
    for _, name := range names {
        node, ok := this.names[name]
        if !ok {
            return "", errors.New("Group does not exist: " + name)
        }
        g.Add(node)
    }

    if err := g.Validate(); err != nil {
        return "", err
    }

    return g.GetRaw(), nil
}

func (this *App) Switch(names ...string) error {
    content, err := this.Generate(names...)
    if err != nil {
        return err
    }

    return ioutil.WriteFile(GetTargetPath(), []byte(content), 0644)
}

func (this *App) SwitchRaw(names ... string) error {
    content, err := this.GenerateRaw(names...)
    if err != nil {
        return err
    }

    return ioutil.WriteFile(GetTargetPath(), []byte(content), 0644)
}