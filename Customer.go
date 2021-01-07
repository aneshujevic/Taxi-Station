package main

import (
	"fmt"
	"sync"
	"time"
)

// Request gets sent to the dispatcher
type Request struct {
	customer 	*Customer
	taxis    	chan *Submittance
}

// Customer struct representing a customer
// with starting location and ending location
// along with other info
type Customer struct {
	id 				string
	locationStart 	Location
	locationEnd 	Location
	durationOfRide 	time.Duration
}

// simulates a phone call for taxi
func (c *Customer) makeRequest(d *Dispatcher) {
	d.enqueueCustomerRequest(c)
}

// simulates a ride home
func (c *Customer) takeARideHome(d *Dispatcher, t *Taxi, w *sync.WaitGroup) {
	if t != nil {
		fmt.Printf("Customer %v has taken %v taxi.\n", c.id, t.id)
		// simulates the time ride takes to get home
		time.AfterFunc(c.durationOfRide * time.Second, func() {
			t.setLocationAtomic(c.locationEnd)
			d.addAvailableTaxiAtomic(t)
			fmt.Printf("Taxi %v is now available.\n", t.id)
			w.Done()
		})
	} else {
		// there was no taxi for us :(
		fmt.Printf("Customer %v has gone home walking..\n", c.id)
		w.Done()
	}
}
