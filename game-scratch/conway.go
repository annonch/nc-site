// Copyright 2018 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	// "bytes"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	// "image"
	// _ "image/png"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	// "github.com/hajimehoshi/ebiten/v2/examples/resources/images"
)

var GlobalCount int

type CellState struct {
	Alive bool
	Image *ebiten.Image
}

// UpdateStrategy defines an interface for different update strategies.
type UpdateStrategy interface {
	Update(ca *CellularAutomaton, x, y int) bool
}

// CellularAutomaton represents the grid of cells in the automaton.
type CellularAutomaton struct {
	Width      int
	Height     int
	Grid       [][]CellState
	Buffer     [][]CellState
	Strategy   UpdateStrategy // Strategy used to update the grid
	DeadImage  *ebiten.Image
	AliveImage *ebiten.Image
}

// NewCellularAutomaton creates a new CellularAutomaton with a given width and height.
func NewCellularAutomaton(width, height int, strategy UpdateStrategy) *CellularAutomaton {
	grid := make([][]CellState, height)
	buffer := make([][]CellState, height)
	for i := range grid {
		grid[i] = make([]CellState, width)
		buffer[i] = make([]CellState, width)
	}

	aliveImage := createCellImage(true)
	deadImage := createCellImage(false)
	return &CellularAutomaton{
		Width:      width,
		Height:     height,
		Grid:       grid,
		Buffer:     buffer,
		Strategy:   strategy,
		AliveImage: aliveImage,
		DeadImage:  deadImage,
	}
}

// Update updates the state of the cellular automaton based on the current strategy.
func (ca *CellularAutomaton) Update() {
	for y := 0; y < ca.Height; y++ {
		for x := 0; x < ca.Width; x++ {
			newState := ca.Strategy.Update(ca, x, y)
			var image *ebiten.Image
			if newState {
				image = ca.AliveImage
			} else {
				image = ca.DeadImage
			}
			ca.Buffer[y][x] = CellState{
				Alive: newState,
				Image: image,
			}
		}
	}

	ca.Grid, ca.Buffer = ca.Buffer, ca.Grid
}

// GameOfLifeStrategy implements the rules for Conway's Game of Life.
type GameOfLifeStrategy struct{}

// Update applies the Game of Life rules to a cell at position (x, y).
func (g *GameOfLifeStrategy) Update(ca *CellularAutomaton, x, y int) bool {
	aliveNeighbors := countAliveNeighbors(ca, x, y)
	currentCellAlive := ca.Grid[y][x].Alive

	if currentCellAlive && (aliveNeighbors < 2 || aliveNeighbors > 3) {
		return false
	} else if !currentCellAlive && aliveNeighbors == 3 {
		return true
	}
	return currentCellAlive
}

// countAliveNeighbors counts the alive neighbors of a cell at position (x, y).
func countAliveNeighbors(ca *CellularAutomaton, x, y int) int {
	alive := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dy == 0 && dx == 0 {
				continue
			}

			nx, ny := x+dx, y+dy
			if nx >= 0 && ny >= 0 && nx < ca.Width && ny < ca.Height && ca.Grid[ny][nx].Alive {
				alive++
			}
		}
	}
	return alive
}

// InitializeGridRandomly initializes the grid of the cellular automaton with random states.
func (ca *CellularAutomaton) InitializeGridRandomly(seed int64) {
	rand.Seed(seed)
	for y := 0; y < ca.Height; y++ {
		for x := 0; x < ca.Width; x++ {
			alive := rand.Intn(2) == 1
			ca.Grid[y][x] = CellState{
				Alive: alive,
				Image: createCellImage(alive),
			}
		}
	}
}

// createCellImage creates an image for a cell based on its state.
func createCellImage(alive bool) *ebiten.Image {
	img := ebiten.NewImage(cellSize, cellSize)
	if alive {
		// img.Fill(color.RGBA{0, 255, 0, 255}) // Green for alive
		// img.Fill(color.RGBA{17, 153, 17, 255}) // dark Green for alive
		// img.Fill(color.RGBA{85, 68, 238, 255}) // vilot for alive
		img.Fill(color.RGBA{0, 0, 0, 255}) // black
	} else {
		img.Fill(color.RGBA{255, 255, 255, 127}) // White for dead
	}
	return img
}

// Game represents the game state.
type Game struct {
	CA *CellularAutomaton
}

// Update updates the game logic.
func (g *Game) Update() error {
	// Check for mouse click
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.addGlider(x/cellSize, y/cellSize) // Assuming each cell is 10x10 pixels
	}

	g.CA.Update()
	GlobalCount++
	return nil
}

func (g *Game) addGlider(x, y int) {
	if x < 1 || x >= g.CA.Width-2 || y < 1 || y >= g.CA.Height-2 {
		return // Not enough space for a glider
	}

	// Coordinates to set cells to alive for forming a glider
	//
	// . x .
	// . . x
	// x x x    (x, y) x x
	gliderCoords := [][]int{
		{x, y}, {x, y + 1}, {x, y + 2}, {x + 1, y + 2}, {x + 2, y + 1},
	}

	for _, coord := range gliderCoords {
		g.CA.Grid[coord[1]][coord[0]].Alive = true
		g.CA.Grid[coord[1]][coord[0]].Image = g.CA.AliveImage
	}
}

// Draw draws the game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	for y, row := range g.CA.Grid {
		for x, cell := range row {
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(x*cellSize), float64(y*cellSize)) // Position each cell
			screen.DrawImage(cell.Image, opts)
		}
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f,  c: %i", ebiten.ActualTPS(), GlobalCount))
}

// Layout returns the screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.CA.Width * cellSize, g.CA.Height * cellSize // Adjust the size based on the grid
}

const (
	screenWidth  = 1024
	screenHeight = 1024
	cellSize     = 8
	tps          = 8
)

func init() {
	GlobalCount = 0
}

// main function to demonstrate the usage
func main() {
	golStrategy := &GameOfLifeStrategy{}
	ca := NewCellularAutomaton(screenWidth/cellSize, screenWidth/cellSize, golStrategy)

	var seed int64
	// seed = 22
	// seed = 42
	// seed = 420
	seed = 4
	ca.InitializeGridRandomly(seed)

	game := &Game{CA: ca}
	ebiten.SetWindowSize(ca.Width*cellSize, ca.Height*cellSize)
	ebiten.SetWindowTitle("Cellular Automaton")
	ebiten.SetTPS(tps)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}

// func init() {
// 	// Decode an image from the image file's byte slice.
// 	// img, _, err := image.Decode(bytes.NewReader(images.Tiles_png))
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// tiles
// 	// Image = ebiten.NewImageFromImage(img)
// }

// type Game struct {
// 	layers [][]int
// }

// func (g *Game) Update() error {
// 	return nil
// }

// func (g *Game) Draw(screen *ebiten.Image) {
// 	w := tilesImage.Bounds().Dx()
// 	tileXCount := w / tileSize

// 	// Draw each tile with each DrawImage call.
// 	// As the source images of all DrawImage calls are always same,
// 	// this rendering is done very efficiently.
// 	// For more detail, see https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Image.DrawImage
// 	const xCount = screenWidth / tileSize
// 	for _, l := range g.layers {
// 		for i, t := range l {
// 			op := &ebiten.DrawImageOptions{}
// 			op.GeoM.Translate(float64((i%xCount)*tileSize), float64((i/xCount)*tileSize))

// 			sx := (t % tileXCount) * tileSize
// 			sy := (t / tileXCount) * tileSize
// 			screen.DrawImage(tilesImage.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
// 		}
// 	}

// 	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
// }

// func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
// 	return screenWidth, screenHeight
// }

// func main() {
// 	g := &Game{
// 		layers: [][]int{
// 			{
// 				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
// 				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
// 				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
// 				243, 218, 243, 243, 243, 243, 243, 243, 243, 243, 243, 218, 243, 244, 243,
// 				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,

// 				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
// 				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
// 				243, 243, 244, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
// 				243, 243, 243, 243, 243, 243, 243, 243, 243, 219, 243, 243, 243, 219, 243,
// 				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,

// 				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
// 				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
// 				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
// 				243, 218, 243, 243, 243, 243, 243, 243, 243, 243, 243, 244, 243, 243, 243,
// 				243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243, 243,
// 			},
// 			{
// 				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
// 				0, 0, 0, 0, 0, 26, 27, 28, 29, 30, 31, 0, 0, 0, 0,
// 				0, 0, 0, 0, 0, 51, 52, 53, 54, 55, 56, 0, 0, 0, 0,
// 				0, 0, 0, 0, 0, 76, 77, 78, 79, 80, 81, 0, 0, 0, 0,
// 				0, 0, 0, 0, 0, 101, 102, 103, 104, 105, 106, 0, 0, 0, 0,

// 				0, 0, 0, 0, 0, 126, 127, 128, 129, 130, 131, 0, 0, 0, 0,
// 				0, 0, 0, 0, 0, 303, 303, 245, 242, 303, 303, 0, 0, 0, 0,
// 				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
// 				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
// 				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,

// 				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
// 				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
// 				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
// 				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
// 				0, 0, 0, 0, 0, 0, 0, 245, 242, 0, 0, 0, 0, 0, 0,
// 			},
// 		},
// 	}

// 	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
// 	ebiten.SetWindowTitle("Tiles (Ebitengine Demo)")
// 	if err := ebiten.RunGame(g); err != nil {
// 		log.Fatal(err)
// 	}
// }
