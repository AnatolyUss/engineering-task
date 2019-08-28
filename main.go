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
 * Clears the terminal.
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

	fmt.Println("Cannot clear your terminal screen.")
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
 * Initializes two-dimensional slice, which will function as actual field.
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
 * If chosen location does not fit the ship, then the function will retry to choose a location up to 100 times.
 * In case suitable location is not found - appropriate message will be printed, and the program will terminate.
 */
func placeShip(field [][]string, shipType uint8) {
	for i := 0; i < 100; i++ { // This is the "retry loop".
		if placeShipAtVertex(field, getShipVertex(field), shipType) {
			return
		}
	}

	fmt.Print("Cannot generate valid position for your vessel.")
	os.Exit(0)
}

/**
 * Checks if adjacent coordinates are free.
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

	return true, getAdjacentCoordinates(x, y, directionToSkip)
}

/**
 * TODO: add description.
 */
func getAdjacentCoordinates(x int, y int, directionToSkip uint8) DirectionToCoordinate {
	directionToCoordinate := make(DirectionToCoordinate)

	if directionToSkip != Directions.West {
		directionToCoordinate[Directions.West] = map[string]int {
			"x": x - 1,
			"y": y,
		}
	}

	if directionToSkip != Directions.East {
		directionToCoordinate[Directions.East] = map[string]int {
			"x": x + 1,
			"y": y,
		}
	}

	if directionToSkip != Directions.North {
		directionToCoordinate[Directions.North] = map[string]int {
			"x": x,
			"y": y + 1,
		}
	}

	if directionToSkip != Directions.South {
		directionToCoordinate[Directions.South] = map[string]int {
			"x": x,
			"y": y - 1,
		}
	}

	if directionToSkip != Directions.NorthWest {
		directionToCoordinate[Directions.NorthWest] = map[string]int {
			"x": x - 1,
			"y": y + 1,
		}
	}

	if directionToSkip != Directions.SouthWest {
		directionToCoordinate[Directions.SouthWest] = map[string]int {
			"x": x - 1,
			"y": y - 1,
		}
	}

	if directionToSkip != Directions.NorthEast {
		directionToCoordinate[Directions.NorthEast] = map[string]int {
			"x": x + 1,
			"y": y + 1,
		}
	}

	if directionToSkip != Directions.SouthEast {
		directionToCoordinate[Directions.SouthEast] = map[string]int {
			"x": x + 1,
			"y": y - 1,
		}
	}

	return directionToCoordinate
}

/**
 * TODO: complete the body, and add the description.
 */
func placeShipAtVertex(field [][]string, vertex map[string]int, shipType uint8) bool {
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

	shipCoordinates := make([]map[string]int, 0)
	// var shipDirection uint8
	var lastLocationAdjacentCoordinates DirectionToCoordinate

	for i := 0; i < shipLength; i++ {
		if i == 0 {
			if ok, adjacentCoordinates := adjacentCoordinatesFree(field, vertex["x"], vertex["y"], Directions.None); ok {
				shipCoordinates = append(shipCoordinates, vertex)
				lastLocationAdjacentCoordinates = adjacentCoordinates
				continue
			} else {
				break
			}
		}

		for direction, lastLocationAdjacentCoordinate := range lastLocationAdjacentCoordinates {
			x, y := lastLocationAdjacentCoordinate["x"], lastLocationAdjacentCoordinate["y"]

			if ok, adjacentCoordinates := adjacentCoordinatesFree(field, x, y, direction); ok {
				shipCoordinates = append(shipCoordinates, lastLocationAdjacentCoordinate)
				lastLocationAdjacentCoordinates = adjacentCoordinates
				// shipDirection = direction
				break
			}
		}
	}

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
	const least = 15
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
 * Prints a horizontal bar.
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
 * ShipLengthsRegistry struct, that mimics an enum.
 */
type ShipTypeRegistry struct {
	Submarine uint8
	Destroyer uint8
	Cruiser uint8
	Carrier uint8
}

/**
 * Initializes ShipType struct, that mimics an enum.
 */
func initializeShipTypeEnum() *ShipTypeRegistry {
	return &ShipTypeRegistry{
		Submarine: 0,
		Destroyer: 1,
		Cruiser:   2,
		Carrier:   3,
	}
}

var Directions = initializeDirectionsEnum()
var ShipTypes = initializeShipTypeEnum()
