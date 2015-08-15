package main

import (
	"flag"
	"html/template"
	"os"
	"io/ioutil"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"regexp"
	"fmt"
	"math/rand"
	"time"
	"github.com/rock/logic/handler"
)

var (
	addr = flag.Bool("addr", false, "find open address and print to final-port.txt")
	templates = template.Must(template.ParseFiles("pages/index.html", "pages/game.html", "pages/result.html"))
	validPath = regexp.MustCompile("^/(rpsls/game|rpsls/new|rpsls/result|rpsls/history)/([a-zA-Z0-9]+)$")
)

type Result struct {
	GameId string
	LastGame string
	PlayerWins int
	ComputerWins int
	Ties int
	Games int
	PlayerPercent float64
	ComputerPercent float64
	TiesPercent float64
}

func save(g *game.Game) error {
	filename := g.GameId + ".txt"
	b, _ := json.Marshal(g)

	return ioutil.WriteFile("data/" + g.GameId + "/" + filename, b, 0600)
}

func loadGame(gameId string) (*game.Game, error) {
 	filename := gameId + ".txt"
 	body, err := ioutil.ReadFile("data/" + gameId + "/" + filename)
 	if err != nil {
 		return nil, err
 	}
 	g := new(game.Game)
 	json.Unmarshal(body, g)
 	return g, nil
}

func generateGameId(l int) string {
	var bytes string
	for i:=0; i<l; i++ {
		bytes += fmt.Sprintf("%d", rand.Intn(100))
	}
	return string(bytes)
}

func renderTemplate(w http.ResponseWriter, tmpl string, g *game.Game) {
	err := templates.ExecuteTemplate(w, tmpl+".html", g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func frontPageHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", nil)
}

func gameHandler(w http.ResponseWriter, r *http.Request, gameId string) {
	g, err := loadGame(gameId)
	if err != nil {
		http.Redirect(w, r, "/rpsls/new/", http.StatusFound)
		return
	}
	renderTemplate(w, "game", g)
}

func newHandler(w http.ResponseWriter, r *http.Request, gameId string) {
	g := game.Game { GameId: gameId, Player: nil, Computer: nil, PlayerWins: 0, ComputerWins: 0, Ties: 0, Games: 0}
	os.Mkdir("data/" + gameId, 0600)
	err := save(&g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/rpsls/game/" + gameId, http.StatusFound)
}

func newHandler2(w http.ResponseWriter, r *http.Request) {
	g := game.Game { GameId: generateGameId(10), Player: nil, Computer: nil, PlayerWins: 0, ComputerWins: 0, Ties: 0, Games: 0}
	os.Mkdir("data/" + g.GameId, 0600)
	err := save(&g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, _ := json.Marshal(g)
	fmt.Fprint(w, string(b))
}

func resultHandler(w http.ResponseWriter, r *http.Request, gameId string) {
	g, err := loadGame(gameId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	p := r.FormValue("choice")
	g.RunRound(p)
	save(g)

	res := Result { GameId : g.GameId, LastGame : g.LastGame, PlayerWins : g.PlayerWins,
		ComputerWins : g.ComputerWins, Ties : g.Ties, Games : g.Games, TiesPercent : g.TiesPercent,
		PlayerPercent : g.PlayerPercent, ComputerPercent : g.ComputerPercent,
	}

	b, _ := json.Marshal(res)
	fmt.Fprint(w, string(b))
	//http.Redirect(w, r, "/rpsls/game/" + gameId, http.StatusFound)
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
	http.HandleFunc("/rpsls/", frontPageHandler)
	http.HandleFunc("/rpsls/new/", newHandler2)
	http.HandleFunc("/rpsls/game/", makeHandler(gameHandler))
	http.HandleFunc("/rpsls/result/", makeHandler(resultHandler))
	http.Handle("/rpsls/css/", http.StripPrefix("/rpsls/css/", http.FileServer(http.Dir("pages"))))
	http.Handle("/rpsls/js/", http.StripPrefix("/rpsls/js/", http.FileServer(http.Dir("pages/js"))))

//	if *addr {
		l, err := net.Listen("tcp", "127.0.0.1:9001")
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile("final.port.txt", []byte(l.Addr().String()), 0644)
		if err != nil {
			log.Fatal(err)
		}
		s := &http.Server{}
		s.Serve(l)
//		return
//	}

//	http.ListenAndServe(":9001", nil)
}