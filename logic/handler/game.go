package game

import (
	"math"
	"math/rand"
)

var (
	choice = [5]string {"rock", "paper", "scissors", "lizard", "spock"}
	moves = [10]Move {
		{Player1: "rock", Player2: "scissors", Flavor: "crushes"},
		{Player1: "rock", Player2: "lizard", Flavor: "crushes"},
		{Player1: "paper", Player2: "rock", Flavor: "covers"},
		{Player1: "paper", Player2: "spock", Flavor: "disproves"},
		{Player1: "scissors", Player2: "paper", Flavor: "cuts"},
		{Player1: "scissors", Player2: "lizard", Flavor: "decapitates"},
		{Player1: "lizard", Player2: "spock", Flavor: "poisons"},
		{Player1: "lizard", Player2: "paper", Flavor: "eats"},
		{Player1: "spock", Player2: "rock", Flavor: "vaporizes"},
		{Player1: "spock", Player2: "scissors", Flavor: "smashes"},
	}
)

type Move struct {
	Player1 string
	Player2 string
	Flavor string
}

type Game struct {
	GameId string
	Player []string
	LastGame string
	Computer []string
	Flavor []string
	PlayerWins int
	PlayerPercent float64
	ComputerWins int
	ComputerPercent float64
	Ties int
	TiesPercent float64
	Games int
}

func (g *Game) RunRound(p string) {
	c := computerChoice()
	g.Player = append(g.Player, p)

	g.Computer = append(g.Computer, c)

	for i := 0; i < 10; i++ {
		if p == c {
			g.Ties++
			g.Flavor = append(g.Flavor, "Ties")
			g.LastGame = p + " Ties " + c
			break
		}
		switch {
			case moves[i].Player1 == p && moves[i].Player2 == c:
				g.PlayerWins++
				g.Flavor = append(g.Flavor, moves[i].Flavor)
				g.LastGame = moves[i].Player1 + " " + moves[i].Flavor + " " + moves[i].Player2 + ", Player Wins"
				break
			case moves[i].Player1 == c && moves[i].Player2 == p:
				g.ComputerWins++
				g.Flavor = append(g.Flavor, moves[i].Flavor)
				g.LastGame = moves[i].Player1 + " " + moves[i].Flavor + " " + moves[i].Player2 + ", Computer Wins"
				break
		}
	}
	g.Games++
	//g.PlayerPercent = calcPercent(g.PlayerWins, g.Games)
	//g.ComputerPercent = calcPercent(g.ComputerWins, g.Games)
	//g.TiesPercent = calcPercent(g.Ties, g.Games)
	g.calcPercent()
}

func (g *Game) calcPercent() {
	if(g.Games == 0) {
		return
	}

	g.PlayerPercent = math.Floor(float64(g.PlayerWins) / float64(g.Games) * 100)
	g.ComputerPercent = math.Floor(float64(g.ComputerWins) / float64(g.Games) * 100)
	g.TiesPercent = math.Floor(float64(g.Ties) / float64(g.Games) * 100)
}

func computerChoice() string {
	return choice[rand.Intn(100) % 5]
}