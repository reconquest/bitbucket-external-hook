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
  bitbucket-external-hook [options] print -b <bitbucket-uri> -p <project> [-r <repo>] -k <hook>
  bitbucket-external-hook [options] list -b <bitbucket-uri> -p <project> [-r <repo>]
  bitbucket-external-hook -h | --help
  bitbucket-external-hook --version

Options:
  -b <bitbucket-uri>    URI to Bitbucket, can include auth info.
  -p <project>          Slug of project.
  -r <repository>       Slug of repository.
  -k <hook>             Hook key.           
  -h --help             Show this screen.
  -e --only-enabled     Show only enabled hooks.
  -c --only-configured  Show only configured hooks.
  --version             Show version.
`
)

type (
	Options struct {
		BitbucketURI   string `docopt:"-b"`
		Project        string `docopt:"-p"`
		Repository     string `docopt:"-r"`
		Hook           string `docopt:"-k"`
		OnlyEnabled    bool
		OnlyConfigured bool
		Print          bool
		List           bool
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

	}

	if err != nil {
		log.Fatal(err)
	}
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
	fmt.Printf("Executable: %v\n", settings.Exe)
	fmt.Printf("SafePath: %v\n", settings.SafePath)

	paramsPrefix := "Params: "
	fmt.Printf(
		"%v%v\n",
		paramsPrefix,
		strings.ReplaceAll(settings.Params, "\n", "\n"+strings.Repeat(" ", len(paramsPrefix))),
	)

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
