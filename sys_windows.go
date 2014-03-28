package ghosts

const CONFIG_DIRECTORY = "ghosts"

const EOL = "\r\n"

func GetSysTargetPath() string {
    // todo: %sysroot%
    return `c:\windows\system32\drivers\etc\hosts`
}