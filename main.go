package main

import(
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func main() {
	fmt.Println("Welcome to Battleship field designer!")
	fmt.Println("Please, populate a battlefield with battle-ships.")
	fmt.Println("Note, a size of the field is random, it can change every time you run this program.")
	fmt.Println("Note, each vessel will be placed at a random location, vertically, horizontally or even diagonally.")
	printCommands()
	field := initializeField(getFieldSize())
	printField(field)
	runTheLoop(field)
}

/**
 * Runs the "loop" of the Battleship field designer.
 */
func runTheLoop(field [][]string) {
	for {
		fmt.Print("Please, enter your command: ")
		var command string
		fmt.Scanf("%s ", &command)

		switch command {
		case "q":
			exitTheProgram()
		case "h":
			printCommands()
			continue
		case "1":
			placeSubmarine(field)
		case "2":
			placeDestroyer(field)
		case "3":
			placeCruiser(field)
		case "4":
			placeCarrier(field)
		default:
			fmt.Printf("Command %s is unknown, please, choose valid one.", command)
			continue
		}

		clearTerminal()
		printField(field)
	}
}

/**
 * Clears the terminal, which creates an effect of shapes been added within a "static screen".
 */
func clearTerminal() {
	clear := map[string]func() {
		"linux": func() {
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()
		},

		"windows": func() {
			cmd := exec.Command("cmd", "/c", "cls")
			cmd.Stdout = os.Stdout
			cmd.Run()
		},
	}

	if f, platformSupported := clear[runtime.GOOS]; platformSupported { // runtime.GOOS: Linux, Windows etc.
		f()
		return
	}

	fmt.Println("Function which clears your terminal screen is not compatible with your OS.")
	fmt.Println("Only Windows and Linux are supported.")
}

/**
 * Prints a list of all available commands.
 */
func printCommands() {
	fmt.Println("Commands:")
	fmt.Println(`"1 + ENTER" - places a submarine at a random location.`)
	fmt.Println(`"2 + ENTER" - places a destroyer at a random location.`)
	fmt.Println(`"3 + ENTER" - places a cruiser at a random location.`)
	fmt.Println(`"4 + ENTER" - places a carrier at a random location.`)
	fmt.Println(`"h + ENTER" - prints the commands list.`)
	fmt.Println(`"q + ENTER" - exists the program.`)
	fmt.Println()
}

/**
 * Exits the program.
 */
func exitTheProgram() {
	os.Exit(0)
}

/**
 * Initializes two-dimensional slice, which will function as an actual battlefield.
 */
func initializeField(size int) [][]string {
	field := make([][]string, size)

	for i := range field {
		field[i] = make([]string, size)

		for j := range field[i] {
			field[i][j] = " "
		}
	}

	return field
}

/**
 * Places a submarine (1 x 1) on the board at a random location.
 */
func placeSubmarine(field [][]string) {
	placeShip(field, ShipTypes.Submarine)
}

/**
 * Places a destroyer (2 x 1) on the board at a random location.
 */
func placeDestroyer(field [][]string) {
	placeShip(field, ShipTypes.Destroyer)
}

/**
 * Places a cruiser (3 x 1) on the board at a random location.
 */
func placeCruiser(field [][]string) {
	placeShip(field, ShipTypes.Cruiser)
}

/**
 * Places a carrier (4 x 1) on the board at a random location.
 */
func placeCarrier(field [][]string) {
	placeShip(field, ShipTypes.Carrier)
}

/**
 * Places a vessel of given type on the board at a random location.
 * If chosen location does not fit the ship, then the function will retry to choose a location up to 200 times.
 * In case suitable location is not found - appropriate message will be printed, and the program will terminate.
 */
func placeShip(field [][]string, shipType uint8) {
	for i := 0; i < 200; i++ { // This is the "retry loop".
		if placeShipAtVertex(field, getShipVertex(field), shipType) {
			return
		}
	}

	fmt.Print("Cannot generate valid position for your vessel.")
	exitTheProgram()
}

/**
 * Places requested vessel at the specified location.
 */
func placeShipAtVertex(field [][]string, vertex map[string]int, shipType uint8) bool {
	shipLength := getShipLength(shipType)
	shipCoordinates := make([]map[string]int, 0)
	var shipDirection uint8
	var lastLocationAdjacentCoordinates DirectionToCoordinate

	for i := 0; i < shipLength; i++ {
		if i == 0 {
			// Check randomly chosen coordinate.
			// In case it is free - allocate it.
			// During the next iteration (if applicable) only surroundings of this coordinate will be checked for allocation.
			if ok, adjacentCoordinates := adjacentCoordinatesFree(field, vertex["x"], vertex["y"], Directions.None); ok {
				shipCoordinates = append(shipCoordinates, vertex)
				lastLocationAdjacentCoordinates = adjacentCoordinates
				continue
			} else {
				return false
			}
		}

		if i == 1 {
			// Check surroundings of previously allocated coordinate to determine a "direction" of a new vessel.
			// Allocate a new coordinate, if possible.
			for direction, lastLocationAdjacentCoordinate := range lastLocationAdjacentCoordinates {
				x, y := lastLocationAdjacentCoordinate["x"], lastLocationAdjacentCoordinate["y"]

				if ok, adjacentCoordinates := adjacentCoordinatesFree(field, x, y, getDirectionToSkip(direction)); ok {
					shipCoordinates = append(shipCoordinates, lastLocationAdjacentCoordinate)
					lastLocationAdjacentCoordinates = adjacentCoordinates
					shipDirection = direction
					break
				}
			}

			if len(shipCoordinates) == 1 {
				// The second coordinate for the vessel is not found.
				// The function must return false to let the caller (func placeShip()) to retry.
				return false
			}

			continue
		}

		// Try to allocate remaining coordinates to complete "vessel's shaping".
		// For now the task is simpler than during first and second iterations,
		// since the "direction" of the ship is determined, hence there is only one coordinate to be checked.
		nextShipCoordinate := lastLocationAdjacentCoordinates[shipDirection]
		x, y := nextShipCoordinate["x"], nextShipCoordinate["y"]

		if ok, adjacentCoordinates := adjacentCoordinatesFree(field, x, y, getDirectionToSkip(shipDirection)); ok {
			shipCoordinates = append(shipCoordinates, nextShipCoordinate)
			lastLocationAdjacentCoordinates = adjacentCoordinates
		}
	}

	// Finally check if the vessel is shaped properly (vessel reached required length).
	if len(shipCoordinates) != shipLength {
		return false
	} else {
		for _, coordinate := range shipCoordinates {
			field[coordinate["x"]][coordinate["y"]] = "*"
		}

		return true
	}
}

/**
 * Accepts ship's direction, and returns the opposite one.
 * It is used to check ship's adjacent coordinates, excluding those already allocated by given ship.
 */
func getDirectionToSkip(shipDirection uint8) uint8 {
	var directionToSkip uint8

	switch shipDirection {
	case Directions.North:
		directionToSkip = Directions.South
	case Directions.South:
		directionToSkip = Directions.North
	case Directions.East:
		directionToSkip = Directions.West
	case Directions.West:
		directionToSkip = Directions.East
	case Directions.NorthWest:
		directionToSkip = Directions.SouthEast
	case Directions.NorthEast:
		directionToSkip = Directions.SouthWest
	case Directions.SouthWest:
		directionToSkip = Directions.NorthEast
	case Directions.SouthEast:
		directionToSkip = Directions.NorthWest
	case Directions.None:
		directionToSkip = Directions.None
	}

	return directionToSkip
}

/**
 * Returns ship's length in accordance to its type.
 */
func getShipLength(shipType uint8) int {
	var shipLength int

	switch shipType {
	case ShipTypes.Submarine:
		shipLength = 1
	case ShipTypes.Destroyer:
		shipLength = 2
	case ShipTypes.Cruiser:
		shipLength = 3
	case ShipTypes.Carrier:
		shipLength = 4
	}

	return shipLength
}

/**
 * Checks if adjacent coordinates are free.
 * If given coordinate (x, y) is approved, then returns a list of adjacent coordinates.
 */
func adjacentCoordinatesFree(field [][]string, x int, y int, directionToSkip uint8) (bool, DirectionToCoordinate) {
	fieldLength := len(field)
	vertexEngaged := field[x][y] == "*"

	if vertexEngaged {
		return false, nil
	}

	directionsValidationMap := map[uint8]bool {
		Directions.West: x > 0 && field[x - 1][y] == "*",
		Directions.East: x + 1 < fieldLength && field[x + 1][y] == "*",
		Directions.North: y + 1 < fieldLength && field[x][y + 1] == "*",
		Directions.South: y > 0 && field[x][y - 1] == "*",
		Directions.NorthWest: x > 0 && y + 1 < fieldLength && field[x - 1][y + 1] == "*",
		Directions.SouthWest: x > 0 && y > 0 && field[x - 1][y - 1] == "*",
		Directions.NorthEast: x + 1 < fieldLength && y + 1 < fieldLength && field[x + 1][y + 1] == "*",
		Directions.SouthEast: x + 1 < fieldLength && y > 0 && field[x + 1][y - 1] == "*",
	}

	for direction, isEngaged := range directionsValidationMap {
		if direction != directionToSkip && isEngaged {
			return false, nil
		}
	}

	return true, getAdjacentCoordinates(x, y, directionToSkip, fieldLength)
}

/**
 * Returns a list of coordinates, adjacent to current one (x, y).
 */
func getAdjacentCoordinates(x int, y int, directionToSkip uint8, fieldLength int) DirectionToCoordinate {
	directionToCoordinate := make(DirectionToCoordinate)

	if directionToSkip != Directions.West && x > 0 {
		directionToCoordinate[Directions.West] = map[string]int {
			"x": x - 1,
			"y": y,
		}
	}

	if directionToSkip != Directions.East && x + 1 < fieldLength {
		directionToCoordinate[Directions.East] = map[string]int {
			"x": x + 1,
			"y": y,
		}
	}

	if directionToSkip != Directions.North && y + 1 < fieldLength {
		directionToCoordinate[Directions.North] = map[string]int {
			"x": x,
			"y": y + 1,
		}
	}

	if directionToSkip != Directions.South && y > 0 {
		directionToCoordinate[Directions.South] = map[string]int {
			"x": x,
			"y": y - 1,
		}
	}

	if directionToSkip != Directions.NorthWest && x > 0 && y + 1 < fieldLength {
		directionToCoordinate[Directions.NorthWest] = map[string]int {
			"x": x - 1,
			"y": y + 1,
		}
	}

	if directionToSkip != Directions.SouthWest && x > 0 && y > 0 {
		directionToCoordinate[Directions.SouthWest] = map[string]int {
			"x": x - 1,
			"y": y - 1,
		}
	}

	if directionToSkip != Directions.NorthEast && x + 1 < fieldLength && y + 1 < fieldLength {
		directionToCoordinate[Directions.NorthEast] = map[string]int {
			"x": x + 1,
			"y": y + 1,
		}
	}

	if directionToSkip != Directions.SouthEast && x + 1 < fieldLength && y > 0 {
		directionToCoordinate[Directions.SouthEast] = map[string]int {
			"x": x + 1,
			"y": y - 1,
		}
	}

	return directionToCoordinate
}

/**
 * Generates a pseudo-random coordinate.
 */
func getShipVertex(field [][]string) map[string]int {
	rand.Seed(time.Now().UnixNano())
	fieldLength := len(field)

	return map[string]int {
		"x": rand.Intn(fieldLength),
		"y": rand.Intn(fieldLength),
	}
}

/**
 * Generated a pseudo-random number, which will be used as the field size.
 */
func getFieldSize() int {
	rand.Seed(time.Now().UnixNano())
	const least = 10
	return rand.Intn(least) + least
}

/**
 * Prints the battlefield.
 */
func printField(field [][]string) {
	fieldLength := len(field)

	for i := 0; i < fieldLength; i++ {
		printHorizontalBar(fieldLength)
		fmt.Print("|")

		for j := 0; j < fieldLength; j++ {
			fmt.Printf(" %s |", field[i][j])
		}

		fmt.Print("\n")
	}

	printHorizontalBar(fieldLength)
}

/**
 * Prints a horizontal bar within the battlefield.
 */
func printHorizontalBar(fieldLength int) {
	for i := 0; i < fieldLength; i++ {
		if i < fieldLength - 1 {
			fmt.Print("----")
			continue
		}

		fmt.Print("-----\n")
	}
}

/**
 * DirectionsRegistry struct, that mimics an enum.
 */
type DirectionsRegistry struct {
	North uint8
	South uint8
	West uint8
	East uint8
	NorthWest uint8
	NorthEast uint8
	SouthWest uint8
	SouthEast uint8
	None uint8
}

/**
 * Initializes DirectionsRegistry struct, that mimics an enum.
 */
func initializeDirectionsEnum() *DirectionsRegistry {
	return &DirectionsRegistry{
		North:     0,
		South:     1,
		West:      2,
		East:      3,
		NorthWest: 4,
		NorthEast: 5,
		SouthWest: 6,
		SouthEast: 7,
		None:      8,
	}
}

/**
 * Map of directory (e.g. West) to closest adjacent coordinate (e.g. { "x": x - 1, "y": y }).
 */
type DirectionToCoordinate map[uint8]map[string]int

/**
 * ShipTypeRegistry struct, that mimics an enum.
 */
type ShipTypeRegistry struct {
	Submarine uint8
	Destroyer uint8
	Cruiser uint8
	Carrier uint8
}

/**
 * Initializes ShipTypeRegistry struct, that mimics an enum.
 */
func initializeShipTypeEnum() *ShipTypeRegistry {
	return &ShipTypeRegistry{
		Submarine: 0,
		Destroyer: 1,
		Cruiser:   2,
		Carrier:   3,
	}
}

// Enables Directions and ShipTypes to be used as regular enums.
// E.g.: Directions.West, ShipTypes.Submarine.
var Directions = initializeDirectionsEnum()
var ShipTypes = initializeShipTypeEnum()
