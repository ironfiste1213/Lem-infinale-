package internal

import (
	"fmt"
	"sort"
	"strings"
)

type SimAnt struct {
	ID       int
	Path     []*Room
	Step     int
	Finished bool
}

func assignAntsSmart(farm *Graph, paths []*Path) []*SimAnt {
	sort.Slice(paths, func(i, j int) bool {
		return paths[i].Len < paths[j].Len
	})

	antsPerPath := make([]int, len(paths))
	remainingAnts := farm.AntCount

	// Assign ants to paths based on minimal cost (path length + number of ants already assigned)
	for remainingAnts > 0 {
		bestIdx := 0
		minCost := antsPerPath[0] + paths[0].Len

		for i := 1; i < len(paths); i++ {
			cost := antsPerPath[i] + paths[i].Len
			if cost < minCost {
				minCost = cost
				bestIdx = i
			}
		}

		antsPerPath[bestIdx]++
		remainingAnts--
	}

	ants := make([]*SimAnt, 0, farm.AntCount)
	antID := 1

	for pathIdx, count := range antsPerPath {
		for j := 0; j < count; j++ {
			ant := &SimAnt{
				ID:   antID,
				Path: paths[pathIdx].Rooms,
				Step: 0,
			}
			ants = append(ants, ant)
			antID++
		}
	}

	return ants
}

func SimulateAntsSmart(farm *Graph, paths []*Path) {
	if len(paths) == 0 {
		fmt.Println("No paths available!")
		return
	}

	ants := assignAntsSmart(farm, paths)
	movingAnts := []*SimAnt{}
	waitingAnts := make([]*SimAnt, len(ants))
	copy(waitingAnts, ants)

	stepCount := 0 // Add step counter

	for {
		moveLine := ""
		occupied := make(map[string]bool)

		// Move ants already walking
		for _, ant := range movingAnts {
			if !ant.Finished && ant.Step < len(ant.Path)-1 {
				nextRoom := ant.Path[ant.Step+1]
				if nextRoom.IsEnd || !occupied[nextRoom.Id] {
					ant.Step++
					moveLine += fmt.Sprintf("L%d-%s ", ant.ID, ant.Path[ant.Step].Id)
					if nextRoom.IsEnd {
						ant.Finished = true
					} else {
						occupied[nextRoom.Id] = true
					}
				}
			}
		}

		// Push as many waiting ants as possible this turn
		newMovingAnts := []*SimAnt{}
		remaining := []*SimAnt{}

		for _, ant := range waitingAnts {
			nextRoom := ant.Path[1]
			if !occupied[nextRoom.Id] || nextRoom.IsEnd {
				ant.Step = 1
				moveLine += fmt.Sprintf("L%d-%s ", ant.ID, nextRoom.Id)
				if !nextRoom.IsEnd {
					occupied[nextRoom.Id] = true
					newMovingAnts = append(newMovingAnts, ant)
				}
			} else {
				remaining = append(remaining, ant)
			}
		}

		waitingAnts = remaining

		for _, ant := range movingAnts {
			if !ant.Finished {
				newMovingAnts = append(newMovingAnts, ant)
			}
		}

		movingAnts = newMovingAnts

		if moveLine != "" {
			fmt.Println(strings.TrimSpace(moveLine))
			stepCount++ // Increment step counter when moves are made
		}

		if len(movingAnts) == 0 && len(waitingAnts) == 0 {
			break
		}
	}

	// Print the total number of steps
	fmt.Printf("Total steps: %d\n", stepCount)
}