# Lem-in Ant Colony Simulation

## Table of Contents
- [Description](#description)
- [Project Structure](#project-structure)
- [Input Format](#input-format)
- [Usage](#usage)
- [Functionality](#functionality)
- [Error Handling](#error-handling)
- [License](#license)

## Description
This project is a Go program that simulates ants moving through a graph representing an ant farm. The graph is defined by rooms and links between them, with designated start and end rooms. The program finds optimal paths for the ants to travel from the start to the end room and simulates their movement.

## Project Structure
```
.
├── README.md                # Project documentation
├── go.mod                   # Go module definition
├── cmd/                     # Additional commands or utilities
│   ├── b                   # Possibly a command or utility
│   ├── main.go              # Alternative main or command entry
│   ├── mars_4000_20_20_95_180_5_no_z(10)  # Data or config file
│   └── test.txt             # Test input or data file
└── internal/                # Internal packages for core logic
    ├── algorithme.go        # Core algorithm for pathfinding and simulation
    ├── models.go            # Internal data models
    ├── output.go            # Internal output handling
    └── parser.go            # Input parsing and validation
```

## Input Format
The program expects an input file describing the ant farm in the following format:

1. The first line contains the number of ants (a positive integer).
2. Room definitions follow, one per line, with the format:
   ```
   room_name x_coordinate y_coordinate
   ```
   - Room names cannot start with 'L' or '#'.
   - Coordinates are integers.
3. Special commands:
   - `##start` indicates the next room is the start room.
   - `##end` indicates the next room is the end room.
4. Link definitions follow, one per line, with the format:
   ```
   room1-room2
   ```
   indicating a bidirectional link between two rooms.

Lines starting with `#` (except `##start` and `##end`) are treated as comments and ignored.

Example:
```
3
##start
A 0 0
B 1 0
##end
C 2 0
A-B
B-C
```

## Usage

### Build
To build the program, run:
```
go build -o lem-in main.go
```

### Run
Run the program with the input file as an argument:
```
./lem-in input.txt
```

The program will parse the input, find optimal paths for the ants, and simulate their movement, printing the simulation output to the console.

## Functionality
- Parses the input file to build a graph of rooms and links.
- Validates the graph structure, ensuring start and end rooms exist and a path connects them.
- Finds groups of disjoint paths from start to end.
- Selects the best group of paths to minimize the number of turns for all ants to reach the end.
- Simulates ants moving along the selected paths and outputs the simulation steps.

## Error Handling
The program handles various input errors, including:
- Empty input file.
- Invalid ant count.
- Missing start or end room.
- Invalid room or link format.
- Duplicate room names.
- No valid path between start and end.

Error messages will be printed to the console if the input is invalid.

## License
This project is open source and free to use.

