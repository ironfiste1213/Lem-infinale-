package internal

import "fmt"

// bfs should look for path if it reach a room that allready used in path go trowt parent of this room after that it can jump
// if we reach the end trout a path that have room allready used backtarkingpath it should delet the link
func FindGroupOfDisjointPath(g *Graph) []*Path {
	GroupOfDisjointPath := []*Path{}
	FirstRound := true
	AllReadyBactraking := false
	for {
		Path, IsBacktrakingPath := g.Bfs(FirstRound, AllReadyBactraking)
		if Path == nil {
			break
		}
		if !IsBacktrakingPath {

			GroupOfDisjointPath = append(GroupOfDisjointPath, Path)
			FirstRound = false
		} else {

		}
	}
	return GroupOfDisjointPath
}

func (g *Graph) RemoveLinks(RoomOne, RoomTwo string) {
	room1, ok1 := g.Rooms[RoomOne]
	room2, ok2 := g.Rooms[RoomTwo]

	if !ok1 || !ok2 {
		fmt.Println("Error: One or both rooms not found in graph")
		return
	}

	delete(room1.Links, RoomTwo) // Remove RoomTwo from RoomOne's links
	delete(room2.Links, RoomOne) // Remove RoomOne from RoomTwo's links
}

func (g *Graph) Bfs(FirstRound bool, AllreadyBacktring bool) (*Path, bool) {
	// Initialize: reset all rooms to unvisited
	g.resetVisited(FirstRound, AllreadyBacktring)

	// Track parent for path reconstruction
	parent := make(map[*Room]*Room)

	// BFS queue
	queue := []*Room{g.StartRoom}
	g.StartRoom.Visited = true
	parent[g.StartRoom] = nil

	fmt.Printf("Starting BFS from room: %s\n", g.StartRoom.Id)

	// BFS to find shortest path
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		fmt.Printf("Exploring from room: %s\n", current.Id)
		// Check all neighbors
		if !AllreadyBacktring {
			if current.Usedinpath {
				if !current.AllowToJump {
					fmt.Println("_______add as backtraking", current.Parent.Id)
					queue = append(queue, current.Parent)
					current.Parent.Visited = true
					current.Parent.AllowToJump = true
					parent[current] = current.Parent
					current.CameFromBacktraking = true
					continue
				} else {
					if current != g.StartRoom && current != g.EndRoom {
						current.Parent.AllowToJump = true
					}
				}
			}
		}
		for _, neighbor := range current.Links {
			// Skip if already visited
			if neighbor.Visited {
				continue
			}
			if neighbor.Blocked {
				continue
			}

			if current == g.StartRoom && neighbor.Usedinpath {
				continue
			}

			// Skip start room (can't go back to start)
			if neighbor == g.StartRoom {
				continue
			}

			fmt.Printf(" -> Checking neighbor: %s\n", neighbor.Id)

			// Mark as visited and set parent
			neighbor.Visited = true
			parent[neighbor] = current

			// Found the end room - reconstruct and return path
			if neighbor == g.EndRoom {

				path, isbacktrakingpath := g.reconstructPath(parent, g.EndRoom)
				fmt.Println("----------End ROOM Found in path isbacktraking == ", isbacktrakingpath)
				fmt.Printf("End room found! Reconstructing path...\n")
				return path, isbacktrakingpath
			}

			// Add to queue for further exploration

			queue = append(queue, neighbor)
		}
	}

	fmt.Println("No path found from start to end")
	return nil, false
}

// reconstructPath builds the path from start to end using parent pointers
func (g *Graph) reconstructPath(parent map[*Room]*Room, endRoom *Room) (*Path, bool) {
	var rooms []*Room
	isbacktrakingpath := false
	// Backtrack from end to start
	current := endRoom
	for current != nil {
		if current.CameFromBacktraking {
			isbacktrakingpath = true
		}
		rooms = append([]*Room{current}, rooms...) // Prepend to get correct order
		current = parent[current]
	}

	path := &Path{
		Rooms:                   rooms,
		Len:                     len(rooms) - 1, // Number of edges
		IsFollowingPathBackward: isbacktrakingpath,
	}
	next := g.StartRoom
	if !path.IsFollowingPathBackward {
		for _, Room := range path.Rooms {
			if Room != g.EndRoom && Room != g.StartRoom {
				Room.Usedinpath = true
				Room.Parent = next
				next = Room
			}
		}
	}
	fmt.Printf("Path reconstructed with %d rooms (length %d)\n", len(rooms), path.Len)
	return path, isbacktrakingpath
}

// resetVisited resets all rooms to unvisited state
func (g *Graph) resetVisited(FirstRound bool, AllReadyBacktraking bool) {
	for _, Room := range g.Rooms {
		if Room.Usedinpath && AllReadyBacktraking {
			Room.Blocked = true
		}
		Room.Visited = false
	}

}
