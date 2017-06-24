package main

import (
	"fmt"
	"regexp"
	"io/ioutil"
	"strings"
	"strconv"
	"flag"
)

var bumpType = 0
var supportedFiles = []string{
	"package.json",
	"config.xml",
}

func _contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func _versionBumpHelper(version string, index int) string {
	versionParts := strings.Split(version, ".")

	if n, err := strconv.ParseInt(versionParts[index], 10, 8); err == nil {
		versionParts[index] = strconv.FormatInt(n + 1, 10)
	}

	return strings.Join(versionParts, ".")
}

func versionPatch(version string) string {
	return _versionBumpHelper(version, 2)
}

func versionMinor(version string) string {
	return _versionBumpHelper(version, 1)
}

func versionMajor(version string) string {
	return _versionBumpHelper(version, 0)
}

func bumpConfigXml()  {
	data, err := ioutil.ReadFile("config.xml")
	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile(`<widget.*version="([0-9]+\.[0-9]+\.[0-9]+)"`)
	version := re.FindStringSubmatch(string(data))[1]
	fmt.Printf("Current version: %s\n", version)

	var newVersion string
	switch bumpType {
	case 0:
		newVersion = versionPatch(version)
		break
	case 1:
		newVersion = versionMinor(version)
		break
	case 2:
		newVersion = versionMajor(version)
		break
	}
	fmt.Printf("New version: %s\n", newVersion)

	re = regexp.MustCompile(version)
	newData := re.ReplaceAllString(string(data), newVersion)

	ioutil.WriteFile("config.xml", []byte(newData), 0644)
}

func bumpPackageJson()  {
	data, err := ioutil.ReadFile("package.json")
	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile(`"version": "([0-9]+\.[0-9]+\.[0-9]+)"`)
	version := re.FindStringSubmatch(string(data))[1]
	fmt.Printf("Current version: %s\n", version)

	var newVersion string
	switch bumpType {
	case 0:
		newVersion = versionPatch(version)
		break
	case 1:
		newVersion = versionMinor(version)
		break
	case 2:
		newVersion = versionMajor(version)
		break
	}
	fmt.Printf("New version: %s\n", newVersion)

	re = regexp.MustCompile(version)
	newData := re.ReplaceAllString(string(data), newVersion)

	ioutil.WriteFile("package.json", []byte(newData), 0644)
}

func bump(name string)  {
	switch name {
	case "config.xml":
		bumpConfigXml()
		break
	case "package.json":
		bumpPackageJson()
		break
	}
}

func main() {
	var patchFlag = flag.Bool("patch", false, "Patch Bump")
	var minorFlag = flag.Bool("minor", false, "Minor Bump")
	var majorFlag = flag.Bool("major", false, "Major Bump")

	flag.Parse()
	fmt.Printf("Patch: %t, Minor: %t, Major: %t\n", *patchFlag, *minorFlag, *majorFlag)
	if *patchFlag {
		bumpType = 0
	}

	if *minorFlag {
		bumpType = 1
	}

	if *majorFlag {
		bumpType = 2
	}

	files, _ := ioutil.ReadDir("./")
	for _, f := range files {
		name := f.Name()

		if _contains(supportedFiles, name) {
			fmt.Printf("Bumping %s\n", name)
			bump(name)
		}
	}
}
