package main

import (
	"sync"
)

type Taxi struct {
	id 			string
	location 	Location
	mutex		sync.RWMutex
}

// evaluate the distance between the current taxi
// and the customer location and submit the info
// through the channel
// w WaitGroup is needed for closing the channel submittances
// after all the taxis have finished submitting
func (t *Taxi) evaluateAndSubmit(customerLocation Location, submittances chan *Submittance, w *sync.WaitGroup) {
	// mutex in case location gets changed in the moment of measuring distance
	t.mutex.RLock()
	distance := t.getDistance(customerLocation)
	t.mutex.RUnlock()
	submittances <- &Submittance{
		taxi:     t,
		distance: distance,
	}
	w.Done()
}

func (t *Taxi) setLocationAtomic(l Location) {
	// mutex so no one can read while we're changing location
	t.mutex.Lock()
	t.location = l
	t.mutex.Unlock()
}

// get distance from current location to the other given location
func (t *Taxi) getDistance(l Location) float64 {
	return t.location.DistanceTo(l)
}
