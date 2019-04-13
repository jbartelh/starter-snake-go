package main

import (
	"github.com/jbartelh/battlesnake-go/api"
	"github.com/stretchr/testify/assert"
	"testing"
)

var distanceTestCases = []struct {
	p1x      int     // Point 1: x-coord
	p1y      int     // Point 1: y-coord
	p2x      int     // Point 2: x-coord
	p2y      int     // Point 2: y-coord
	expected float64 // expected distance
}{
	{2, 1, 4, 5, 4.47213},
	{0, 0, 0, 2, 2},
	{0, 0, 2, 0, 2},
	{4,3,9,5,5.385},
}

const epsilon float64 = 0.001

func TestDistanceBetweenTwoCoords(t *testing.T) {
	for _, tt := range distanceTestCases {
		a := api.Coord{tt.p1x, tt.p1y}
		b := api.Coord{tt.p2x, tt.p2y}

		assert.InEpsilon(t, tt.expected, distance(&a, &b), epsilon)
	}
}

func TestDistanceBetweenTwoCoordsSwapped(t *testing.T) {
	for _, tt := range distanceTestCases {
		a := api.Coord{tt.p1x, tt.p1y}
		b := api.Coord{tt.p2x, tt.p2y}

		assert.InEpsilon(t, tt.expected, distance(&b, &a), epsilon)
	}
}
