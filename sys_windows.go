package ghosts

import (
	"os/user"
	"path/filepath"
)

const CONFIG_DIRECTORY = "ghosts"

const EOL = "\r\n"

func GetSysTargetPath() string {
    // todo: %sysroot%
    return `c:\windows\system32\drivers\etc\hosts`
}

// Get home dir of current user
func GetSysHomePath() (string, error) {
    u, err := user.Current()
    if err != nil {
        return "", err
    }

    return u.HomeDir, nil
}

func GetSysConfigPath() (string, error) {
    p, err := GetSysHomePath()
    if err != nil {
        return err
    }

    return filepath.Join(p, CONFIG_DIRECTORY, "ghosts.yml")
}