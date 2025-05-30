package internal

import (
	"fmt"
	"sort"
)

func FindBestGroup(antCount int, ALLgROUPS [][]*Path) []*Path {

	var bestGroup []*Path
	bestTurns := 1e9

	for _, group := range ALLgROUPS {
		turns := calculateTurns(antCount, group)
		if float64(turns) < bestTurns {
			bestTurns = float64(turns)
			bestGroup = group
		}
	}

	return bestGroup
}

func calculateTurns(antCount int, paths []*Path) int {
	if len(paths) == 0 {
		return 1e9 // big number
	}

	// Paths should be sorted by length ascending (shorter first)
	sort.Slice(paths, func(i, j int) bool {
		return paths[i].Len < paths[j].Len
	})

	turns := 0
	extraAnts := antCount

	// While ants are not all assigned
	for extraAnts > 0 {
		turns++
		for range paths {
			if extraAnts > 0 {
				extraAnts--
			}
		}
	}

	// turns += (length of longest path - 1)
	return turns + paths[len(paths)-1].Len - 1
}

// Add this method to your Graph struct
func (g *Graph) Copy() *Graph {
	// Create new graph
	newGraph := &Graph{
		Rooms:    make(map[string]*Room),
		AntCount: g.AntCount,
	}

	// First, copy all rooms
	for roomID, room := range g.Rooms {
		newRoom := &Room{
			Id:                  room.Id,
			IsStart:             room.IsStart,
			IsEnd:               room.IsEnd,
			Links:               make(map[string]*Room),
			Visited:             false,
			Parent:              nil, 
			Usedinpath:          false,
			X:                   room.X,
			Y:                   room.Y,
			AllowToJump:         false,
			CameFromBacktraking: false,
			Forclinks:           nil,
			ParentInbfs: nil,
		}
		newGraph.Rooms[roomID] = newRoom
	}

	// Set StartRoom and EndRoom references
	if g.StartRoom != nil {
		newGraph.StartRoom = newGraph.Rooms[g.StartRoom.Id]
	}
	if g.EndRoom != nil {
		newGraph.EndRoom = newGraph.Rooms[g.EndRoom.Id]
	}

	// Then copy all links (using the new room references)
	for roomID, room := range g.Rooms {
		for linkID, linkedRoom := range room.Links {
			if linkedRoom != nil {
				newGraph.Rooms[roomID].Links[linkID] = newGraph.Rooms[linkID]
			}
		}
	}

	// Copy paths if needed
	for _, path := range g.Paths {
		newPath := &Path{
			Rooms: make([]*Room, len(path.Rooms)),
		}
		for i, room := range path.Rooms {
			newPath.Rooms[i] = newGraph.Rooms[room.Id]
		}
		newGraph.Paths = append(newGraph.Paths, newPath)
	}

	return newGraph
}
func FindAllGroupsOfPath(g *Graph) [][]*Path {
	groupsofgroups := [][]*Path{}
	groupoflink := [][]string{}

	for {
		// Create a fresh copy each time
		CopyGraph := g.Copy()

		// Remove all previously found conflicting links
		for _, link := range groupoflink {
			if link != nil {
				RoomOne := link[0]
				Roomtow := link[1]
				CopyGraph.RemoveLinks(RoomOne, Roomtow)
			}
		}

		Paths, links := FindGroupOfDisjointPath(CopyGraph)
		if len(Paths) == 0 {
			return groupsofgroups
		}

		// Only add links if they exist and we found paths
		if links != nil {
			groupoflink = append(groupoflink, links)
		}
		groupsofgroups = append(groupsofgroups, Paths)

		// If no conflicting links were found, we're done
		if links == nil {
			return groupsofgroups
		}
	}
}
func FindGroupOfDisjointPath(g *Graph) ([]*Path, []string) {
	GroupOfDisjointPath := []*Path{}
	for {
		Path, links := g.Bfs()
		if Path != nil && Path.Len != 2 {
			if links != nil {
				return GroupOfDisjointPath, links
			}
		}
		if Path == nil {
			return GroupOfDisjointPath, links
		}
		GroupOfDisjointPath = append(GroupOfDisjointPath, Path)
	}

}

func (g *Graph) RemoveLinks(RoomOne, RoomTwo string) {
	room1, ok1 := g.Rooms[RoomOne]
	room2, ok2 := g.Rooms[RoomTwo]

	if !ok1 || !ok2 {
		fmt.Println("Error: One or both rooms not found in graph")
		return
	}

	delete(room1.Links, RoomTwo) //
	//fmt.Println("Remove , ", RoomTwo, "from", RoomOne, "'s , links")
	delete(room2.Links, RoomOne)
	//fmt.Println("Remove ", RoomOne, " from", RoomTwo, "'s links")
}

func (g *Graph) Bfs() (*Path, []string) {
	// Initialize: reset all rooms to unvisited
	g.resetVisited()
	// BFS queue
	queue := []*Room{g.StartRoom}
	g.StartRoom.Visited = true
	g.StartRoom.ParentInbfs = nil
	// BFS to find shortest path
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		// Check all neighbors
		if current.Usedinpath {
			if !current.AllowToJump {
				queue = append(queue, current.Parent)
				current.Parent.Visited = true
				current.Visited = true
				current.Parent.AllowToJump = true
				current.Parent.ParentInbfs = current
				current.CameFromBacktraking = true
				current.Parent.CameFromBacktraking = true
				current.Forclinks = []string{current.Id, current.Parent.Id}
				current.Parent.Forclinks = current.Forclinks
				continue
			} else {
				if current.Parent != nil {
					current.CameFromBacktraking = true
					current.Forclinks = current.ParentInbfs.Forclinks
				}
			}
		}
		for _, neighbor := range current.Links {
			if neighbor.Usedinpath && neighbor.Parent == g.StartRoom {
				continue
			}
			// Skip if already visited
			if neighbor.Visited {
				continue
			}

			if current == g.StartRoom && neighbor.Usedinpath {
				continue
			}

			// Skip start room (can't go back to start)
			if neighbor == g.StartRoom {
				continue
			}

			// Mark as visited and set parent
			neighbor.Visited = true
			neighbor.ParentInbfs = current
			if current.CameFromBacktraking {
				neighbor.CameFromBacktraking = true
				neighbor.Forclinks = current.Forclinks
			}
			// Found the end room - reconstruct and return path
			if neighbor == g.EndRoom {

				path, links := g.reconstructPath(g.EndRoom)
				g.EndRoom.ParentInbfs = current
				return path, links
			}
			// Add to queue for further exploration

			queue = append(queue, neighbor)
		}
	}

	//("No path found from start to end")
	return nil, nil
}

// reconstructPath builds the path from start to end using parent pointers
func (g *Graph) reconstructPath(endRoom *Room) (*Path, []string) {
	var rooms []*Room

	if endRoom.CameFromBacktraking {
		return nil, endRoom.Forclinks
	}
	// Backtrack from end to start
	current := endRoom
	

	for current != nil {
		if current != g.StartRoom && current != g.EndRoom {
			current.Usedinpath = true
			current.Parent = current.ParentInbfs
		}
		rooms = append([]*Room{current}, rooms...) // Prepend to get correct order
		current = current.ParentInbfs
	}
	path := &Path{
		Rooms:                   rooms,
		Len:                     len(rooms) - 1, // Number of edges
	}
	return path,  nil

}

// resetVisited resets all rooms to unvisited state
func (g *Graph) resetVisited() {
	for _, Room := range g.Rooms {
		Room.CameFromBacktraking = false
		Room.ParentInbfs = nil
		Room.AllowToJump = false
		Room.Forclinks = nil
		Room.Visited = false
		if !Room.Usedinpath {
			Room.Parent = nil
		}
	}
}
