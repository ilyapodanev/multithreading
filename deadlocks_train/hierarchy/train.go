package hierarchy

import (
	. "multithreading/deadlocks_train/common"
	"sort"
	"time"
)

func lockIntersectionsInDistance(id, reserveStart, reserveEnd int, crossings []*Crossing) {
	var IntersectionsToLock []*Intersection
	for _, crossing := range crossings {
		if reserveEnd >= crossing.Position && reserveStart <= crossing.Position && crossing.Intersection.LockedBy != id {
			IntersectionsToLock = append(IntersectionsToLock, crossing.Intersection)
		}
	}

	sort.Slice(IntersectionsToLock, func(i, j int) bool {
		return IntersectionsToLock[i].Id < IntersectionsToLock[j].Id
	})

	for _, it := range IntersectionsToLock {
		it.Mutex.Lock()
		it.LockedBy = id
		time.Sleep(10 + time.Millisecond)
	}
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
				crossing.Intersection.LockedBy = -1
				crossing.Intersection.Mutex.Unlock()
			}
		}
		time.Sleep(30 * time.Millisecond)
	}
}
