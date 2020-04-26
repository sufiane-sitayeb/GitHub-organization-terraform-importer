package main

import (
	"flag"
	"fmt"
	"os"

	"github_terraform_importer/internal/create"
	"github_terraform_importer/internal/fetch"
	"github_terraform_importer/internal/generalio"
)

func init() {
	// check if both env vars exist
	if !generalio.EnvExist("githubToken") && !generalio.EnvExist("githubOrgANIZATION") {
		// If one of the two env vars is not set then exit
		fmt.Printf(" *** One of the following env vars has not been set: githubToken or githubOrgANIZATION *** \n")
		os.Exit(1)
	}
}

func main() {

	// Init github and Auth libraries
	githubOrg, ctx, client := fetch.InitLib()

	Import := flag.Bool("import", false, "If set, it also imports the terraform resources.")
	AutoApprove := flag.Bool("auto-approve", false, "Skip interactive approval before importing.")
	TerraformPath := flag.String("terraform_path", "./", "Absolute path where to create the terraform folder.")
	BackendTemplatePath := flag.String("backend_template", "", "Absolute path of the template file you want to use.")

	flag.Parse()

	if *AutoApprove && *Import == false {
		fmt.Printf(" *** -auto-approve must have also -import set *** \n")
		os.Exit(1)
	}

	//  This maps the resource with the correspondent function for retrieving the data
	type fnT = func() interface{}
	resourcesToFunction := map[string]fnT{
		"users": func() interface{} { return fetch.Users(githubOrg, ctx, client) },
		"teams": func() interface{} { return fetch.Teams(githubOrg, ctx, client) },
		"repos": func() interface{} { return fetch.Repos(githubOrg, ctx, client) },
	}

	//  Generate the terraform files and import
	for resource, fn := range resourcesToFunction {

		fmt.Printf("Generating Terraform %s file...\n", resource)
		data := fn()

		create.FilesAndImport(resource, githubOrg, data, *Import, *AutoApprove, *TerraformPath, *BackendTemplatePath)
	}
}
