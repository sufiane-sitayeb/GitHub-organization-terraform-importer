package gh

import (
    "fmt"
    "os"
    "strings"
    "github.com/google/go-github/v29/github"
    "github_terraform_importer/lib/util"
)

type TeamMember struct {
    User        *User
    Membership  *github.Membership
}

type Team struct {
    *github.Team
    Members     []TeamMember
}

func (t Team) TerraformName() string {
    return strings.Replace(t.GetName(), " ", "_", -1)
}

func (t Team) MembersString() string {
    var members []string
    for _,m := range t.Members {
        members = append(members, m.User.GetLogin())
    }
    return util.Serialize(members)
}

func api_fetch_teams() (ret []Team) {
    for {
        opt := &github.ListOptions{
            PerPage: 5,
        }

        teams, resp, err := Github.Client.Teams.ListTeams(Github.Ctx, Github.Org, opt)
        if err != nil {
            fmt.Printf("ERROR: %s\n", err)
            os.Exit(1)
        }
        for _,t := range teams {
            ret = append(ret, Team {t, []TeamMember{}})
        }
        if resp.NextPage == 0 {
            break
        }
        opt.Page = resp.NextPage
break
    }

    // Get members of each team
    for i, _ := range ret {
        team := &ret[i]

        opt := &github.TeamListTeamMembersOptions{
            ListOptions: github.ListOptions{PerPage: 10},
        }

        for {
            members, resp, err := Github.Client.Teams.ListTeamMembers(Github.Ctx, team.GetID(), opt)

            if err != nil {
                fmt.Printf("ERROR: %s\n", err)
                os.Exit(1)
            }

            for _,m := range members {
                membership,_,_ := Github.Client.Teams.GetTeamMembership(Github.Ctx, team.GetID(), m.GetLogin())
                ret[i].Members = append(ret[i].Members, TeamMember {&User{m},membership})
            }

            if resp.NextPage == 0 {
                break
            }

            opt.Page = resp.NextPage
break
        }
    }

    return ret
}

type TeamCollection struct {
    all []Team
}

func (c TeamCollection) All() []Team {
    if c.all == nil {
        c.all = api_fetch_teams()
    }
    return c.all
}

func (c TeamCollection) Export() {
    util.TemplateToFile("configs/templates/teams.tmpl", "terraform/teams.tf", c.All())
}
