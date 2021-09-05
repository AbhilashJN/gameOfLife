package main

import (
	"flag"
	"log"

	"github.com/pkg/profile"
)

func initWorld(w *World) {
	// simple glider pattern
	w[1][2] = 1
	w[2][3] = 1
	w[3][1] = 1
	w[3][2] = 1
	w[3][3] = 1

}

func main() {
	// setup profiler in trace mode
	defer profile.Start(profile.TraceProfile, profile.ProfilePath(".")).Stop()

	// initialize world with some pattern
	var world World
	initWorld(&world)
	var nextWorldState World
	var calcNextState func(*World, *World)

	// check method flag value
	method := flag.Int("m", -1, "Select state update method")
	flag.Parse()

	if *method == -1 {
		log.Fatal(`
			Usage: gameOfLife -m <method_number>
			Flags: 
				-m :  Method used for updating state.
				Values
					1 : update cells sequentially.
					2 : update all cells using one goroutine per cell.
					3 : update each row using one goroutine per row.
					4 : update using workers, one goroutine per worker. Each worker processes one row at a time.
		`)
	}

	switch *method {
	case 1:
		calcNextState = calculateNextWorldState //sequential
	case 2:
		calcNextState = calculateNextWorldStateCellParallel // cell parallel
	case 3:
		calcNextState = calculateNextWorldStateRowParallel // row parallel
	case 4:
		calcNextState = calculateNextWorldStateRowWorker // worker pattern
	}

	// world.display()
	calcNextState(&world, &nextWorldState)

	// // Run 100 iterations
	// for i := 0; i < 100; i++ {
	// 	time.Sleep(1 * time.Second)
	// 	calcNextState(&world, &nextWorldState)
	// 	world = nextWorldState
	// 	world.display()
	// }

}
