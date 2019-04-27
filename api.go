package main

import "github.com/bndr/gopencils"

type API struct {
	*gopencils.Resource
	*Remote
}

func NewAPI(remote *Remote) *API {
	url := remote.scheme + "://" + remote.host + "/rest/api/1.0" +
		"/" + remote.projectType + "/" + remote.project

	if remote.repo != "" {
		url += "/repos/" + remote.repo
	}

	api := &API{
		Resource: gopencils.Api(
			url,
			&gopencils.BasicAuth{Username: remote.user, Password: remote.pass},
		),
		Remote: remote,
	}

	return api
}

func (api *API) GetHookSettings(key string) (*HookSettings, error) {
	response, err := request(
		"GET",
		api.Res("settings").
			Res("hooks").Res(key).
			Res("settings", &HookSettings{}),
	)
	if err != nil {
		return nil, err
	}

	return response.(*HookSettings), nil
}

func (api *API) GetHooks() ([]*Hook, error) {
	response, err := request(
		"GET",
		api.Res("settings").
			Res("hooks", &ResponseHooks{}),
	)
	if err != nil {
		return nil, err
	}

	return response.(*ResponseHooks).Values, nil
}

func (api *API) GetHook(key string) (*Hook, error) {
	rawResponse, err := request(
		"GET",
		api.Res("settings").
			Res("hooks").Res(key, &Hook{}),
	)
	if err != nil {
		return nil, err
	}

	response := rawResponse.(*Hook)
	return response, nil
}

func (api *API) EnableHook(key string) error {
	_, err := request(
		"PUT",
		api.Res("settings").
			Res("hooks").Res(key).
			Res("enabled", &Hook{}),
	)

	return err
}

func (api *API) DisableHook(key string) error {
	_, err := request(
		"DELETE",
		api.Res("settings").
			Res("hooks").Res(key).
			Res("enabled", &Hook{}),
	)

	return err
}

func (api *API) SetHookSettings(key string, settings *HookSettings) error {
	_, err := request(
		"PUT",
		api.Res("settings").
			Res("hooks").Res(key).
			Res("settings", &HookSettings{}),
		settings,
	)

	return err
}
