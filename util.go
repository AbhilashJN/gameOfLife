package main

import (
	"fmt"
	"sync"
)

const (
	maxRows = 3000
	maxCols = 3000
)

type World [maxRows][maxCols]int8

// print current world to stdout. Use small grid sizes to view output properly on stdout.
func (w *World) display() {
	for row := 0; row < maxRows; row++ {
		for col := 0; col < maxCols; col++ {
			if w[row][col] == 1 {
				fmt.Printf(" O ")
			} else {
				fmt.Printf(" - ")
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n\n\n")
}

// // Synthetic work load
// func calculateNextStepAt(currentWorldState *World, x int, y int) int8 {
// 	var temp int
// 	for i := 0; i < 100; i++ {
// 		temp += i / ((x * y) + 1)
// 	}
// 	return int8(temp % 2)
// }

func calculateNextStepAt(currentWorldState *World, x int, y int) int8 {
	isCurrentCellLive := currentWorldState[x][y] == 1
	var aliveNeighbours int8 = 0
	for row := x - 1; row <= x+1; row++ {
		for col := y - 1; col <= y+1; col++ {
			if row < 0 || row >= maxRows || col < 0 || col >= maxCols {
				continue
			}
			if row == x && col == y {
				continue
			}
			aliveNeighbours += currentWorldState[row][col]
		}
	}

	/*
		Rules:
		   Any live cell with two or three live neighbours survives.
		   Any dead cell with three live neighbours becomes a live cell.
		   All other live cells die in the next generation. Similarly, all other dead cells stay dead.
	*/
	switch {
	case isCurrentCellLive && (aliveNeighbours == 2 || aliveNeighbours == 3):
		return 1
	case !isCurrentCellLive && aliveNeighbours == 3:
		return 1
	default:
		return 0
	}
}

func calculateNextWorldState(currentWorldState *World, nextWorldState *World) {
	for row := 0; row < maxRows; row++ {
		for col := 0; col < maxCols; col++ {
			nextWorldState[row][col] = calculateNextStepAt(currentWorldState, row, col)
		}
	}
}

func calculateNextWorldStateCellParallel(currentWorldState *World, nextWorldState *World) {
	var wg sync.WaitGroup
	wg.Add(maxRows * maxCols)
	for row := 0; row < maxRows; row++ {
		for col := 0; col < maxCols; col++ {
			go func(row, col int, w *World) {
				nextWorldState[row][col] = calculateNextStepAt(w, row, col)
				wg.Done()
			}(row, col, currentWorldState)
		}
	}
	wg.Wait()
}

func calculateNextWorldStateRowParallel(currentWorldState *World, nextWorldState *World) {
	var wg sync.WaitGroup
	wg.Add(maxRows)
	for row := 0; row < maxRows; row++ {
		go func(row int, w *World) {
			for col := 0; col < maxCols; col++ {
				nextWorldState[row][col] = calculateNextStepAt(w, row, col)
			}
			wg.Done()
		}(row, currentWorldState)
	}
	wg.Wait()
}

func calculateNextWorldStateRowWorker(currentWorldState *World, nextWorldState *World) {
	c := make(chan int, maxRows)
	workers := 8
	var wg sync.WaitGroup
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func(w *World) {
			for row := range c {
				for col := 0; col < maxCols; col++ {
					nextWorldState[row][col] = calculateNextStepAt(w, row, col)
				}
			}
			wg.Done()
		}(currentWorldState)
	}
	for row := 0; row < maxRows; row++ {
		c <- row
	}
	close(c)
	wg.Wait()
}
