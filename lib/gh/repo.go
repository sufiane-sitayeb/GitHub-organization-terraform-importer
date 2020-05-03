package gh

import (
    "fmt"
    "os"
    "strings"
    "github.com/google/go-github/v29/github"
    "github_terraform_importer/lib/util"
)

type Repo struct {
    *github.Repository
}

func (r Repo) TerraformName() string {
    return strings.Replace(r.GetName(), " ", "_", -1)
}

func (r Repo) TopicsString() string {
    return util.Serialize(r.Topics)
}

func api_fetch_repos() (ret []Repo) {
    opt := &github.RepositoryListByOrgOptions{
        ListOptions: github.ListOptions{PerPage: 20},
    }

    // Retrieve all repos
    for {
        repos, resp, err := Github.Client.Repositories.ListByOrg(Github.Ctx, Github.Org, opt)
        if err != nil {
            fmt.Printf("ERROR: %s\n", err)
            os.Exit(1)
        }
        for _, r := range repos {
            ret = append(ret, Repo {r})
        }
        if resp.NextPage == 0 {
            break
        }
        opt.Page = resp.NextPage
break
    }

    return ret
}

type RepoCollection struct {
    all []Repo
}

func (c RepoCollection) All() []Repo {
    if c.all == nil {
        c.all = api_fetch_repos()
    }
    return c.all
}

func (c RepoCollection) Export() {
    util.TemplateToFile("configs/templates/repos.tmpl", "terraform/repos.tf", c.All())
}
