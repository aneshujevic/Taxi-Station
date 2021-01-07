package main

import "sync"

// instantiate some customers, taxis and dispatcher
// and test them out
func main() {
	c1 := Customer{
		id:             "Harry Potter",
		locationStart:  Location{1,1},
		locationEnd:    Location{5,5},
		durationOfRide: 10,
	}

	c2 := Customer{
		id:             "Yugi",
		locationStart:  Location{1,1},
		locationEnd:    Location{5,5},
		durationOfRide: 10,
	}

	c3 := Customer{
		id:             "Yami",
		locationStart:  Location{1,1},
		locationEnd:    Location{5,5},
		durationOfRide: 15,
	}
	
	t1 := Taxi{
		id:        "Kadric",
		location:  Location{1,2},
		mutex:     sync.RWMutex{},
	}
	
	t2 := Taxi{
		id:        "Zuti",
		location:  Location{4,4},
		mutex:     sync.RWMutex{},
	}
	
	d := Dispatcher{
		availableTaxis:     []*Taxi{&t1, &t2},
		taxiMutex:          sync.RWMutex{},
		requests: 			make(chan *Request, 5),
		customerTaxiMutex:  sync.RWMutex{},
	}

	// wg WaitGroup responsible for closing the dispatchers requests channel
	// otherwise we get a deadlock
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		c2.makeRequest(&d)
		wg.Done()
	}()

	go func() {
		c1.makeRequest(&d)
		wg.Done()
	}()

	go func() {
		c3.makeRequest(&d)
		wg.Done()
	}()

	wg.Wait()
	close(d.requests)
	d.broadcastRequestsToTaxis()
}
