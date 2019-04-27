package main

import (
	"net/url"
	"regexp"
)

var (
	reParserWebURL = regexp.MustCompile(
		`://([^/]+)/((users|projects)/([^/]+))/repos/([^/]+)`,
	)
)

type Remote struct {
	scheme      string
	host        string
	project     string
	projectType string
	repo        string
	user        string
	pass        string
}

func GetRemote(opts Options) (*Remote, error) {
	uri, err := url.Parse(opts.BitbucketURI)
	if err != nil {
		return nil, err
	}

	var user string
	var pass string
	if uri.User != nil {
		user = uri.User.Username()
		pass, _ = uri.User.Password()
	}

	remote := &Remote{
		scheme:      uri.Scheme,
		host:        uri.Host,
		project:     opts.Project,
		projectType: "projects",
		repo:        opts.Repository,
		user:        user,
		pass:        pass,
	}

	return remote, nil
}

func getRepoWebURL(host, project, projectType, repo string) string {
	return "http://" + host + "/" +
		projectType + "/" + project + "/repos/" + repo
}
