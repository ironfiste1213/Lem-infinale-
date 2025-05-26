package internal

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ParseInput reads and parses the lem-in input file
func ParseInput(filename string) (*Graph, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("ERROR: invalid data format, cannot open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	graph := &Graph{
		Rooms:     make(map[string]*Room),
		Paths:     []*Path{},
	}

	// Parse ant count (first line)
	if !scanner.Scan() {
		return nil, fmt.Errorf("ERROR: invalid data format, no ant count")
	}
	
	antCount, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil || antCount <= 0 {
		return nil, fmt.Errorf("ERROR: invalid data format, invalid number of ants")
	}
	graph.AntCount = antCount

	var isStart, isEnd bool
	var pendingLinks []string

	// Parse rooms and links
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines and comments (unless they're special commands)
		if line == "" || (strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "##")) {
			continue
		}

		// Handle special commands
		if line == "##start" {
			isStart = true
			continue
		} else if line == "##end" {
			isEnd = true
			continue
		}

		// Check if it's a link (contains '-')
		if strings.Contains(line, "-") {
			pendingLinks = append(pendingLinks, line)
			continue
		}

		// Parse room
		room, err := parseRoom(line, isStart, isEnd)
		if err != nil {
			return nil, err
		}

		// Check for duplicate rooms
		if _, exists := graph.Rooms[room.Id]; exists {
			return nil, fmt.Errorf("ERROR: invalid data format, duplicate room: %s", room.Id)
		}

		graph.Rooms[room.Id] = room

		// Set start/end room references
		if room.IsStart {
			if graph.StartRoom != nil {
				return nil, fmt.Errorf("ERROR: invalid data format, multiple start rooms")
			}
			graph.StartRoom = room
		}
		if room.IsEnd {
			if graph.EndRoom != nil {
				return nil, fmt.Errorf("ERROR: invalid data format, multiple end rooms")
			}
			graph.EndRoom = room
		}

		// Reset flags
		isStart = false
		isEnd = false
	}

	// Validate that we have start and end rooms
	if graph.StartRoom == nil {
		return nil, fmt.Errorf("ERROR: invalid data format, no start room found")
	}
	if graph.EndRoom == nil {
		return nil, fmt.Errorf("ERROR: invalid data format, no end room found")
	}

	// Process links
	err = processLinks(graph, pendingLinks)
	if err != nil {
		return nil, err
	}

	return graph, nil
}

// parseRoom parses a room line and creates a Room struct
func parseRoom(line string, isStart, isEnd bool) (*Room, error) {
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return nil, fmt.Errorf("ERROR: invalid data format, invalid room format: %s", line)
	}

	name := parts[0]
	
	// Validate room name
	if err := validateRoomName(name); err != nil {
		return nil, err
	}

	// Parse coordinates
	x, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("ERROR: invalid data format, invalid x coordinate: %s", parts[1])
	}

	y, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("ERROR: invalid data format, invalid y coordinate: %s", parts[2])
	}

	return &Room{
		Id:           name,
		Links:        make(map[string]*Room),
		IsStart:      isStart,
		IsEnd:        isEnd,
		Visited:      false,
		X:            x,
		Y:            y,
	}, nil
}

// validateRoomName checks if room name follows the rules
func validateRoomName(name string) error {
	if name == "" {
		return fmt.Errorf("ERROR: invalid data format, empty room name")
	}
	
	// Room cannot start with 'L'
	if strings.HasPrefix(name, "L") {
		return fmt.Errorf("ERROR: invalid data format, room name cannot start with 'L': %s", name)
	}
	
	// Room cannot start with '#'
	if strings.HasPrefix(name, "#") {
		return fmt.Errorf("ERROR: invalid data format, room name cannot start with '#': %s", name)
	}
	
	// Room cannot contain spaces
	if strings.Contains(name, " ") {
		return fmt.Errorf("ERROR: invalid data format, room name cannot contain spaces: %s", name)
	}
	
	return nil
}

// processLinks creates connections between rooms
func processLinks(graph *Graph, links []string) error {
	for _, link := range links {
		parts := strings.Split(link, "-")
		if len(parts) != 2 {
			return fmt.Errorf("ERROR: invalid data format, invalid link format: %s", link)
		}

		room1Name := strings.TrimSpace(parts[0])
		room2Name := strings.TrimSpace(parts[1])

		// Get rooms
		room1, exists1 := graph.Rooms[room1Name]
		if !exists1 {
			return fmt.Errorf("ERROR: invalid data format, unknown room in link: %s", room1Name)
		}

		room2, exists2 := graph.Rooms[room2Name]
		if !exists2 {
			return fmt.Errorf("ERROR: invalid data format, unknown room in link: %s", room2Name)
		}

		// Check for self-links
		if room1 == room2 {
			return fmt.Errorf("ERROR: invalid data format, room cannot link to itself: %s", room1Name)
		}

		// Check for duplicate links
		if isAlreadyLinked(room1, room2) {
			return fmt.Errorf("ERROR: invalid data format, duplicate link: %s-%s", room1Name, room2Name)
		}

		// Create bidirectional link
		room1.Links[room2.Id] = room2
		room2.Links[room1.Id] = room1

		// Update LinksToRm map
		graph.Rooms[room1.Id] = room1
		graph.Rooms[room2.Id] = room2
	}

	return nil
}

// isAlreadyLinked checks if two rooms are already connected
func isAlreadyLinked(room1, room2 *Room) bool {
	for _, link := range room1.Links {
		if link == room2 {
			return true
		}
	}
	return false
}

// ValidateGraph performs final validation on the parsed graph
func ValidateGraph(graph *Graph) error {
	// Check if there's at least one path from start to end
	if !HasPathToEnd(graph) {
		return fmt.Errorf("ERROR: invalid data format, no path from start to end")
	}
	
	return nil
}

// hasPathToEnd performs a simple BFS to check if end is reachable from start
func HasPathToEnd(graph *Graph) bool {
	if graph.StartRoom == nil || graph.EndRoom == nil {
		return false
	}

	// Reset visited status
	for _, room := range graph.Rooms {
		room.Visited = false
	}

	queue := []*Room{graph.StartRoom}
	graph.StartRoom.Visited = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == graph.EndRoom {
			for _, room := range graph.Rooms {
				room.Visited = false
			}
			return true
		}

		for _, neighbor := range current.Links {
			if !neighbor.Visited {
				neighbor.Visited = true
				queue = append(queue, neighbor)
			}
		}
	}
	for _, room := range graph.Rooms {
		room.Visited = false
	}

	return false
}

// PrintOriginalInput prints the original input format (for output display)
func PrintOriginalInput(graph *Graph, originalLines []string) {
	for _, line := range originalLines {
		fmt.Println(line)
	}
	fmt.Println() // Empty line before ant movements
}

// ReadOriginalLines reads the file and keeps original lines for output
func ReadOriginalLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}