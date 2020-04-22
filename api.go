package main

import (
	"strconv"

	"github.com/bndr/gopencils"
	"github.com/kovetskiy/stash"
)

type API struct {
	*gopencils.Resource
	*Remote
}

func NewAPI(remote *Remote) *API {
	url := remote.scheme + "://" + remote.host + remote.path + "/rest/api/1.0"

	api := &API{
		Resource: gopencils.Api(
			url,
			&gopencils.BasicAuth{Username: remote.user, Password: remote.pass},
		),
	}

	return api
}

func (api *API) sub(project, repository string) *gopencils.Resource {
	resource := api.Resource
	if project != "" {
		resource = resource.Res("projects").Res(project)
	}
	if repository != "" {
		resource = resource.Res("repos").Res(repository)
	}
	return resource
}

func (api *API) GetHookSettings(project, repository, key string) (*HookSettings, error) {
	resource := api.sub(project, repository).Res("settings").
		Res("hooks").Res(key).
		Res("settings", &HookSettings{})

	response, err := request("GET", resource)
	if err != nil {
		return nil, err
	}

	return response.(*HookSettings), nil
}

func (api *API) StreamRepositories(process func(stash.Repository)) error {
	start := 0
	limit := 50
	more := true

	for more {
		reply, err := request(
			"GET",
			api.Res("repos", &stash.Repositories{}),
			map[string]string{
				"start": strconv.Itoa(start),
				"limit": strconv.Itoa(limit),
			},
		)
		if err != nil {
			return err
		}

		response := reply.(*stash.Repositories)
		for _, repo := range response.Repository {
			process(repo)
		}

		more = !response.IsLastPage
		start = response.NextPageStart
	}

	return nil
}

func (api *API) GetHooks(project, repository string) ([]*Hook, error) {
	response, err := request(
		"GET",
		api.sub(project, repository).Res("settings").
			Res("hooks", &ResponseHooks{}),
	)
	if err != nil {
		return nil, err
	}

	return response.(*ResponseHooks).Values, nil
}

func (api *API) GetHook(project, repository, key string) (*Hook, error) {
	rawResponse, err := request(
		"GET",
		api.sub(project, repository).Res("settings").
			Res("hooks").Res(key, &Hook{}),
	)
	if err != nil {
		return nil, err
	}

	response := rawResponse.(*Hook)
	return response, nil
}

func (api *API) EnableHook(project, repository, key string) error {
	_, err := request(
		"PUT",
		api.sub(project, repository).Res("settings").
			Res("hooks").Res(key).
			Res("enabled", &Hook{}),
	)

	return err
}

func (api *API) DisableHook(project, repository, key string) error {
	_, err := request(
		"DELETE",
		api.sub(project, repository).Res("settings").
			Res("hooks").Res(key).
			Res("enabled", &Hook{}),
	)

	return err
}

func (api *API) SetHookSettings(project, repository, key string, settings *HookSettings) error {
	_, err := request(
		"PUT",
		api.sub(project, repository).Res("settings").
			Res("hooks").Res(key).
			Res("settings", &HookSettings{}),
		settings,
	)

	return err
}
