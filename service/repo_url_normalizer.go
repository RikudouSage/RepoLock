package service

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"go.chrastecky.dev/repolock/helper"
)

type RepositoryURLNormalizer interface {
	Normalize(uri string) (string, error)
}

func NewRepositoryURLNormalizer() RepositoryURLNormalizer {
	return &repositoryURLNormalizer{
		scpRegex: regexp.MustCompile("^(?:(?P<Username>[^@]+)@)?(?P<Host>[^:]+):(?P<Repo>.+)$"),
	}
}

type repositoryURLNormalizer struct {
	scpRegex *regexp.Regexp
}

func (receiver *repositoryURLNormalizer) Normalize(uri string) (string, error) {
	remote := strings.TrimSpace(uri)
	if remote == "" {
		return "", fmt.Errorf("invalid repository URL: %s", uri)
	}

	if strings.Contains(remote, "://") {
		parsed, err := url.Parse(remote)
		if err != nil {
			return "", fmt.Errorf("invalid repository URL: %s", parsed)
		}

		switch parsed.Scheme {
		case "http", "https", "ssh":
		default:
			return "", fmt.Errorf("unsupported repository scheme %q", parsed.Scheme)
		}

		host := parsed.Hostname()
		if host == "" {
			return "", fmt.Errorf("remote has no host")
		}

		return receiver.normalize(host, parsed.Path), nil
	}

	matches := helper.RegexNamedSubmatch(receiver.scpRegex, remote)
	if matches == nil {
		return "", fmt.Errorf("invalid repository URL: %s", remote)
	}

	return receiver.normalize(matches["Host"], matches["Repo"]), nil
}

func (*repositoryURLNormalizer) normalize(host, repo string) string {
	host = strings.ToLower(host)

	repo = strings.TrimPrefix(repo, "/")
	repo = strings.TrimSuffix(repo, "/")

	repo = strings.TrimSuffix(repo, ".git")

	return host + ":" + repo
}
