package util

import (
    "fmt"
    "path"
    "log"
    "os"
    "encoding/json"
    "text/template"
    "github_terraform_importer/internal/generalio"
)

func Serialize(a interface{}) string {
    byteArray, err := json.Marshal(a)
    if err != nil { fmt.Println(err) }
    return string(byteArray)
}

func TemplateToFile(path_tpl string, path_out string, data interface{}) {
    dir_out := path.Dir(path_out)
    template_name := path.Base(path_tpl)
    tmpl, err := template.New(template_name).ParseFiles(path_tpl)
    if err != nil { panic(err) }

    // Create directory that will hold the tf files and set the tf file name used by the template
    if generalio.CreateDirIfNotExist(dir_out) {
        if err != nil {
            log.Println("Create file: ", err)
            return
        }
    }

    fh, err := os.Create(path_out)

    // Execute template
    err = tmpl.Execute(fh, data)
    if err != nil { panic(err) }
}
