package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/beevik/guid"
)

// Question é uma estrutura de dados para
//dar suporte ao programa
type Question struct {
	Q string
	A string
}

//Game é a estrutura básica de controle do nosso jogo
type Game struct {
	Questions map[string]*Question
	Score     int
	Size      int
	Arquivo   string
}

// New starts a new game definition
func (g *Game) New(filename string) *Game {
	filename = strings.TrimSpace(filename)

	g.File(filename).Dict()
	return g
}

//Start sets the Game strucuture and runs a new game
func (g *Game) Start() *Game {
	for _, v := range g.Questions {

		fmt.Printf("What is %s? ", v.Q)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		answer := scanner.Text()
		strings.TrimSpace(answer)

		if answer == v.A {
			g.Score++
		}
	}

	return g
}

//Dict builds a map of questions
func (g *Game) Dict() *Game {

	file, err := os.Open(g.Arquivo)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	defer file.Close()

	csvReader := csv.NewReader(file)
	questions := make(map[string]*Question)

	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		uuid := guid.NewString()
		question := &Question{line[0], line[1]}

		questions[uuid] = question

	}
	g.Questions = questions
	return g

}

//File sets the file to be used in the Score
//problems.csv is the default if no other filepaths are provided
func (g *Game) File(filename string) *Game {

	if isValid, _ := regexp.MatchString(`[a-zA-Z0-9]{3,}\.csv`, filename); isValid {
		g.Arquivo = filename

	} else {
		g.Arquivo = "problems.csv"
	}
	return g

}

//ShowScore Prints the score of the player
func (g *Game) ShowScore() {
	fmt.Printf("Your score is %d out of %d\n", g.Score, len(g.Questions))
}

func main() {

	file := flag.String("file", "", "Path of the file to use in the game")
	timer := flag.String("timer", "30", "Define the duration in seconds of the game. Default:30s")
	flag.Parse()

	duration, err := time.ParseDuration(*timer)
	if err != nil {
		duration = 30 * time.Second
	}
	game := &Game{}
	t := time.AfterFunc(duration, func() {
		fmt.Println()
		game.ShowScore()
		os.Exit(1)
	})
	defer t.Stop()

	game.New(*file)
	//fmt.Printf("Running the file: %s\n", game.Arquivo)
	game.Start()

}
