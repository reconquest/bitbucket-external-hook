package main

import (
	"fmt"
	"strings"

	"github.com/kovetskiy/stash"
	"github.com/reconquest/karma-go"
)

func handleCount(api *API, opts Options) error {
	log.Infof(nil, "receiving a total number of repositories")

	total := 0
	err := api.StreamRepositories(func(repo stash.Repository) {
		total++
	})
	if err != nil {
		return karma.Format(
			err,
			"unable to obtain total number of repositories",
		)
	}

	log.Infof(
		nil, "total number of repositories: %d", total,
	)

	log.Infof(
		nil,
		"started measuring total number of repositories with enabled/configured hooks",
	)

	totalEnabled := 0
	totalConfigured := 0
	tableEnabled := map[string]int{}
	tableConfigured := map[string]int{}

	index := 0
	err = api.StreamRepositories(func(repo stash.Repository) {
		index++

		hooks, err := api.GetHooks(repo.Project.Key, repo.Slug)
		if err != nil {
			log.Errorf(
				err,
				"unable to obtain hooks for repo: %s/%s",
				repo.Project.Key,
				repo.Name,
			)
		}

		for _, hook := range hooks {
			if strings.HasPrefix(hook.Details.Key, opts.Hook) {
				if hook.Enabled {

					var state string
					if hook.Configured {
						totalConfigured++

						if _, ok := tableConfigured[hook.Details.Key]; !ok {
							tableConfigured[hook.Details.Key] = 0
						} else {
							tableConfigured[hook.Details.Key]++
						}

						state = "enabled and configured"
					} else {
						totalEnabled++

						if _, ok := tableEnabled[hook.Details.Key]; !ok {
							tableEnabled[hook.Details.Key] = 0
						} else {
							tableEnabled[hook.Details.Key]++
						}

						state = "enabled (inherited from project)"
					}

					log.Debugf(
						nil,
						"%s/%s has %s hook: %s",
						repo.Project.Key,
						repo.Slug,
						hook.Details.Key,
						state,
					)
				}
			}
		}

		if index%50 == 0 || index == total {
			log.Infof(
				nil,
				"%.2f%% | %d of %d repositories processed",
				(float64(index) / float64(total) * 100),
				index,
				total,
			)
		}
	})
	if err != nil {
		return karma.Format(
			err,
			"unable to stream repositories",
		)
	}

	keys := []string{}
	for key, _ := range tableEnabled {
		keys = append(keys, key)
	}
	for key, _ := range tableConfigured {
		found := false
		for _, item := range keys {
			if key == item {
				break
			}
		}
		if !found {
			keys = append(keys, key)
		}
	}

	for i, key := range keys {
		if i != 0 {
			fmt.Println()
		}

		fmt.Printf("Hook: %s\n", key)

		configured, _ := tableConfigured[key]
		fmt.Printf("Enabled & Configured on %d repositories\n", configured)

		enabled, _ := tableEnabled[key]
		fmt.Printf("Enabled (inherited from project) on %d repositories\n", enabled)
	}

	return nil
}
