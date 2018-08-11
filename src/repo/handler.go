package repo

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

)

func getCommitHistoryOfFile(path string) []*object.Commit {
	objects := []*object.Commit{}

	ref, err := repo.Head()
	if err == nil {
		cIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
		if err == nil {
			cIter.ForEach(filterByChangesToPath(repo, path, func(c *object.Commit) error {
				objects = append(objects, c)
				return nil
			}))
		}
	}

	return objects
}