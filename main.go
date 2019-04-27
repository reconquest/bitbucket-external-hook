package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/reconquest/karma-go"
)

var (
	version = "[manual build]"
	usage   = "bitbucket-external-hook " + version + `

Usage:
  bitbucket-external-hook [options] list -b <bitbucket-uri> -p <project> [-r <repo>]
  bitbucket-external-hook [options] print -b <bitbucket-uri> -p <project> [-r <repo>] <hook>
  bitbucket-external-hook [options] enable -b <bitbucket-uri> -p <project> [-r <repo>] <hook>
  bitbucket-external-hook [options] disable -b <bitbucket-uri> -p <project> [-r <repo>] <hook>
  bitbucket-external-hook [options] set -b <bitbucket-uri> -p <project> [-r <repo>] <hook> [-e <path>] [-s] [<param>...]
  bitbucket-external-hook -h | --help
  bitbucket-external-hook --version

Options:
  -b <bitbucket-uri>      URI to Bitbucket, can include auth info.
  -p <project>            Slug of project.
  -r <repository>         Slug of repository.
  <hook>                  Hook key.           
  -h --help               Show this screen.
  -o --only-enabled       Show only enabled hooks.
  -c --only-configured    Show only configured hooks.
  -e --executable <path>  Set hook executable.
  -s --safepath           Look in safe path.
  <param>                 Use param for hook.
  --version               Show version.
`
)

type (
	Options struct {
		BitbucketURI string `docopt:"-b"`
		Project      string `docopt:"-p"`
		Repository   string `docopt:"-r"`
		Hook         string `docopt:"<hook>"`

		List           bool
		OnlyEnabled    bool
		OnlyConfigured bool

		Print   bool
		Enable  bool
		Disable bool

		Set        bool
		Executable string
		SafePath   bool     `docopt:"--safepath"`
		Params     []string `docopt:"<param>"`
	}
)

func main() {
	args, err := docopt.ParseArgs(usage, nil, version)
	if err != nil {
		panic(err)
	}

	var opts Options
	err = args.Bind(&opts)
	if err != nil {
		log.Fatal(err)
	}

	remote, err := GetRemote(opts)
	if err != nil {
		log.Fatalf(err, "unable to parse URI")
	}

	api := NewAPI(remote)

	switch {
	case opts.List:
		err = handleList(api, opts)
	case opts.Print:
		err = handlePrint(api, opts)
	case opts.Enable:
		err = handleEnable(api, opts)
	case opts.Disable:
		err = handleDisable(api, opts)
	case opts.Set:
		err = handleSet(api, opts)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func handleSet(api *API, opts Options) error {
	settings := &HookSettings{
		Exe:      opts.Executable,
		SafePath: opts.SafePath,
		Params:   strings.Join(opts.Params, "\n"),
	}

	err := api.SetHookSettings(opts.Hook, settings)
	if err != nil {
		return karma.Format(
			err,
			"unable to set hook settings",
		)
	}

	printHookSettings(settings)

	return nil
}

func handleEnable(api *API, opts Options) error {
	err := api.EnableHook(opts.Hook)
	if err != nil {
		return karma.Format(
			err,
			"unable to enable hook",
		)
	}

	return nil
}

func handleDisable(api *API, opts Options) error {
	err := api.DisableHook(opts.Hook)
	if err != nil {
		return karma.Format(
			err,
			"unable to enable hook",
		)
	}

	return nil
}

func handlePrint(api *API, opts Options) error {
	hook, err := api.GetHook(opts.Hook)
	if err != nil {
		return karma.Format(
			err,
			"unable to get hook",
		)
	}

	printHook(hook)

	if !hook.Configured && !hook.Enabled {
		return nil
	}

	settings, err := api.GetHookSettings(opts.Hook)
	if err != nil {
		return karma.Format(
			err,
			"unable to get hook settings",
		)
	}

	fmt.Println()

	printHookSettings(settings)

	return nil
}

func handleList(api *API, opts Options) error {
	hooks, err := api.GetHooks()
	if err != nil {
		return karma.Format(
			err,
			"unable to get hooks",
		)
	}

	sort.SliceStable(hooks, func(i, j int) bool {
		return hooks[i].Details.Key < hooks[j].Details.Key
	})

	id := 0
	for _, hook := range hooks {
		if opts.OnlyEnabled && !hook.Enabled {
			continue
		}
		if opts.OnlyConfigured && !hook.Configured {
			continue
		}

		if id > 0 {
			fmt.Println()
		}
		id++

		printHook(hook)
	}

	return nil
}

func printHook(hook *Hook) {
	fmt.Printf("Key: %v\n", hook.Details.Key)
	fmt.Printf("Name: %v\n", hook.Details.Name)
	fmt.Printf("Type: %v\n", hook.Details.Type)
	fmt.Printf("Version: %v\n", hook.Details.Version)
	fmt.Printf("Scope: %v\n", hook.Scope.Type)
	fmt.Printf("Configured: %v\n", hook.Configured)
	fmt.Printf("Enabled: %v\n", hook.Enabled)
}

func printHookSettings(settings *HookSettings) {
	fmt.Printf("Executable: %v\n", settings.Exe)
	fmt.Printf("SafePath: %v\n", settings.SafePath)

	paramsPrefix := "Params: "
	fmt.Printf(
		"%v%v\n",
		paramsPrefix,
		strings.ReplaceAll(settings.Params, "\n", "\n"+strings.Repeat(" ", len(paramsPrefix))),
	)
}
