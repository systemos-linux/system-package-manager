package database

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/schollz/progressbar/v3"
	"systemos.org/pkg/archive"
	"systemos.org/pkg/common"
	"systemos.org/pkg/semver"
)

const CONFIG_VALUE string = "https://github.com/systemos-linux/packages  offical  stable"

type Repository struct {
	url    string
	origin string
	branch string
}

type PackageNotFoundError struct {
	name    string
	version string
}

func (p *PackageNotFoundError) Error() string {
	return fmt.Sprintf("Package %s@%s not found in repository.", p.name, p.version)
}

func get_repository_config() (url, origin, branch string) {

	data, err := os.ReadFile("/etc/pkg/config")
	common.CheckIfError(err)
	if err != nil {
		fmt.Println("Package configuration not found, assuming defaults...")
		data = []byte(CONFIG_VALUE)
	}

	s := strings.Split(string(data), "  ")
	url = s[0]
	origin = s[1]
	branch = s[2]
	return
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
	common.CheckIfError(err)

	w, err := repo.Worktree()
	common.CheckIfError(err)

	bar := progressbar.DefaultBytes(
		-1,
		"Updating Database...",
	)

	err = w.Pull(&git.PullOptions{
		RemoteName: r.origin,
		Progress:   bar,
	})
	common.CheckIfError(err)

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(r.branch),
	})
	common.CheckIfError(err)

	// Print the latest commit that was just pulled
	ref, err := repo.Head()
	common.CheckIfError(err)

	commit, err := repo.CommitObject(ref.Hash())
	common.CheckIfError(err)

	v, err := repo.TagObject(commit.Hash)
	common.CheckIfError(err)

	fmt.Printf("System Database updated to %s", v.Name)

}

func (r *Repository) Get(package_name, version string) (*archive.Archive, error) {
	package_name = strings.Trim(strings.ToLower(package_name), " ")
	package_path := fmt.Sprintf("/var/pkg/database/%s/%s/pkg.json", package_name, version)
	_, err := os.Stat(package_path)
	if err != nil {
		return nil, &PackageNotFoundError{
			name:    package_name,
			version: version,
		}
	} else {
		return archive.Load(package_path), nil
	}
}

func (r *Repository) ListVersionsForPackage(package_name string) (versions []string, err error) {
	package_name = strings.Trim(strings.ToLower(package_name), " ")
	package_path := fmt.Sprintf("/var/pkg/database/%s/", package_name)
	folders, err := os.ReadDir(package_path)
	common.CheckIfError(err)
	if len(folders) > 0 {
		for _, folder := range folders {
			info, err := folder.Info()
			common.CheckIfError(err)
			versions = append(versions, info.Name())
		}
	}
	semver.Sort(versions)
	return
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
	common.CheckIfError(err)

	w, _ := r.Worktree()
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
	})
	common.CheckIfError(err)

}
