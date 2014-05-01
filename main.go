package main

import (
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
	"fmt"
	"math/rand"
	"time"
	"encoding/json"
)

var (
	addr = flag.Bool("addr", false, "find open address and print to final-port.txt")
	templates = template.Must(template.ParseFiles("tmpl/index.html", "tmpl/game.html", "tmpl/result.html"))
	validPath = regexp.MustCompile("^/(game|new|result)/([a-zA-Z0-9]+)$")
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
	ComputerWins int
	Ties int
	Games int
}

func (g *Game) save() error {
	filename := g.GameId + ".txt"
	b, _ := json.Marshal(g)
	return ioutil.WriteFile("data/"+filename, b, 0600)
}

func loadGame(gameId string) (*Game, error) {
 	filename := gameId + ".txt"
 	body, err := ioutil.ReadFile("data/"+filename)
 	if err != nil {
 		return nil, err
 	}
 	var g Game
 	json.Unmarshal(body, &g)
 	return &g, nil
}

func generateGameId(l int) string {
	var bytes string
	for i:=0; i<l; i++ {
		bytes += fmt.Sprintf("%d", rand.Intn(100))
	}
	return string(bytes)
}

func computerChoice() string {
	return choice[rand.Intn(100) % 5]
}

func renderTemplate(w http.ResponseWriter, tmpl string, g *Game) {
	err := templates.ExecuteTemplate(w, tmpl+".html", g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func frontPageHandler(w http.ResponseWriter, r *http.Request) {
	g := &Game { GameId: generateGameId(10)}
	renderTemplate(w, "index", g)
}

func gameHandler(w http.ResponseWriter, r *http.Request, gameId string) {
	g, err := loadGame(gameId)
	if err != nil {
		http.Redirect(w, r, "/new/"+gameId, http.StatusFound)
		return
	}
	renderTemplate(w, "game", g)
}

func newHandler(w http.ResponseWriter, r *http.Request, gameId string) {
	g := &Game { GameId: gameId, Player: nil, Computer: nil, PlayerWins: 0, ComputerWins: 0, Ties: 0, Games: 0}
	err := g.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/game/"+gameId, http.StatusFound)
}

func resultHandler(w http.ResponseWriter, r *http.Request, gameId string) {
	p := r.FormValue("choice")
	c := computerChoice()
	g, err := loadGame(gameId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
	g.save()
	http.Redirect(w, r, "/game/"+gameId, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
	flag.Parse()
	rand.Seed(time.Now().UTC().UnixNano())
	http.HandleFunc("/", frontPageHandler)
	http.HandleFunc("/new/", makeHandler(newHandler))
	http.HandleFunc("/game/", makeHandler(gameHandler))
	http.HandleFunc("/result/", makeHandler(resultHandler))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("tmpl"))))

	if *addr {
		l, err := net.Listen("tcp", "81.229.174.91:0")
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile("final.port.txt", []byte(l.Addr().String()), 0644)
		if err != nil {
			log.Fatal(err)
		}
		s := &http.Server{}
		s.Serve(l)
		return
	}

	http.ListenAndServe(":8080", nil)
}