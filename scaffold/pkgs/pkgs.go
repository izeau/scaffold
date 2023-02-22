// Package pkgs contains functions for parsing remote urls and checking if a
// directory is a git repository.
//
// TODO: pkgs name should be changed to something more descriptive
package pkgs

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

// ParseRemote parses a URL and returns a filesystem path representing the
// repository.
//
// Examples:
//
//	ParseRemote("https://github.com/hay-kot/scaffold-go-cli")
//	github.com
//	└── hay-kot
//	    └── scaffold-go-cli
//			└── repository files
func ParseRemote(urlStr string) (string, error) {
	url, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse url: %w", err)
	}

	host := url.Host
	split := strings.Split(url.Path, "/")

	// Remove .git from repo name if it exists but keeps @tag or @branch intact
	split[len(split)-1] = strings.Replace(split[len(split)-1], ".git", "", 1)

	if len(split) < 3 {
		return "", fmt.Errorf("invalid url")
	}

	user := split[1]
	repo := split[2]

	return filepath.Join(host, user, repo), nil
}

// IsRemote checks if the string is a remote url or an alias for a remote url
// if it is a remote url, it returns the url. If the string uses and alias
// it returns the expanded url
//
// Examples:
//
//	isRemote(gh:foo/bar) -> https://github.com/foo/bar, true
func IsRemote(str string, shorts map[string]string) (expanded string, ok bool) {
	if strings.HasPrefix(str, "http") {
		return str, true
	}

	split := strings.Split(str, ":")

	if len(split) == 2 {
		short := split[0]

		for k, v := range shorts {
			if k == short {
				out, err := url.JoinPath(v, split[1])
				if err != nil {
					return "", false
				}

				return out, true
			}
		}
	}

	return "", false
}

// Update updates a git repository to the latest commit
func Update(path string) (updated bool, err error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return false, err
	}

	w, err := repo.Worktree()
	if err != nil {
		return false, err
	}

	err = w.Pull(&git.PullOptions{})
	if err != nil {
		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
