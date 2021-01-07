package main

import "math"

// point in a 2d
type Location struct {
	X float64
	Y float64
}

// calculate the distance from current location to the other location
func (l *Location) DistanceTo(o Location) float64 {
	return math.Sqrt(math.Pow(l.X - o.X, 2) + math.Pow(l.Y - o.Y, 2))
}
