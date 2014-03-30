// +build darwin dragonfly freebsd linux netbsd openbsd

package ghosts

import (
    "os"
    "os/user"
    "path/filepath"
)

const CONFIG_DIRECTORY = ".ghosts"
const EOL = "\n"

func GetSysTargetPath() string {
    return "/etc/hosts"
}

// Get home dir of current user (or sudo_user)
func GetSysHomePath() (string, error) {
    if sudoUser := os.Getenv("SUDO_USER"); sudoUser != "" {
        u, err := user.Lookup(sudoUser)
        if err != nil {
            return "", err
        }

        return u.HomeDir, nil
    }

    u, err := user.Current()
    if err != nil {
        return "", err
    }

    return u.HomeDir, nil
}

func GetSysConfigPath() (string, error) {
    p, err := GetSysHomePath()
    if err != nil {
        return "", err
    }

    return filepath.Join(p, CONFIG_DIRECTORY, "ghosts.yml"), nil
}