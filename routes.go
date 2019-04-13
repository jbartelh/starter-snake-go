package main

import (
	"log"
	"math"
	"net/http"

	"github.com/jbartelh/battlesnake-go/api"
)

const (
	up    = "up"
	down  = "down"
	left  = "left"
	right = "right"
)

var moves = [...]string{up, down, left, right}

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
	snakeHead := decoded.You.Body[0]
	log.Printf("snakehead: (%v,%v)", snakeHead.X, snakeHead.Y)

	possibleMoves := make(chan map[string]bool)
	go func(){
		possibleMoves <-validateAllMoves(&snakeHead, &decoded.Board)
	}()

	nearestFoodCoord := nearestFood(&snakeHead, decoded.Board.Food)
	log.Printf("nearestFoodCoord: (%v,%v)", nearestFoodCoord.X, nearestFoodCoord.Y)

	move := moveTowards(&snakeHead, nearestFoodCoord)
	//respond(res, validateMove(&snakeHead, move, decoded.Board.Snakes))
	pms := <-possibleMoves
	log.Printf("posible moves: %v", pms)
	if len(pms) == 0 {
		log.Printf("Game Lost!")
		respond(res, moveResponse("up"))
		return
	}
	m, ok := pms[move]
	if ok {
		log.Printf("move is fine %v", m)
		respond(res, moveResponse(move))
		return
	}
	for m := range pms {
		log.Printf("use alternative move")
		respond(res, moveResponse(m))
		return
	}

}

// checks for all four moves, if they are possible or not.
func validateAllMoves(pos *api.Coord, board *api.Board) map[string]bool {
	possibleMoves := make(map[string]bool)
	for _, m := range moves {
		if validateMove(pos, m, board) {
			possibleMoves[m] = true
		}
	}
	return possibleMoves
	//close(out)
}

// checks if the move is possible. Validates that the coordinate is not blocked by a snake and is inbound
func validateMove(pos *api.Coord, move string, board *api.Board) bool {
	moveCoord := moveToCoord(&move, pos)
	if coordOutOfBound(&moveCoord, board){
		log.Printf("move is out of bound: m(%v) coords(%v)", move, moveCoord)
		return false
	}
	for _, snake := range board.Snakes{
		for _, snakeBody := range snake.Body {
			if moveCoord == snakeBody {
				return false
			}
		}
	}
	log.Printf("move possible: m(%v) coords(%v)", move, moveCoord)
	return true
}


// checks if the given coordinates are inside the board
func coordOutOfBound(coord *api.Coord, board *api.Board) bool {
	switch {
	case coord.X < 0:
		return true
	case coord.Y < 0:
		return true
	case coord.X >= board.Width:
		return true
	case coord.Y >= board.Height:
		return true
	default:
		return false
	}
}

// converts a move (up|down|left|right) into its corresponding coordinate
func moveToCoord(move *string, pos *api.Coord) api.Coord {
	switch *move {
	case right:
		return api.Coord{pos.X + 1, pos.Y}
	case left:
		return api.Coord{pos.X - 1, pos.Y}
	case up:
		return api.Coord{pos.X, pos.Y - 1}
	case down:
		return api.Coord{pos.X, pos.Y + 1}
	default:
		panic("shall never happen")
	}
}

// returns the next move towards the given coordinate
func moveTowards(from, to *api.Coord) string {
	dis := distanceCoord(from, to)
	if math.Abs(float64(dis.X)) > math.Abs(float64(dis.Y)) {
		// X is greater -> move horizontally (left|right)
		if dis.X < 0 {
			return right
		} else {
			return left
		}
	} else {
		// Y is greater -> move vertically (up|down)
		if dis.Y < 0 {
			return down
		} else {
			return up
		}
	}
}

// returns the nearest of the given food to the given position
func nearestFood(pos *api.Coord, foods []api.Coord) *api.Coord {
	nearestFood := &(foods[0])
	smallestDistance := distance(pos, nearestFood)
	for _, food := range foods[1:] {
		if foodDistance := distance(pos, &food); foodDistance < smallestDistance {
			smallestDistance = foodDistance
			nearestFood = &food
		}
	}
	return nearestFood
}

// sqrt(a^2 + a^2), where a is abs(from.x - to.x) and b is same with y
func distance(from, to *api.Coord) float64 {
	return math.Sqrt(math.Pow(math.Abs(float64(from.X-to.X)), 2) + math.Pow(math.Abs(float64(from.Y-to.Y)), 2))
}

// substracts the "to" vector from the "from" vector elementwise
func distanceCoord(from, to *api.Coord) api.Coord {
	return api.Coord{from.X - to.X, from.Y - to.Y}
}

func End(res http.ResponseWriter, req *http.Request) {
	return
}

func Ping(res http.ResponseWriter, req *http.Request) {
	return
}
