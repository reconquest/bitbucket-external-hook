# bitbucket-external-hook

## Install

You can build the program using `go get`
```
go get -v github.com/reconquest/bitbucket-external-hook
```

Or you can use GoBinaries.com service to obtain the binary:

```
curl -sf https://gobinaries.com/reconquest/bitbucket-external-hook | sh
```


## Usage

```
bitbucket-external-hook [options] list -b <bitbucket-uri> -p <project> [-r <repo>]
bitbucket-external-hook [options] print -b <bitbucket-uri> -p <project> [-r <repo>] <hook>
bitbucket-external-hook [options] enable -b <bitbucket-uri> -p <project> [-r <repo>] <hook>
bitbucket-external-hook [options] disable -b <bitbucket-uri> -p <project> [-r <repo>] <hook>
bitbucket-external-hook [options] set -b <bitbucket-uri> -p <project> [-r <repo>] <hook> [-e <path>] [-s] [<param>...]
bitbucket-external-hook count -b <bitbucket-uri> [--debug] [<hook>]
```

```
  -b <bitbucket-uri>      URI to Bitbucket, can include auth info.
  -p <project>            Project key.
  -r <repository>         Repository key, optional.
  <hook>                  Hook key (Java thing).
  -o --only-enabled       Show only enabled hooks.
  -c --only-configured    Show only configured hooks.
  -e --executable <path>  Path to hook executable.
  -s --safepath           Look for <path> in safe directory (bitbucket home).
  <param>                 Add param to hook, can be specified multiple times.
```

`-b <bitbucket-uri>` can include authorization credentials, like as following:
`-b http://admin:password@bitbucket.local:4000/`

### List hooks in a project

```
bitbucket-external-hook list -b http://admin:admin@bitbucket.local:1234 -p hooking
```

### List hooks in a repository

```
bitbucket-external-hook list -b http://admin:admin@bitbucket.local:1234 -p hooking -r withhooks
```

### Print hook info/settings in a repository

```
bitbucket-external-hook print -b http://admin:admin@bitbucket.local:1234 -p hooking -r withhooks com.ngs.stash.externalhooks.external-hooks:external-pre-receive-hook
```

### Enable a hook in a project

```
bitbucket-external-hook enable -b http://admin:admin@bitbucket.local:1234 -p hooking com.ngs.stash.externalhooks.external-hooks:external-pre-receive-hook
```

### Enable a hook in a repository

```
bitbucket-external-hook enable -b http://admin:admin@bitbucket.local:1234 -p hooking -r withhooks com.ngs.stash.externalhooks.external-hooks:external-pre-receive-hook
```

### Disable a hook in a project

```
bitbucket-external-hook disable -b http://admin:admin@bitbucket.local:1234 -p hooking com.ngs.stash.externalhooks.external-hooks:external-pre-receive-hook
```

### Disable a hook in a repository

```
bitbucket-external-hook disable -b http://admin:admin@bitbucket.local:1234 -p hooking -r withhooks com.ngs.stash.externalhooks.external-hooks:external-pre-receive-hook
```

### Set hook settings (ExternalHooks)

```
bitbucket-external-hook set -b http://admin:admin@bitbucket.local:1234 -p hooking -r withhooks com.ngs.stash.externalhooks.external-hooks:external-pre-receive-hook \
    -e test1.sh \
    -s \
    param1 param2
```

- `-e` means path to executable
- `-s` means to look for test1.sh in safe directory (bitbucket home or shared/ in dc)
- `param1` and `param2` are just strings that will be contactenated into one
    string with newlines and sent as `Params` field.

### Measure hooks usage

Print a total number of repositories where hooks are enabled and configured or
inherited from project level.

```
bitbucket-external-hook count -b http://admin:admin@bitbucket.local
```

Also, you can specify a hook prefix:

```
bitbucket-external-hook count -b http://admin:admin@bitbucket.local com.ngs.stash.externalhooks.external-hooks
```
