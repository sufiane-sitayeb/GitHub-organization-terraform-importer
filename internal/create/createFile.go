package create

import (
    "fmt"
    "github_terraform_importer/internal/generalio"
    "github_terraform_importer/internal/tfimport"
    "log"
    "os"
    "text/template"
)

// tfFile creates the terraform file using templates
func tfFile(resource string, templatePath string, destPath string, destinationFolder string, data interface{}) {

    // Get full destPath of the go file
    dir, err := os.Getwd()
    if err != nil {
        log.Fatal(err)
    }
    dir = dir + "/"

    // Get template file
    tmpl, err := template.New(resource + ".tmpl").ParseFiles(dir + templatePath)
    if err != nil {
        panic(err)
    }

    // Set output File
    terraformDir := destPath + "terraform/github/" + destinationFolder + "/"
    var repoFile *os.File

    // Create directory that will hold the tf files and set the tf file name used by the template
    if generalio.CreateDirIfNotExist(terraformDir) {
        repoFile, err = os.Create(terraformDir + resource + ".tf")
        if err != nil {
            log.Println("Create file: ", err)
            return
        }
    }

    // Execute template
    err = tmpl.Execute(repoFile, data)
    if err != nil {
        panic(err)
    }
}

// FilesAndImport creates the terraform files and imports them if needed
func FilesAndImport(resource string, githubOrg string, data interface{}, Import bool, AutoApprove bool, TerraformPath string, BackendTemplatePath string) {

    templatesPath := "configs/templates/" + resource + ".tmpl"

    // Generate terraform file
    tfFile(resource, templatesPath, TerraformPath, resource, data)

    if BackendTemplatePath != "" {
        tfFile("backend", BackendTemplatePath, TerraformPath, resource, "github-"+resource)
    }
    // Copy module files
    err := generalio.CopyDir("./configs/tf_files/modules/github/"+resource, TerraformPath+"./terraform/modules/"+resource)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    // Importing the data
    if Import {
        tfimport.Files(resource, data, githubOrg, TerraformPath, AutoApprove)
    }
}
