package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/schollz/progressbar/v3"
)

const CONFIG_VALUE string = "https://github.com/systemos-linux/packages  offical  stable"

type Repository struct {
	url    string
	origin string
	branch string
}

func get_repository_config() (url, origin, branch string) {

	data, err := os.ReadFile("/etc/pkg/config")
	CheckIfError(err)

	if len(data) == 0 {
		data = []byte(CONFIG_VALUE)
	}
	s := strings.Split(string(data), "  ")
	url = s[0]
	origin = s[1]
	branch = s[2]
}

func new_repository() *Repository {
	u, o, b := get_repository_config()

	return &Repository{
		url:    u,
		origin: o,
		branch: b,
	}
}

func (r *Repository) Init() {
	_, err := git.PlainOpen("/var/pkg/database")
	if err == git.ErrRepositoryNotExists {
		clone_repository(r.url, r.origin, r.branch)
	}
}

func (r *Repository) Update() {

	repo, err := git.PlainOpen("/var/pkg/database")
	CheckIfError(err)

	w, err := repo.Worktree()
	CheckIfError(err)

	bar := progressbar.DefaultBytes(
		-1,
		"Updating Database...",
	)

	err = w.Pull(&git.PullOptions{
		RemoteName: r.origin,
		Progress:   bar,
	})
	CheckIfError(err)

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(r.branch),
	})
	CheckIfError(err)

	// Print the latest commit that was just pulled
	ref, err := repo.Head()
	CheckIfError(err)

	commit, err := repo.CommitObject(ref.Hash())
	CheckIfError(err)

	v, err := repo.TagObject(commit.Hash)
	CheckIfError(err)

	fmt.Printf("System Database updated to %s", v.Name)

}

func clone_repository(url, origin, branch string) {

	bar := progressbar.DefaultBytes(
		-1,
		"Installing Database...",
	)

	r, err := git.PlainClone("/var/pkg/database", false, &git.CloneOptions{
		RemoteName: origin,
		URL:        url,
		Progress:   bar,
		Tags:       git.AllTags,
	})
	CheckIfError(err)

	w, _ := r.Worktree()
	CheckIfError(w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
	}))

}
