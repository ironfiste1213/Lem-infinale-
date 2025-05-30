package main

import (
	"fmt"
	"lem-in/internal"
	"os"
)

func main() {
	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("ERROR : can't open file ---> ", filename)
	}
	Graph, err := internal.Parser(file)
	if err != nil {
		fmt.Println("ERROR : INVALID DATA FORMAT! ,", err)
		return
	}

	    Groupsofpaths := internal.FindAllGroupsOfPath(Graph)
		Group := internal.FindBestGroup(Graph.AntCount, Groupsofpaths)
		for _, line := range Graph.File {
			fmt.Println(line)
		}
		fmt.Println("")
		internal.SimulateAntsSmart(Graph, Group)
	}

