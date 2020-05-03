package gh

import (
    "fmt"
    "os"
    "strings"
    "github.com/google/go-github/v29/github"
    "github_terraform_importer/lib/util"
)

type User struct {
    *github.User
}

func (u User) TerraformName() string {
    return strings.Replace(u.GetLogin(), " ", "_", -1)
}

func (u User) Role() string {
    if u.GetSiteAdmin() {
        return "admin"
    } else {
        return "member"
    }
}

func api_fetch_users() (ret []User) {
    opt := &github.ListMembersOptions{
        ListOptions: github.ListOptions{PerPage: 20},
    }

    for {
        users, resp, err := Github.Client.Organizations.ListMembers(Github.Ctx, Github.Org, opt)
        if err != nil {
            fmt.Printf("ERROR: %s\n", err)
            os.Exit(1)
        }

        for _, u := range users {
            ret = append(ret, User {u})
        }

        if resp.NextPage == 0 {
            break
        }

        opt.Page = resp.NextPage
break
    }

    return ret
}

type UserCollection struct {
    all []User
}

func (c UserCollection) All() []User {
    if c.all == nil {
        c.all = api_fetch_users()
    }
    return c.all
}

func (c UserCollection) Export() {
    util.TemplateToFile("configs/templates/users.tmpl", "terraform/users.tf", c.All())
}
