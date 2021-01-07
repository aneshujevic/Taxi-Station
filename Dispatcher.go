package main

import (
	"fmt"
	"sync"
	"time"
)

// Time each customer spends waiting for the dispatcher to see
// it there is any available taxi for it
const TIMEOUT = 5

type Dispatcher struct {
	availableTaxis    []*Taxi
	taxiMutex         sync.RWMutex
	requests          chan *Request
	customerTaxiMutex sync.RWMutex
}

// push customer request into dispatchers requests queue
func (d *Dispatcher) enqueueCustomerRequest(p *Customer) {
	// In case available number of available taxis gets changed
	d.taxiMutex.RLock()
	taxiNo := len(d.availableTaxis)
	d.taxiMutex.RUnlock()

	d.requests <- &Request{
		customer: p,
		taxis:    make(chan *Submittance, taxiNo),
	}
}

// broadcast request to all the available taxis
func (d *Dispatcher) broadcastRequestsToTaxis() {
	wg := sync.WaitGroup{}
	for request := range d.requests {
		// Wait for TIMEOUT seconds before checking if any taxi is available
		fmt.Println("Customer " + request.customer.id + " is waiting for dispatcher..")
		time.Sleep(time.Duration(TIMEOUT) * time.Second)

		// w WaitGroup is responsible for closing channel for
		// the certain request
		w := sync.WaitGroup{}
		for _, taxi := range d.availableTaxis {
			w.Add(1)
			go taxi.evaluateAndSubmit(request.customer.locationStart, request.taxis, &w)
		}

		go func() {
			w.Wait()
			close(request.taxis)
		}()

		// when closest taxi is found, remove it from list of available ones
		// wg WaitGroup responsible for waiting until all requests finish with their rides
		// so we can see the printed messages :)
		chosenTaxi := d.findClosestTaxi(request)
		d.removeAvailableTaxiAtomic(chosenTaxi)
		wg.Add(1)
		request.customer.takeARideHome(d, chosenTaxi, &wg)
	}
	wg.Wait()
}

// find the taxi that is closest to the requested location
// and is willing to take the passenger
func (d *Dispatcher) findClosestTaxi(ctp *Request) *Taxi {
	closestSubmittance := <-ctp.taxis
	for submit := range ctp.taxis {
		if submit.distance < closestSubmittance.distance {
			closestSubmittance = submit
		}
	}

	if closestSubmittance != nil {
		return closestSubmittance.taxi
	}
	return nil
}

// add taxi to the list of available taxis
func (d *Dispatcher) addAvailableTaxiAtomic(t *Taxi) {
	d.taxiMutex.Lock()
	d.availableTaxis = append(d.availableTaxis, t)
	d.taxiMutex.Unlock()
}

// remove taxi from the list of available taxis
// the algorithm is like this, otherwise we get memory leak
func (d *Dispatcher) removeAvailableTaxiAtomic(t *Taxi) {
	d.taxiMutex.Lock()
	for i, taxi := range d.availableTaxis {
		if taxi == t {
			length := len(d.availableTaxis)
			copy(d.availableTaxis[i:], d.availableTaxis[i+1:])
			d.availableTaxis[length-1] = nil
			d.availableTaxis = d.availableTaxis[:length-1]
		}
	}
	d.taxiMutex.Unlock()
}
