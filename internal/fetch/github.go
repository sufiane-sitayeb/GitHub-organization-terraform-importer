package fetch

import (
    "context"
    "fmt"
    "os"
    "strconv"
    "strings"

    "github.com/google/go-github/v29/github"
    "golang.org/x/oauth2"
)

// ReposT This struct hold only the informations we need for the repos
type ReposT struct {
    GithubModuleName    string
    Name                string
    Description         string
    Homepage            string
    Private             bool
    HasIssues           bool
    HasProjects         bool
    HasWiki             bool
    AllowMergeCommit    bool
    AllowSquashMerge    bool
    AllowRebaseMerge    bool
    HasDownloads        bool
    DeleteBranchOnMerge bool
    DefaultBranch       string
    Archived            bool
    Topics              string
}

// TeamsT This struct hold only the informations we need for the users
type TeamsT struct {
    GithubModuleName string
    ID               int64
    Name             string
    Description      string
    Privacy          string
    MembersString    string
    MemberRoleMap    map[string]string
}

// InitLib initializes both auth and github library.
func InitLib() (githubOrg string, ctx context.Context, client *github.Client) {

    githubToken := os.Getenv("GITHUB_TOKEN")
    githubOrg = os.Getenv("GITHUB_ORGANIZATION")

    ctx = context.Background()
    ts := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: githubToken},
    )
    tc := oauth2.NewClient(ctx, ts)

    client = github.NewClient(tc)

    return
}

// Teams gets all users from the organization
func Teams(ctx context.Context, githubOrg string, client *github.Client) (teams []TeamsT) {

    var allTeams []*github.Team

    // Get teams
    for {

        // Set Pagination
        opt := &github.ListOptions{
            PerPage: 20,
        }

        teamPage, resp, err := client.Teams.ListTeams(ctx, githubOrg, opt)
        if err != nil {
            fmt.Printf("ERROR: %s\n", err)
            os.Exit(1)
        }
        allTeams = append(allTeams, teamPage...)
        if resp.NextPage == 0 {
            break
        }
        opt.Page = resp.NextPage
    }

    // Get members of each team
    for _, team := range allTeams {

        var memberRoleString string
        var memberRoleMap = make(map[string]string)

        groupMemberOptions := [2]string{"member", "maintainer"}
        groupMemberByRole := make(map[string][]*github.User)

        for _, role := range groupMemberOptions {
            // Set Pagination
            opt := &github.TeamListTeamMembersOptions{
                Role:        role,
                ListOptions: github.ListOptions{PerPage: 20},
            }

            for {
                teamMembersPage, resp, err := client.Teams.ListTeamMembers(ctx, team.GetID(), opt)
                if err != nil {
                    fmt.Printf("ERROR: %s\n", err)
                    os.Exit(1)
                }

                groupMemberByRole[role] = append(groupMemberByRole[role], teamMembersPage...)
                if resp.NextPage == 0 {
                    break
                }
                opt.Page = resp.NextPage
            }

            for _, user := range groupMemberByRole[role] {
                memberRoleString += strconv.Quote(*user.Login) + " : " + strconv.Quote(role) + ", "
                memberRoleMap[*user.Login] = role
            }
        }
        teams = append(teams, TeamsT{
            // Replacing all spaces in the team name with _. This is used as terraform module name
            GithubModuleName: strings.Replace(team.GetName(), " ", "_", -1),
            ID:               team.GetID(),
            Name:             team.GetName(),
            // This escapes every quote from the variable we use for the description
            Description:   strings.Replace(team.GetDescription(), "\"", "\\\"", -1),
            Privacy:       team.GetPrivacy(),
            MembersString: memberRoleString,
            MemberRoleMap: memberRoleMap,
        })
    }
    return teams
}

// Repos returns all repositories from a given organization in a list of struct with only the parameters needed and used by the template.
func Repos(ctx context.Context, githubOrg string, client *github.Client) (repos []ReposT) {

    // Manage GitHub pagination api
    opt := &github.RepositoryListByOrgOptions{
        ListOptions: github.ListOptions{PerPage: 50},
    }

    var allRepos []*github.Repository

    // Retrieve all repos
    for {
        repos, resp, err := client.Repositories.ListByOrg(ctx, githubOrg, opt)
        if err != nil {
            fmt.Printf("ERROR: %s\n", err)
            os.Exit(1)
        }
        allRepos = append(allRepos, repos...)
        if resp.NextPage == 0 {
            break
        }
        opt.Page = resp.NextPage
    }

    // Extract from allRepos only the informations we need
    for _, repo := range allRepos {

        // I have to do this because the ListByOrg doesn't return any has* (has_issues, has_projects etc) values
        // https://developer.github.com/v3/repos/#get-a-repository
        repoDetails, _, err := client.Repositories.Get(ctx, githubOrg, repo.GetName())
        if err != nil {
            fmt.Printf("ERROR: %s\n", err)
            os.Exit(1)
        }

        var topics []string
        // Add quotes so we can printe them for each repo in the list
        for _, topic := range repoDetails.Topics {
            topics = append(topics, strconv.Quote(topic))
        }

        // Convert the slice in a string with comma in between strings
        // ie: we end up with something like this "csharp","extensions","grpc"
        topicsString := strings.Join(topics, ",")

        // Create repos that holds only the informations we need on the template.
        repos = append(repos, ReposT{
            // This removes all dots from the variable we use for the module name
            GithubModuleName: strings.Replace(repoDetails.GetName(), ".", "-", -1),
            Name:             repoDetails.GetName(),
            // This escapes every quote from the variable we use for the description
            Description:         strings.Replace(repoDetails.GetDescription(), "\"", "\\\"", -1),
            Homepage:            repoDetails.GetHomepage(),
            Private:             repoDetails.GetPrivate(),
            HasIssues:           repoDetails.GetHasIssues(),
            HasProjects:         repoDetails.GetHasProjects(),
            HasWiki:             repoDetails.GetHasWiki(),
            AllowMergeCommit:    repoDetails.GetAllowMergeCommit(),
            AllowSquashMerge:    repoDetails.GetAllowSquashMerge(),
            AllowRebaseMerge:    repoDetails.GetAllowRebaseMerge(),
            HasDownloads:        repoDetails.GetHasDownloads(),
            DeleteBranchOnMerge: repoDetails.GetDeleteBranchOnMerge(),
            DefaultBranch:       repoDetails.GetDefaultBranch(),
            Archived:            repoDetails.GetArchived(),
            Topics:              topicsString,
        })
        // fmt.Printf("repo: %s Merge, commit: %t\n", *repoDetails.Name, repoDetails.GetAllowMergeCommit())
    }
    return repos
}
