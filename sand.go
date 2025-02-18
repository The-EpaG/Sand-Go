package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	gridWidth     = 100
	gridHeight    = 60
	maxGenerations = 1000
	
	// Colors
	colorWhite    = termbox.ColorWhite
	colorYellow   = termbox.ColorYellow
	colorDefault  = termbox.ColorDefault

	// Characters
	charEmpty         = ' '
	charSand          = '█'
	charLeftTopCorner = '╔'
	charRightTopCorner = '╗'
	charLeftBottomCorner = '╚'
	charRightBottomCorner = '╝'
	charHorizontalBorder = '═'
	charVerticalBorder   = '║'
)

type Cell struct {
	filled bool
}

func drawGrid(grid [][]Cell) {
	// Draw corners
	termbox.SetCell(0, 0, charLeftTopCorner, colorWhite, colorDefault)
	termbox.SetCell(gridWidth+1, 0, charRightTopCorner, colorWhite, colorDefault)
	termbox.SetCell(0, gridHeight+1, charLeftBottomCorner, colorWhite, colorDefault)
	termbox.SetCell(gridWidth+1, gridHeight+1, charRightBottomCorner, colorWhite, colorDefault)

	// Draw horizontal borders
	for x := 1; x <= gridWidth; x++ {
		termbox.SetCell(x, 0, charHorizontalBorder, colorWhite, colorDefault)
		termbox.SetCell(x, gridHeight+1, charHorizontalBorder, colorWhite, colorDefault)
	}

	// Draw grid and side borders
	for y := 0; y < gridHeight; y++ {
		// Draw side borders
		termbox.SetCell(0, y+1, charVerticalBorder, colorWhite, colorDefault)
		termbox.SetCell(gridWidth+1, y+1, charVerticalBorder, colorWhite, colorDefault)
		
		for x := 0; x < gridWidth; x++ {
			character := charEmpty
			fgColor := colorDefault
			if grid[y][x].filled {
				character = charSand
				fgColor = colorYellow
			}
			termbox.SetCell(x+1, y+1, character, fgColor, colorDefault)
		}
	}
}

func main() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	// Initialize the grid
	grid := make([][]Cell, gridHeight)
	for i := range grid {
		grid[i] = make([]Cell, gridWidth)
	}

	// Get the number of sand particles from the command line
	args := os.Args
	if len(args) != 2 {
		fmt.Println("Usage:", os.Args[0], "<number_of_sand_particles>")
		return
	}
	particleCount, err := strconv.Atoi(args[1])
	if err != nil || particleCount <= 0 || particleCount > gridWidth*gridHeight/4 {
		fmt.Println("Invalid number of sand particles. Must be a positive integer and less than or equal to", gridWidth*gridHeight/4)
		return
	}

	// Initialize the grid and randomly place sand particles
	for i := 0; i < particleCount; i++ {
		x := rand.Intn(gridWidth)
		y := rand.Intn(gridHeight)
		grid[y][x].filled = true
	}

	// Simulate sand falling
	previousGrid := copyGrid(grid)

	// Draw the grid
	termbox.Clear(colorDefault, colorDefault)
	drawGrid(grid)

	termbox.Flush()
	time.Sleep(100 * time.Millisecond)

	for generation := 0; generation < maxGenerations; generation++ {
		// Clear the screen
		termbox.Clear(colorDefault, colorDefault)

		// Update the grid state
		for i := gridHeight - 2; i >= 0; i-- {
			for j := 0; j < gridWidth; j++ {
				if grid[i][j].filled {
					if i+1 < gridHeight && !grid[i+1][j].filled { // Downward
						grid[i+1][j].filled = true
						grid[i][j].filled = false
					} else if i+1 < gridHeight && j-1 >= 0 && !grid[i+1][j-1].filled { // Downward left
						grid[i+1][j-1].filled = true
						grid[i][j].filled = false
					} else if i+1 < gridHeight && j+1 < gridWidth && !grid[i+1][j+1].filled { // Downward right
						grid[i+1][j+1].filled = true
						grid[i][j].filled = false
					}
				}
			}
		}

		// Draw the grid
		drawGrid(grid)

		termbox.Flush()
		time.Sleep(50 * time.Millisecond)

		if gridsEqual(grid, previousGrid) {
			time.Sleep(1000 * time.Millisecond)
			break
		}
		previousGrid = copyGrid(grid)
	}
}

func copyGrid(grid [][]Cell) [][]Cell {
	newGrid := make([][]Cell, gridHeight)
	for i := range newGrid {
		newGrid[i] = make([]Cell, gridWidth)
		for j := range newGrid[i] {
			newGrid[i][j] = grid[i][j]
		}
	}
	return newGrid
}

func gridsEqual(grid1, grid2 [][]Cell) bool {
	for i := 0; i < gridHeight; i++ {
		for j := 0; j < gridWidth; j++ {
			if grid1[i][j].filled != grid2[i][j].filled {
				return false
			}
		}
	}
	return true
}

