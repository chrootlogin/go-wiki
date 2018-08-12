package repo

import (
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// Original author: https://github.com/orirawlings, https://github.com/src-d/go-git/issues/562#issuecomment-324819816

type memo map[plumbing.Hash]plumbing.Hash

// filterByChangesToPath provides a CommitIter callback that only invokes
// a delegate callback for commits that include changes to the content of path.
func filterByChangesToPath(r *git.Repository, path string, callback func(*object.Commit) error) func(*object.Commit) error {
	m := make(memo)
	return func(c *object.Commit) error {
		if err := ensure(m, c, path); err != nil {
			return err
		}
		if c.NumParents() == 0 && !m[c.Hash].IsZero() {
			// c is a root commit containing the path
			return callback(c)
		}
		// Compare the path in c with the path in each of its parents
		for _, p := range c.ParentHashes {
			if _, ok := m[p]; !ok {
				pc, err := r.CommitObject(p)
				if err != nil {
					return err
				}
				if err := ensure(m, pc, path); err != nil {
					return err
				}
			}
			if m[p] != m[c.Hash] {
				// contents at path are different from parent
				return callback(c)
			}
		}
		return nil
	}
}

// ensure our memoization includes a mapping from commit hash
// to the hash of path contents.
func ensure(m memo, c *object.Commit, path string) error {
	if _, ok := m[c.Hash]; !ok {
		t, err := c.Tree()
		if err != nil {
			return err
		}
		te, err := t.FindEntry(path)
		if err == object.ErrDirectoryNotFound {
			m[c.Hash] = plumbing.ZeroHash
			return nil
		} else if err != nil {
			if !strings.ContainsRune(path, '/') {
				// path is in root directory of project, but not found in this commit
				m[c.Hash] = plumbing.ZeroHash
				return nil
			}
			return err
		}
		m[c.Hash] = te.Hash
	}
	return nil
}