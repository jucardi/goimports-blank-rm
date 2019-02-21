package main

import (
	"github.com/jucardi/go-logger-lib/log"
	"github.com/jucardi/go-osx/paths"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var regexImports = regexp.MustCompile(`(import \()(\s|.)+?(\))`)

func main() {
	if len(os.Args) < 2 {
		println("invalid argument")
	}
	processDir(os.Args[1])
}

func processDir(dir string) {
	var directories []os.FileInfo

	objs, err := ioutil.ReadDir(dir)

	if err != nil {
		log.FatalErr(err)
	}

	for _, f := range objs {
		if f.IsDir() {
			directories = append(directories, f)
			continue
		}
		extSplit := strings.Split(f.Name(), ".")
		if strings.ToLower(extSplit[len(extSplit)-1]) == "go" && strings.ToLower(extSplit[len(extSplit)-2]) != "pb"{
			processFile(paths.Combine(dir, f.Name()))
		}
	}

	for _, d := range directories {
		if strings.TrimSpace(d.Name()) == "" {
			continue
		}
		dirName := paths.Combine(dir, d.Name())
		processDir(dirName)
	}
}

func processFile(file string) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.FatalErr(err)
	}
	log.Debug("file: ", file)
	imports := regexImports.FindString(string(data))
	lines := strings.Split(imports, "\n")
	var newImports []string
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		newImports = append(newImports, line)
	}
	str := strings.Join(newImports, "\n")
	if str == imports {
		return
	}
	log.Info(file)
	result := regexImports.ReplaceAllString(string(data), str)
	log.PanicErr(ioutil.WriteFile(file, []byte(result), 0644))
}
