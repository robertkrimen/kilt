package main

import (
	"fmt"
	"github.com/robertkrimen/kilt"
	"go/build"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
)

const (
	otherName = "kilt"
	otherGo   = "kilt.go"
	otherURL  = "http://raw.github.com/robertkrimen/kilt/master/kilt.go"
)

func main() {
	thisPath, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		log.Fatal(err)
	}

	thisImport := func() string {
		for _, base := range build.Default.SrcDirs() {
			cleanBase := filepath.Clean(base)
			if cleanBase != "/" {
				cleanBase += "/"
			}
			if strings.HasPrefix(thisPath, cleanBase) {
				return thisPath[len(cleanBase):]
			}
		}
		return ""
	}()

	// TODO Add option to force... how do we detect the current package name?
	if thisImport == "" {
		log.Fatal("No Go package detected in the current directory (not in $GOROOT or $GOPATH)")
	}

	thisImport = thisImport                          // .../<this>
	otherName := otherName                           // kilt
	otherPath := filepath.Join(thisPath, otherName)  // ./kilt
	otherGoPath := filepath.Join(otherPath, otherGo) // ./kilt/kilt.go

	{
		fmt.Fprintf(os.Stderr, kilt.StringTrimSpace(`
---
The following file will be written (or over-written) (from "%s"):

    ./kilt/kilt.go      %s

        `), otherURL, otherGoPath)

		fmt.Fprintf(os.Stderr, "\n---\nContinue? (yN) ")
		yesNo := ""
		fmt.Scanln(&yesNo)
		if !regexp.MustCompile(`^\s*[yY]\s*$`).MatchString(yesNo) {
			fmt.Fprintf(os.Stderr, "Cancel.\n")
			os.Exit(0)
		}
	}

	// <package>/kilt
	{
		path := otherPath
		meta, err := os.Stat(path)
		if err == nil {
			if !meta.IsDir() {
				log.Fatalf("%s: is not a directory", path)
			}
		} else if os.IsNotExist(err) {
			err := os.Mkdir(path, 0777)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}

	// <package>/kilt/kilt.go
	{
		path := otherGoPath
		meta, err := os.Stat(path)
		if err == nil {
			if meta.IsDir() {
				log.Fatalf("%s: is a directory", path)
			}
		} else if os.IsNotExist(err) {
		} else {
			log.Fatal(err)
		}

		tmp, err := ioutil.TempFile("", "other.go.")
		fatal := func(err error) {
			os.Remove(tmp.Name())
			log.Fatal(err)
		}
		if err != nil {
			fatal(err)
		}

		response, err := http.Get(otherURL)
		if err != nil {
			log.Fatal(err)
		} else if response.StatusCode != 200 {
			log.Fatalf("%s: %s", otherURL, response.Status)
		}
		defer response.Body.Close()
		_, err = io.Copy(tmp, response.Body)
		if err != nil {
			fatal(err)
		}

		err = os.Rename(tmp.Name(), otherGoPath)
		if err != nil {
			fatal(err)
		}

		// Kind of a hack
		// Might be a race condition, but we're not goroutining, so...
		umask := syscall.Umask(0)
		syscall.Umask(umask)
		os.Chmod(otherGoPath, os.FileMode(0666-umask))
	}
}
