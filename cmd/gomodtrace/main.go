package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/godepsresolve/gomodtrace"
)

const minimalRatio = 0.9

func tracePaths(parent, target string) {
	input := readInput()
	graph := gomodtrace.ParseGraph(input)
	index := gomodtrace.BuildGraphIndex(graph)
	if _, ok := index[parent]; !ok {
		fmt.Println(prepareExtendedMessage(index, minimalRatio, parent, "parent"))
		return
	}
	if _, ok := index[target]; !ok {
		fmt.Println(prepareExtendedMessage(index, minimalRatio, target, "target"))
		return
	}

	paths := index.FindPaths(parent, target, nil)
	log.Printf("%v\n", paths)
	involvedLibraries := paths.ListInvolvedLibraries()
	fmt.Println(graph.WithOnly(involvedLibraries))
}

// stringMatchRatio calculates match ratio between two strings.
func stringMatchRatio(a, b string) float64 {
	if a == b {
		return 1.0
	}

	minLenSequence := min(len(a), len(b))
	if minLenSequence == 0 {
		return 0
	}

	found := float64(0)
	for i := 0; i < minLenSequence; i++ {
		if a[i] == b[i] {
			found++
		}
	}

	return found / float64(minLenSequence)
}

// prepareExtendedMessage prepares an extended message based on the
// string matching ratio.
func prepareExtendedMessage(
	index gomodtrace.NodeIndex,
	minRatio float64,
	element,
	elementType string,
) string {
	message := "check input"
	for k := range index {
		if stringMatchRatio(k, element) > minRatio {
			message = fmt.Sprintf("did you mean '%s'?", k)
			break
		}
	}

	return fmt.Sprintf(
		"No %s package: '%s' found in package index, %s",
		elementType,
		element,
		message,
	)
}

func readInput() []string {
	var buf []string
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		buf = append(buf, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	if len(buf) == 0 {
		fmt.Println("No input provided")
		help()
	}
	return buf
}

func hasStdInput() bool {
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

func help() {
	flag.Usage()
	os.Exit(1)
}

func setUsage() {
	flag.Usage = func() {
		fmt.Println("Usage of gomodtrace: go mod graph | gomodtrace [OPTION]... PARENT_PACKAGE DEPENDENT_PACKAGE")
		flag.PrintDefaults()
	}
}

func getArgs() (string, string, bool) {
	needVerbose := flag.Bool("v", false, "use verbose mode")
	flag.Parse()
	parent, target := flag.Arg(0), flag.Arg(1)
	if parent == "" || target == "" {
		fmt.Println("Required arguments PARENT_PACKAGE or DEPENDENT_PACKAGE were not provided")
		help()
	}
	return parent, target, *needVerbose
}

func main() {
	setUsage()

	parent, target, needVerbose := getArgs()
	if !needVerbose {
		log.SetOutput(io.Discard)
	}

	if !hasStdInput() {
		fmt.Println("Provide input by `go mod graph | gomodtrace ROOT_PACKAGE DEPENDENT_PACKAGE` or `gomodtrace ROOT_PACKAGE DEPENDENT_PACKAGE < gomodgraph.txt`")
		help()
	}
	tracePaths(parent, target)
}
