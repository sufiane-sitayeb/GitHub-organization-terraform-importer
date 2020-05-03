package generalio

import (
    "fmt"
    "io"
    "os"
)

// CopyFile copies files
func CopyFile(source string, dest string) (err error) {
    sourcefile, err := os.Open(source)
    if err != nil {
        return err
    }

    defer sourcefile.Close()

    destfile, err := os.Create(dest)
    if err != nil {
        return err
    }

    defer destfile.Close()

    _, err = io.Copy(destfile, sourcefile)
    if err == nil {
        sourceinfo, err := os.Stat(source)
        if err != nil {
            err = os.Chmod(dest, sourceinfo.Mode())
        }
    }

    return
}

// CopyDir copies directories recursively
func CopyDir(source string, dest string) (err error) {

    // get properties of source dir
    sourceinfo, err := os.Stat(source)
    if err != nil {
        return err
    }

    // create dest dir

    err = os.MkdirAll(dest, sourceinfo.Mode())
    if err != nil {
        return err
    }

    directory, _ := os.Open(source)

    objects, err := directory.Readdir(-1)

    for _, obj := range objects {

        sourcefilepointer := source + "/" + obj.Name()

        destinationfilepointer := dest + "/" + obj.Name()

        if obj.IsDir() {
            // create sub-directories - recursively
            err = CopyDir(sourcefilepointer, destinationfilepointer)
            if err != nil {
                fmt.Println(err)
            }
        } else {
            // perform copy
            err = CopyFile(sourcefilepointer, destinationfilepointer)
            if err != nil {
                fmt.Println(err)
            }
        }
    }
    return
}

// CreateDirIfNotExist checks if folder exists and if not, create it (used when we write templates)
func CreateDirIfNotExist(dir string) bool {
    if _, err := os.Stat(dir); os.IsNotExist(err) {
        err = os.MkdirAll(dir, 0755)
        if err != nil {
            panic(err)
        }
    }
    return true
}
