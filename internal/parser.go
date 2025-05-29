package internal

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

)

// Various error types for specific parsing errors
var (
	ErrEmptyInput        = errors.New("empty input")
	ErrInvalidAntCount   = errors.New("invalid ant count")
	ErrMissingStartRoom  = errors.New("missing start room")
	ErrMissingEndRoom    = errors.New("missing end room")
	ErrNoPathExists      = errors.New("no valid path exists between start and end")
	ErrInvalidRoomFormat = errors.New("invalid room format")
	ErrInvalidLinkFormat = errors.New("invalid link format")
	ErrDuplicateRoom     = errors.New("duplicate room name")
	ErrRoomNotFound      = errors.New("referenced room does not exist")
)

// Parse reads the entire input and constructs an AntFarm object
func Parser(r io.Reader) (*Graph, error) {
	scanner := bufio.NewScanner(r)
	farm := &Graph{
		Rooms: make(map[string]*Room),
		File: []string{},
	}

	// First pass: Parse the ant count
	if !scanner.Scan() {
		return nil, ErrEmptyInput
	}

	line := scanner.Text()
	farm.File = append(farm.File, line)

	antCount, err := strconv.Atoi(line)
	if err != nil || antCount <= 0 {
		return nil, fmt.Errorf("%w: %s", ErrInvalidAntCount, line)
	}
	farm.AntCount = antCount

	// Second pass: Parse rooms and links
	var expectingStart, expectingEnd bool
	parsingLinks := false

	for scanner.Scan() {
		line = scanner.Text()
		farm.File = append(farm.File, line)

		// Skip comments that aren't commands
		if strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "##") {
			continue
		}

		// Check for start/end commands
		if line == "##start" {
			expectingStart = true
			continue
		}
		if line == "##end" {
			expectingEnd = true
			continue
		}

		// If we find a link definition, switch to parsing links
		if strings.Contains(line, "-") && !parsingLinks {
			parsingLinks = true
		}

		if parsingLinks {
			// Parse links between rooms
			if err := parseLink(farm, line); err != nil {
				return nil, err
			}
		} else {
			// Parse room definitions
			if err := parseRoom(farm, line, expectingStart, expectingEnd); err != nil {
				return nil, err
			}
			expectingStart = false
			expectingEnd = false
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Validate the farm structure
	if err := validateFarm(farm); err != nil {
		return nil, err
	}

	return farm, nil
}

// parseRoom parses a room definition line
func parseRoom(farm *Graph, line string, isStart, isEnd bool) error {
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return fmt.Errorf("%w: %s", ErrInvalidRoomFormat, line)
	}

	name := parts[0]
	// Room names can't start with 'L' or '#'
	if strings.HasPrefix(name, "L") || strings.HasPrefix(name, "#") {
		return fmt.Errorf("%w: room name can't start with 'L' or '#'", ErrInvalidRoomFormat)
	}

	x, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("%w: invalid x coordinate", ErrInvalidRoomFormat)
	}

	y, err := strconv.Atoi(parts[2])
	if err != nil {
		return fmt.Errorf("%w: invalid y coordinate", ErrInvalidRoomFormat)
	}

	// Check for duplicate room names
	if _, exists := farm.Rooms[name]; exists {
		return fmt.Errorf("%w: %s", ErrDuplicateRoom, name)
	}

	room := &Room{
		Id:        name,
		X:           x,
		Y:           y,
		IsStart:     isStart,
		IsEnd:       isEnd,
		Links: map[string]*Room{},
	}

	farm.Rooms[name] = room

	if isStart {
		farm.StartRoom = room
	}
	if isEnd {
		farm.EndRoom = room
	}

	return nil
}

// parseLink parses a link definition line
func parseLink(farm *Graph, line string) error {
	parts := strings.Split(line, "-")
	if len(parts) != 2 {
		return fmt.Errorf("%w: %s", ErrInvalidLinkFormat, line)
	}

	room1Name := parts[0]
	room2Name := parts[1]

	room1, exists1 := farm.Rooms[room1Name]
	if !exists1 {
		return fmt.Errorf("%w: %s", ErrRoomNotFound, room1Name)
	}

	room2, exists2 := farm.Rooms[room2Name]
	if !exists2 {
		return fmt.Errorf("%w: %s", ErrRoomNotFound, room2Name)
	}

	// Add the connection in both directions (undirected graph)
	room1.Links[room2Name] = room2
	room2.Links[room1Name] = room1

	return nil
}

// validateFarm ensures the farm structure is valid
func validateFarm(farm *Graph) error {
	if farm.StartRoom == nil {
		return ErrMissingStartRoom
	}
	if farm.EndRoom == nil {
		return ErrMissingEndRoom
	}

	// ✅ New: Check if a path exists from start to end using BFS
	visited := make(map[string]bool)
	queue := []*Room{farm.StartRoom}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == farm.EndRoom {
			return nil // ✅ Path exists
		}

		for _, neighbor := range current.Links {
			if !visited[neighbor.Id] {
				visited[neighbor.Id] = true
				queue = append(queue, neighbor)
			}
		}
	}

	// If we finish BFS without finding EndRoom
	return ErrNoPathExists
}