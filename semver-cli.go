package main

import (
	"bufio"
	"os"
	"strings"
	"strconv"
	"fmt"
)

const (
	ComponentMajor    = iota
	ComponentMinor    = iota
	ComponentRevision = iota
)

func main() {
	changeType := determineChangeType()
	components := readComponents()

	for i, component := range components {
		if i == changeType {
			versionInt, err := strconv.Atoi(component)
			if nil != err {
				panic(fmt.Sprintf("could not parse version component: %s", component))
			}
			components[i] = strconv.Itoa(versionInt + 1)
		} else {
			components[i] = component
		}
	}

	fmt.Println(strings.Join(components, "."))
}

func readComponents() []string {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		exitWithUsage("error: no version string piped to stdin", 1)
	}
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	version := scanner.Text()
	components := strings.Split(version, ".")
	return components
}

func determineChangeType() int {
	if len(os.Args) < 2 {
		exitWithUsage("semver-cli - a CLI for modifying semver strings", 0)
	}

	switch os.Args[1] {
	case "-M":
		return ComponentMajor
	case "-m":
		return ComponentMinor
	case "-r":
		return ComponentRevision
	default:
		exitWithUsage(fmt.Sprintf("unsupported option: %s", os.Args[1]), 1)
		return -1
	}
}

func exitWithUsage(msg string, code int) {
	fmt.Println(msg)
	printUsage()
	os.Exit(code)
}

func printUsage() {
	fmt.Print(`
Usage:
		
  echo "1.2.3" | semver-cli OPTION
		
Options:
  -M  change major
  -m  change minor
  -r  change revision
`)
}
