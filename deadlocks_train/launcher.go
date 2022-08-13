package main

import (
	"log"
	. "multithreading/deadlocks_train/common"
	. "multithreading/deadlocks_train/deadlock"
	"sync"

	"github.com/hajimehoshi/ebiten"
)

var (
	trains        [4]*Train
	intersections [4]*Intersection
)

const TrainLength = 70

func update(screen *ebiten.Image) error {
	if !ebiten.IsDrawingSkipped() {
		DrawTracks(screen)
		DrawIntersections(screen)
		DrawTrains(screen)
	}
	return nil
}

func main() {
	for i := 0; i < 4; i++ {
		trains[i] = &Train{Id: i, TrainLength: TrainLength, Front: 0}
	}

	for i := 0; i < 4; i++ {
		intersections[i] = &Intersection{Id: i, Mutex: sync.Mutex{}, LockedBy: -1}
	}

	go MoveTrain(trains[0], 300, []*Crossing{{Position: 125, Intersection: intersections[0]},
		{Position: 175, Intersection: intersections[1]}})

	go MoveTrain(trains[1], 300, []*Crossing{{Position: 125, Intersection: intersections[1]},
		{Position: 175, Intersection: intersections[2]}})

	go MoveTrain(trains[2], 300, []*Crossing{{Position: 125, Intersection: intersections[2]},
		{Position: 175, Intersection: intersections[3]}})

	go MoveTrain(trains[3], 300, []*Crossing{{Position: 125, Intersection: intersections[3]},
		{Position: 175, Intersection: intersections[0]}})

	if err := ebiten.Run(update, 320, 320, 3, "Trains in a box"); err != nil {
		log.Fatal(err)
	}
}
