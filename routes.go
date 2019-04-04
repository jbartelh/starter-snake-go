package main

import (
	"log"
	"math"
	"net/http"

	"github.com/jbartelh/battlesnake-go/api"
)

const (
	up = "up"
	down = "down"
	left = "left"
	right = "right"
)
// creates a moveResponse out of a move (up|down|left|right)
func moveResponse(move string) api.MoveResponse {
	log.Printf("move: " + move)
	return api.MoveResponse{
		Move: move,
	}
}

func Index(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Battlesnake documentation can be found at <a href=\"https://docs.battlesnake.io\">https://docs.battlesnake.io</a>."))
}

func Start(res http.ResponseWriter, req *http.Request) {
	decoded := api.SnakeRequest{}
	err := api.DecodeSnakeRequest(req, &decoded)
	if err != nil {
		log.Printf("Bad start request: %v", err)
	}
	dump(decoded)

	respond(res, api.StartResponse{
		Color: "#75CEDD",
	})
}

func Move(res http.ResponseWriter, req *http.Request) {
	decoded := api.SnakeRequest{}
	err := api.DecodeSnakeRequest(req, &decoded)
	if err != nil {
		log.Printf("Bad move request: %v", err)
	}
	//dump(decoded)
	snakeHead := decoded.You.Body[0]
	log.Printf("snakehead: (%v,%v)", snakeHead.X, snakeHead.Y)
	nearestFoodCoord := nearestFood(snakeHead, decoded.Board.Food)
	log.Printf("nearestFoodCoord: (%v,%v)", nearestFoodCoord.X, nearestFoodCoord.Y)
	respond(res, moveTowards(snakeHead, nearestFoodCoord))
}

func moveTowards(from, to api.Coord) api.MoveResponse {
	dis := distanceCoord(from, to)
	if math.Abs(float64(dis.X)) > math.Abs(float64(dis.Y)) {
		// X is greater -> move horizontally (left|right)
		if dis.X < 0 {
			return moveResponse(right)
		} else {
			return moveResponse(left)
		}
	} else {
		// Y is greater -> move vertically (up|down)
		if dis.Y < 0 {
			return moveResponse(down)
		} else {
			return moveResponse(up)
		}
	}

}

// returns the nearest of the given food to the given position
func nearestFood(pos api.Coord, coord []api.Coord) api.Coord {
	nearestFood :=  coord[0]
	smallestDistance := distance(pos, nearestFood)
	for _, food := range coord[1:] {
		if foodDistance := distance(pos, food); foodDistance < smallestDistance {
			smallestDistance = foodDistance
			nearestFood = food
		}
	}
	return nearestFood
}

// sqrt(a^2 + a^2), where a is abs(from.x - to.x) and b is same with y
func distance(from, to api.Coord) float64 {
	return math.Sqrt(math.Pow(math.Abs(float64(from.X - to.X)), 2) + math.Pow(math.Abs(float64(from.Y-to.Y)), 2))
}

func distanceCoord(from, to api.Coord) api.Coord {
	return api.Coord{from.X - to.X, from.Y - to.Y}
}

func End(res http.ResponseWriter, req *http.Request) {
	return
}

func Ping(res http.ResponseWriter, req *http.Request) {
	return
}
