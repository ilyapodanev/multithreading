package arbitrator

import (
	. "multithreading/deadlocks_train/common"
	"sync"
	"time"
)

var (
	controller = sync.Mutex{}
	cond       = sync.NewCond(&controller)
)

func allFree(intersectionsToLock []*Intersection) bool {
	for _, it := range intersectionsToLock {
		if it.LockedBy >= 0 {
			return false
		}
	}
	return true
}

func lockIntersectionsInDistance(id, reserveStart, reserveEnd int, crossings []*Crossing) {
	var IntersectionsToLock []*Intersection
	for _, crossing := range crossings {
		if reserveEnd >= crossing.Position && reserveStart <= crossing.Position && crossing.Intersection.LockedBy != id {
			IntersectionsToLock = append(IntersectionsToLock, crossing.Intersection)
		}
	}

	controller.Lock()
	for !allFree(IntersectionsToLock) {
		cond.Wait()
	}
	for _, it := range IntersectionsToLock {
		it.LockedBy = id
		time.Sleep(10 + time.Millisecond)
	}
	controller.Unlock()
}

func MoveTrain(train *Train, distance int, crossings []*Crossing) {
	for train.Front < distance {
		train.Front += 1
		for _, crossing := range crossings {
			if train.Front == crossing.Position {
				lockIntersectionsInDistance(train.Id, crossing.Position, crossing.Position+train.TrainLength, crossings)
			}
			back := train.Front - train.TrainLength
			if back == crossing.Position {
				controller.Lock()
				crossing.Intersection.LockedBy = -1
				cond.Broadcast()
				controller.Unlock()
			}
		}
		time.Sleep(30 * time.Millisecond)
	}
}