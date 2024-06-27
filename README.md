# Tic-Tac-Toe AI

## Disclaimer

This project is just a playground for me to experiment with Go routines. This project is ideal for high loading testing
due to extremely high CPU usage.

## The challenge

To build a Tic-Tac-Toe AI that never loses you should predict all the possible moves and outcomes of the game. Sounds
easy for a 3x3 board, but what about a 5x5 board? Or even a 7x7 board?
The number of possible moves and outcomes grows exponentially with the board size, making it impossible to calculate all
the possible moves and outcomes in real-time.

The formula to calculate approximate number of possible moves is:

```
n = board_size * board_size
m = n!
```

Where `n` is the number of cells on the board and `m` is the number of possible moves.

So in a 3x3 board, there are 9 cells and 362,880 possible moves. In a 5x5 board, there are 25 cells and
15,511,210,043,330,985,984,000,000 possible moves. Now it's harder, right?

Each game have also win length that is the number of cells in a row that a player needs to win. For example, in a 3x3
board, the win length is 3. In a 5x5 board, the win length is 5. And you need to calculate all the possible moves and
outcomes for each win length.

## The solution

The solution is to pre-calculate all the possible moves and outcomes for a given board size and store them in map files.
To build a map file, we need to calculate all the possible moves and outcomes for each game state. This is a
time-consuming
process, but it only needs to be done once for each board size. Once the map files are built, the AI can quickly look up
the optimal move for any given game state.

So in the `internal/map_builder` package, we have a `MapBuilder` struct that makes the queue with channels and
goroutines
and `internal/map_storage` stores those maps in the file system. To prevent to much fil writes per second, we have a
splitter that splits the map into smaller parts and writes them to the file system.

Each game board can be represented with a string. For example, a 3x3 board with the following state:

```
X O X
O X O
X O X
```

Can be represented as:

```
"XOXOXOXOX"
```

Also we mark empty cells with a `_` symbol and store the winner at the start of the line. So the empty 3x3 board can be
represented as:

```
"_ _________"
```

In the file system we can split all games in a such way:

```
maps/
├── progress
│   ├── _ ____X____ - the progress of the map building for requested game
├── 3x3_3 - 3x3 board with win length 3
│   ├── ___ - all possible moves and outcomes for the map with "___" in the beggining
│   ├── __O - all possible moves and outcomes for the map with "__O" in the beggining
│   ├── __X - all possible moves and outcomes for the map with "__X" in the beggining
│   ├── _O_ - all possible moves and outcomes for the map with "_O_" in the beggining
│   ├── ... - and so on
├── 5x5_5 - 3x3 board with win length 5
│   ├── ___ - all possible moves and outcomes for the map with "___" in the beggining
│   │   ├── ___ - all possible moves and outcomes for the map with "______" in the beggining
│   │   ├── __O - all possible moves and outcomes for the map with "___O__" in the beggining
│   │   ├── ... - and so on
│   ├── __O - all possible moves and outcomes for the map with "__O" in the beggining
│   │   ├── ___ - all possible moves and outcomes for the map with "__O___" in the beggining
│   │   ├── __O - all possible moves and outcomes for the map with "__OO__" in the beggining
│   │   ├── ... - and so on
│   ├── __X - all possible moves and outcomes for the map with "__X" in the beggining
│   │   ├── ___ - all possible moves and outcomes for the map with "__X___" in the beggining
│   │   ├── __O - all possible moves and outcomes for the map with "__XO__" in the beggining
│   │   ├── ... - and so on
│   ├── _O_ - all possible moves and outcomes for the map with "_O_" in the beggining
│   │   ├── ___ - all possible moves and outcomes for the map with "_O____" in the beggining
│   │   ├── __O - all possible moves and outcomes for the map with "_O_O__" in the beggining
│   │   ├── ... - and so on
│   ├── ... - and so on
├── ... - and so on   
```

## Overview

The Tic-Tac-Toe AI project is a web-based application that allows users to play the classic game of Tic-Tac-Toe against
an AI opponent. The AI uses pre-calculated game maps to determine the best move, ensuring a challenging and engaging
experience for the player.

## Technologies

| Name          | Version |
|---------------|---------|
| Go            | 1.20.6  |
| gin-gonic/gin | 1.9.1   |

## Project Structure

```
/Users/serdnaley/GolandProjects/tic-tac-toe-ai
├── static - contains static files (css, js)
│   ├── css - contains css files
│   └── js - contains javascript files
│       └── components - contains javascript components
├── maps - contains precalculated game maps
└── internal - contains internal go code
    ├── game - contains game logic
    ├── map_builder - contains logic for building game maps
    ├── map_storage - contains logic for storing and retrieving game maps
    ├── map_reader - contains logic for reading game maps
    ├── server - contains server routes and handlers
    └── util - contains utility functions
```

## Business Logic

### Game Logic

The core of the project is the game logic, which is responsible for managing the state of the game, determining valid
moves, and checking for win conditions. The game can be played on different board sizes, such as 3x3 or 5x5.

### AI Opponent

The AI opponent uses pre-calculated game maps to determine the best move. These maps are generated and stored on the
server, allowing the AI to quickly look up the optimal move for any given game state.

### Game Maps

Game maps are essential for the AI's decision-making process. They are built for different board sizes and stored on the
server. The maps contain information about the chances of winning, losing, or drawing for any given game state.

### Client-Side

The client-side is implemented in JavaScript and provides a user interface for playing the game. It communicates with
the server to get the AI's next move, the status of the game map building process, and other game-related information.

### Server-Side

The server-side is implemented in Go and provides several API endpoints:

- `GET /api/health` - Checks the health of the server.
- `GET /api/maps/status` - Gets the status of the game map building process.
- `POST /api/maps/build` - Builds a game map for a specific board size.
- `GET /api/chances` - Gets the chances of winning, losing, or drawing for a given game state.
- `GET /api/next-move` - Gets the next best move for the AI opponent.

### Sequence Diagrams

#### Get Map Status

```sequence
Browser->>Server: GET /api/maps/status?game={game_state}
Server->>File system: Check for map file
File system->>Server: Map file exists or not
Server->>Browser: Map status (progress, ready/not ready)
```

#### Build Map

```sequence
Browser->>Server: POST /api/maps/build?game={game_state}
Server->>File system: Create map file
Server->>File system: Write game data to map file
File system->>Server: Success
Server->>Browser: Map building started
```

#### Get Chances

```sequence
Browser->>Server: GET /api/chances?game={game_state}
Server->>File system: Read map file
File system->>Server: Game data from map
Server->>Browser: Chances of winning, losing, drawing
```

#### Get Next Move

```sequence
Browser->>Server: GET /api/next-move?game={game_state}
Server->>File system: Read map file
File system->>Server: Game data from map
Server->>Browser: Next move coordinates
```

## Setting Up the Development Environment

1. Install Go: Download and install the Go programming language from the official website: https://go.dev/
2. Set up Go workspace: Create a workspace directory and set the GOPATH environment variable to point to it.
3. Install Gin framework: `go get github.com/gin-gonic/gin`

## Running the Project in the Development Environment

### Prerequisites

No prerequisites.

### Execution Instructions

- **Locally**:

```shell
go run main.go
```

This will start the server on `http://localhost:4000` by default. You can access the client-side application by
navigating to `http://localhost:4000/static` in your web browser.

## Testing

The project includes unit tests for the game logic, map building, map storage, and map reading components.

## Deployment

The project can be deployed as a standalone web server.

## Notes

- The project requires a Go development environment to be set up.
- The server must be running for the client to function properly.
