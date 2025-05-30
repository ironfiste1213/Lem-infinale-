package main

import (
	"fmt"
	"lem-in/internal"
	"os"
	"time"
)

func main() {
	t := time.Now()
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run . input.txt")
		os.Exit(1)
	}
	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("ERROR : CAN'T OPEN FILE!  ---> ", filename, err)
		os.Exit(1)
	}
	defer file.Close()
	Graph, err := internal.Parser(file)
	if err != nil {
		fmt.Println("ERROR : INVALID DATA FORMAT! ,", err)
		os.Exit(1)
	}
	Groupsofpaths := internal.FindAllGroupsOfPath(Graph)
	Group := internal.FindBestGroup(Graph.AntCount, Groupsofpaths)
	for _, line := range Graph.File {
		fmt.Println(line)
	}
	fmt.Println("")
	internal.SimulateAntsSmart(Graph, Group)
	fmt.Println(time.Since(t))
}
