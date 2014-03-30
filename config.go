package ghosts

import (
    "os"
    "io/ioutil"
    "errors"
    "fmt"
    "strings"
    "path/filepath"

    "github.com/go-yaml/yaml"
)

// Global environment
var Environment = struct {
    ConfigPath string
    TargetPath string
} {

}

func GetTargetPath() string {
    if Environment.TargetPath == "" {
        return GetSysTargetPath()
    }

    return Environment.TargetPath
}

func GetConfigPath() string {
    if Environment.ConfigPath == "" {
        p, err :=  GetSysConfigPath()
        if err != nil {
            panic(err)
        }

        return p
    }

    return Environment.ConfigPath
}

type Config struct {
    Path string
    data map[string]interface{}
    names map[string]HostNode
    pending map[string]bool
}

// Load the config file and parse it
func (this *Config) Parse() (map[string]HostNode, error) {
    content, err := ioutil.ReadFile(this.Path)
    if err != nil {
        return nil, err
    }

    err = yaml.Unmarshal(content, &this.data)
    if err != nil {
        return nil, err
    }

    this.names = make(map[string]HostNode)
    this.pending = make(map[string]bool)

    for k, _ := range this.data {
        n, err := this.ParseName(k)
        if err != nil {
            return nil, err
        }

        this.names[k] = n
    }

    return this.names, nil
}

// Parse a group entry
func (this *Config) ParseName(name string) (HostNode, error) {
    // parsed
    n, ok := this.names[name]
    if ok {
       return n, nil 
    }

    _, ok = this.pending[name]
    if ok {
        return nil, errors.New("Bad reference: " + name)
    }

    this.pending[name] = true

    defer delete(this.pending, name)

    data, _ := this.data[name]
    switch v := data.(type) {
    case string:
        // it's url
        if isUrl(v) {
            return &HostFileRemote {FileUrl: v}, nil
        }

        // it's name
        _, ok := this.data[v]
        if ok {
            return this.ParseName(v)
        }

        // test scope
        n := &HostScope {Scope: v}
        if err := n.Validate(); err == nil && len(n.GetList()) > 0 {
            return n, nil
        }

        // test file
        if f := detectFile(v); f != "" {
            n := &HostFile {FilePath: f}
            return n, nil
        }

        return nil, errors.New("Invalid config entry: " + v)
    case []interface{}:
        list, err := convertList(v)
        if err != nil {
            return nil, err
        }

        g := &HostGroup{}
        for _, v := range list {
            // it's url
            if isUrl(v) {
                g.Add(&HostFileRemote {FileUrl: v})
                continue
            }

            // it's name
            _, ok := this.data[v]
            if ok {
                n, err := this.ParseName(v)
                if err != nil {
                    return nil, err
                }
                g.Add(n)
                continue
            }

            // test scope
            n := &HostScope {Scope: v}
            if err := n.Validate(); err == nil && len(n.GetList()) > 0 {
                g.Add(n)
                continue
            }

            // test file
            if f := detectFile(v); f != "" {
                n := &HostFile {FilePath: f}
                g.Add(n)
                continue
            }

            return nil, errors.New("Invalid config entry: " + v)
        }

        return g, nil
    }
    return nil, errors.New("Invalid config entry: " + name)
}

// Convert []interface{} to []string
func convertList(l []interface{}) ([]string, error) {
    result := make([]string, len(l))
    for k, v := range l {
        s, ok :=  v.(string)
        if !ok {
            return nil, fmt.Errorf("Invalid config entry: %v", v)
        }

        result[k] = s
    }

    return result, nil
}

// Check if it's a file
func detectFile(s string) string {
    if strings.IndexAny(s, `?<>\:*|”\r\n`) != -1 {
        return ""
    }

    // convert to abs path (relative to the config file)
    if !filepath.IsAbs(s) {
        p := GetConfigPath()
        p = filepath.Dir(p)
        s = filepath.Join(p, s)
    }

    // not a file
    fi, err := os.Stat(s)
    if err != nil || fi.IsDir() {
        return ""
    }

    // ok, it's a file
    return s
}

// Simple url validation
func isUrl(s string) bool {
    if len(s) <= 8 {
        return false
    }
    if strings.ToLower(s[0:7]) == "http://" || strings.ToLower(s[0:8]) == "https://" {
        return true
    }
    return false
}