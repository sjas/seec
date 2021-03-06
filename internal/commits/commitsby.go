package commits

import (
	"fmt"

	"seec/internal/gh"
	"seec/internal/users"

	"github.com/VonC/godbg"
)

var Pdbg *godbg.Pdbg

type CommitsByAuthor struct {
	author string
	cbd    []*CommitsByDate
}
type CommitsByAuthors map[users.User]*CommitsByAuthor

// Because of seec 709cd912d4663af87903d3d278a3bab9d4d84153
type CommitsByDate struct {
	date    string
	commits []*gh.Commit
}

func (cba *CommitsByAuthor) String() string {
	res := ""
	first := true
	for i, acbd := range cba.cbd {
		if !first {
			res = res + ", "
		}
		first = false
		if i == len(cba.cbd)-1 && i > 0 {
			res = res + "and "
		}
		res = res + acbd.String()
	}
	return fmt.Sprintf("%s=>%s", cba.author, res)
}

func (cbd *CommitsByDate) String() string {
	res := ""
	first := true
	for i, commit := range cbd.commits {
		if !first {
			res = res + ", "
		}
		first = false
		if i == len(cbd.commits)-1 && i > 0 {
			res = res + "and "
		}
		res = res + (*commit.SHA)[:7]
	}
	return fmt.Sprintf("%s (%s)", res, cbd.date)
}

func NewCommitsByAuthor(authorname string) *CommitsByAuthor {
	return &CommitsByAuthor{authorname, []*CommitsByDate{}}
}

func NewCommitsByDate(commit *gh.Commit) *CommitsByDate {
	return &CommitsByDate{commit.AuthorDate(), []*gh.Commit{commit}}
}

func (cbas CommitsByAuthors) Add(somecbas CommitsByAuthors) {
	for authorName, pcommitsByAuthor := range somecbas {
		acommitsByAuthor := cbas[authorName]
		if acommitsByAuthor == nil {
			cbas[authorName] = pcommitsByAuthor
		} else {
			acommitsByAuthor.addCba(pcommitsByAuthor)
			// cbas[authorName] = acommitsByAuthor
			Pdbg.Pdbgf("Put commits '%s' for author '%s'", acommitsByAuthor.String(), authorName)
		}
	}
	Pdbg.Pdbgf("ADD RES '%s'", cbas)
}

func (cba *CommitsByAuthor) addCba(acba *CommitsByAuthor) {
	Pdbg.Pdbgf("addCba '%s' to cba '%s'", acba, cba)
	for _, acbd := range acba.cbd {
		date := acbd.date
		found := false
		for _, cbd := range cba.cbd {
			Pdbg.Pdbgf("addCba cbd '%s' vs acba '%s'", cbd.date, date)
			if cbd.date == date {
				found = true
				cbd.commits = append(cbd.commits, acbd.commits...)
				break
			}
		}
		if !found {
			Pdbg.Pdbgf("addCba date not found => add '%s' to '%s'", acbd, cba)
			cba.cbd = append(cba.cbd, acbd)
		}
	}
}

func (cba *CommitsByAuthor) AddCommit(commit *gh.Commit) {
	date := commit.AuthorDate()
	for _, cbd := range cba.cbd {
		Pdbg.Pdbgf("ADDCOMMIT: cbd date '%s' vs commit date '%s'", cbd.date, date)
		if cbd.date == date {
			cbd.commits = append(cbd.commits, commit)
			return
		}
	}
	cba.cbd = append(cba.cbd, NewCommitsByDate(commit))
}

func (cba *CommitsByAuthor) CommitsByDate() []*CommitsByDate {
	return cba.cbd
}

func (cbd *CommitsByDate) Commits() []*gh.Commit {
	return cbd.commits
}

func (cbd *CommitsByDate) Date() string {
	return cbd.date
}
