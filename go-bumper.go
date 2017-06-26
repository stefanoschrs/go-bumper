package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

var supportedFiles = map[string]string{
	"package.json": "package.json",
	"config.xml":   "config.xml",
}

func versionBumpHelper(version string, index int) string {
	versionParts := strings.Split(version, ".")

	if n, err := strconv.ParseInt(versionParts[index], 10, 8); err == nil {
		versionParts[index] = strconv.FormatInt(n+1, 10)
	}

	return strings.Join(versionParts, ".")
}

func versionPatch(version string) string {
	return versionBumpHelper(version, 2)
}

func versionMinor(version string) string {
	return versionBumpHelper(version, 1)
}

func versionMajor(version string) string {
	return versionBumpHelper(version, 0)
}

func bumpFile(bumpType int, fileName string, pattern string) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile(pattern)
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

	ioutil.WriteFile(fileName, []byte(newData), 0644)
}

func bump(name string, bumpType int) {
	switch name {
	case "config.xml":
		bumpFile(bumpType, "config.xml",
			`<widget.*version="([0-9]+\.[0-9]+\.[0-9]+)"`)
		break
	case "package.json":
		bumpFile(bumpType, "package.json",
			`"version": "([0-9]+\.[0-9]+\.[0-9]+)"`)
		break
	}
}

func main() {
	bumpType := parseFlags()

	files, _ := ioutil.ReadDir("./")
	for _, f := range files {
		name := f.Name()

		if _, ok := supportedFiles[name]; ok {
			fmt.Printf("Bumping %s\n", name)
			bump(name, bumpType)
		}
	}
}

func parseFlags() (bumpType int) {
	patchFlag := flag.Bool("patch", false, "Patch Bump")
	minorFlag := flag.Bool("minor", false, "Minor Bump")
	majorFlag := flag.Bool("major", false, "Major Bump")

	flag.Parse()
	fmt.Printf("Patch: %t, Minor: %t, Major: %t\n",
		*patchFlag, *minorFlag, *majorFlag)

	if *patchFlag {
		bumpType = 0
	} else if *minorFlag {
		bumpType = 1
	} else if *majorFlag {
		bumpType = 2
	} else {
		bumpType = 0
	}

	return
}
