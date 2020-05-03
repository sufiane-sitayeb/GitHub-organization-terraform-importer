package tfimport

import (
    "fmt"
    "github_terraform_importer/internal/fetch"
    "github_terraform_importer/internal/generalio"
    "os/exec"
    "strconv"
)

// Files imports the terraform files just created
func Files(resource string, data interface{}, githubOrg string, path string, autoApprove bool) {
    if generalio.InputFromUser(resource, autoApprove) {
        fmt.Printf("\nContinuing and importing %s\n", resource)
        // Terraform init the path
        cmd := exec.Command("terraform", "init")
        cmd.Dir = path + "terraform/github/" + resource
        out, err := cmd.Output()
        outSring := string(out[:])
        if err != nil {
            panic(err)
        }
        fmt.Print(outSring)

        switch v := data.(type) {
        case []fetch.ReposT:
            fmt.Printf(" *** \n\nImporting each %s one by one...it may take a while...\n\n *** \n\n", resource)
            for _, repo := range v {
                cmd := exec.Command("terraform", "import", "module.github_repo_"+repo.GithubModuleName+".github_repository.repo", repo.Name)
                cmd.Dir = path + "/terraform/github/" + resource
                out, err := cmd.Output()
                outSring := string(out[:])
                if err != nil {
                    fmt.Printf("Unable to import %s: %s due to %s\n", resource, repo.GithubModuleName, err)
                }
                fmt.Print(outSring)
            }
        case []fetch.UsersT:
            fmt.Printf(" *** \n\nImporting each %s one by one...it may take a while...\n\n *** \n\n", resource)
            for _, user := range v {
                cmd := exec.Command("terraform", "import", "module.github_user_"+user.GithubModuleName+".github_membership.user", githubOrg+":"+user.Login)
                cmd.Dir = path + "terraform/github/" + resource
                out, err := cmd.Output()
                outSring := string(out[:])
                if err != nil {
                    fmt.Printf("Unable to import %s: %s due to %s\n", resource, user.GithubModuleName, err)
                }
                fmt.Print(outSring)
            }
        case []fetch.TeamsT:
            fmt.Printf(" *** \n\nImporting each %s one by one...it may take a while...\n\n *** \n\n", resource)
            for _, team := range v {

                // Converting the base64 ID into string
                TeamIDStr := strconv.FormatInt(team.ID, 10)

                cmd := exec.Command("terraform", "import", "module.github_team_"+team.GithubModuleName+".github_team.team", TeamIDStr)
                cmd.Dir = path + "terraform/github/" + resource
                out, _ := cmd.Output()
                outSring := string(out[:])
                if err != nil {
                    fmt.Printf("Unable to import %s: %s due to %s\n", resource, team.GithubModuleName, err)
                }
                fmt.Print(outSring)

                for user := range team.MemberRoleMap {
                    cmd := exec.Command("terraform", "import", "module.github_team_"+team.GithubModuleName+".github_team_membership.team_membership["+"\""+user+"\""+"]", TeamIDStr+":"+user)
                    cmd.Dir = path + "terraform/github/" + resource
                    out, err := cmd.Output()
                    outSring := string(out[:])
                    if err != nil {
                        fmt.Printf("Unable to import %s: %s due to %s\n", resource, team.GithubModuleName, err)
                    }
                    fmt.Print(outSring)

                }
                fmt.Print(outSring)
            }
        }
    }
}
