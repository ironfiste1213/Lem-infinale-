package main

import (
	"fmt"
	"lem-in/internal"
	"os"
)

func main() {
	filename := os.Args[1]
	file, err := os.Open(filename)
	Graph, err := internal.Parser(file)
	if err != nil {
		fmt.Println("erro get file", err)
		return
	}

	Groupsofpaths := internal.FindAllGroupsOfPath(Graph)
	biggroup := 0
	for _, paths := range Groupsofpaths {
		fmt.Println("group of pathS ,LEN() ", len(paths))
		if biggroup < len((paths)) {
			biggroup = len(paths)
		}
			for index, path := range paths {
				fmt.Println("len of path ", index,"is ",len(path.Rooms))
				for _, RoomsName := range path.Rooms {
					fmt.Print(RoomsName.Id, "-->")
				}
				fmt.Println("")
			}
			fmt.Println("")
		}
		fmt.Println("max group of path have :", biggroup, "path")

		Group := internal.FindBestGroup(Graph.AntCount, Groupsofpaths)
		fmt.Println("best group len is ", len(Group))
		for _, line := range Graph.File {
			fmt.Println(line)
		}
		fmt.Println("")
		internal.SimulateAntsSmart(Graph, Group)
	}

