package ghosts

import (
    "errors"
    "regexp"
    "strings"
    "bytes"
    "net/http"

    "io"
    "bufio"
    "io/ioutil"
)

// http://stackoverflow.com/questions/106179/regular-expression-to-match-hostname-or-ip-address
const ValidIpAddressRegex = `^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`;
const ValidHostnameRegex = `^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`;
var ValidIpAddressRegexCompile *regexp.Regexp
var ValidHostnameRegexCompile *regexp.Regexp

const HostLineRegex = `^([a-zA-Z0-9\._-]+)\s+([a-zA-Z0-9\._-]+)$`
var HostLineRegexCompile *regexp.Regexp


type HostNode interface {
    Validate() error
    GetList() []HostStruct
    GetRaw() string
}

type HostStruct struct {
    Host string
    IP string
}

// Simple entry node
type HostEntry struct {
    Host string
    IP string
}

func (this *HostEntry) Validate() error {
    if !isValidHost(this.Host) {
        return errors.New("Invalid host name: " + this.Host)
    }

    if !isValidIP(this.IP) {
        return errors.New("Invalid IP address: " + this.IP)
    }

    return nil
}

func (this *HostEntry) GetList() []HostStruct {
    return []HostStruct {
        HostStruct {
            Host: this.Host,
            IP: this.IP,
        },
    }
}

func (this *HostEntry) GetRaw() string {
    return this.IP + "\t" + this.Host
}


// The scope node, could be content of a host file
type HostScope struct {
    Scope string
    parsed []HostStruct
}

func (this *HostScope) Validate() (err error) {
    r := strings.NewReader(this.Scope)
    this.parsed, err = parseReader(r)
    
    return err
}

func (this *HostScope) GetList() []HostStruct {
    return this.parsed
}

func (this *HostScope) GetRaw() string {
    return this.Scope
}

// The file node
type HostFile struct {
    FilePath string
    parsed []HostStruct
    content []byte
}

func (this *HostFile) Validate() (err error) {
    this.content, err = ioutil.ReadFile(this.FilePath)
    if err != nil {
        return err
    }

    r := bytes.NewReader(this.content)

    this.parsed, err = parseReader(r)

    return err
}

func (this *HostFile) GetList() []HostStruct {
    return this.parsed
}

func (this *HostFile) GetRaw() string {
    return string(this.content)
}

// The remote file node
type HostFileRemote struct {
    FileUrl string
    parsed []HostStruct
    content []byte
}

func (this *HostFileRemote) Validate() error {
    resp, err := http.Get(this.FileUrl)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return errors.New("GET " + this.FileUrl + " failed with status: " + resp.Status)
    }
    this.content, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    r := bytes.NewReader(this.content)
    this.parsed, err = parseReader(r)

    return err
}

func (this *HostFileRemote) GetList() []HostStruct {
    return this.parsed
}

func (this *HostFileRemote) GetRaw() string {
    return string(this.content)
}

type HostGroup struct {
    Children []HostNode
}

func (this *HostGroup) AddNode(n HostNode) {
    this.Children = append(this.Children, n)
}

func (this *HostGroup) Validate() error {
    for _, n := range this.Children {
        err := n.Validate()
        if err != nil {
            return err
        }
    }

    return nil
}

func (this *HostGroup) GetList() []HostStruct {
    var hosts []HostStruct
    for _, n := range this.Children {
        hosts = append(hosts, n.GetList()...)
    }

    return hosts
}

func (this *HostGroup) GetRaw() string {
    var s []string

    for _, n := range this.Children {
        s = append(s, n.GetRaw())
    }

    return strings.Join(s, "\n")
}

func isValidHost(host string) bool {
    if ValidHostnameRegexCompile == nil {
        ValidHostnameRegexCompile = regexp.MustCompile(ValidHostnameRegex)
    }

    return ValidHostnameRegexCompile.MatchString(host)
}

func isValidIP(ip string) bool {
    if ValidIpAddressRegexCompile == nil {
        ValidIpAddressRegexCompile = regexp.MustCompile(ValidIpAddressRegex)
    }

    return ValidIpAddressRegexCompile.MatchString(ip)
}

// Parse a reader, might be string, opened file or opened http response body
func parseReader(r io.Reader) ([]HostStruct, error) {
    if HostLineRegexCompile == nil {
        HostLineRegexCompile = regexp.MustCompile(HostLineRegex)
    }

    var hosts []HostStruct

    scanner := bufio.NewScanner(r)
    for scanner.Scan() {
        line := scanner.Text()

        // Remove spaces and comments
        i := strings.IndexRune(line, '#')
        if i != -1 {
            line = line[0:i]
        }

        line = strings.Trim(line, " \t\r")

        // Only parse none empty line
        if line != "" {
            matches := HostLineRegexCompile.FindStringSubmatch(line)
            if len(matches) == 0 {
                return nil, errors.New("Invalid line: " + line)
            }

            ip, host := matches[1], matches[2]
            if !isValidHost(host) {
                return nil, errors.New("Invalid hostname: " + host)
            }

            if !isValidIP(ip) {
                return nil, errors.New("Invalid IP: " + ip)
            }

            hosts = append(hosts, HostStruct {
                Host: host,
                IP: ip,
            })
        }
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return hosts, nil
}