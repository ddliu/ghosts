// +build darwin dragonfly freebsd linux netbsd openbsd

package ghosts

const CONFIG_DIRECTORY = ".ghosts"
const EOL = "\n"

func GetSysTargetPath() string {
    return "/etc/hosts"
}