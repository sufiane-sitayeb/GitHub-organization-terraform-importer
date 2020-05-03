package main

import (
    // "flag"
    "fmt"
    "os"
    "github_terraform_importer/lib/gh"
    "github_terraform_importer/internal/generalio"
)

func init() {
    // check if both env vars exist
    if !generalio.EnvExist("GITHUB_TOKEN") && !generalio.EnvExist("GITHUB_ORGANIZATION") {
        // If one of the two env vars is not set then exit
        fmt.Printf(" *** One of the following env vars has not been set: GITHUB_TOKEN or GITHUB_ORGANIZATION *** \n")
        os.Exit(1)
    }

    // Init github and Auth libraries
    gh.Initialize(os.Getenv("GITHUB_ORGANIZATION"), os.Getenv("GITHUB_TOKEN"))
}

func main() {
    // Import := flag.Bool("import", false, "If set, it also imports the terraform resources.")
    // AutoApprove := flag.Bool("auto-approve", false, "Skip interactive approval before importing.")
    // TerraformPath := flag.String("terraform_path", "./", "Absolute path where to create the terraform folder.")
    // BackendTemplatePath := flag.String("backend_template", "", "Absolute path of the template file you want to use.")
    //
    // flag.Parse()
    //
    // if *AutoApprove && *Import == false {
    //     fmt.Printf(" *** -auto-approve must have also -import set *** \n")
    //     os.Exit(1)
    // }

    // for _, u := range gh.Github.Users.All() {
    //     fmt.Printf("%-28s %s\n", u.GetLogin(), u.TerraformName())
    // }

    gh.Github.Users.Export()
    // gh.Github.Teams.Export()
    gh.Github.Repos.Export()
}
