package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var (
	Mapping = map[string]func(args ...string) (string, error){
		"run": runGame,
	}
)

type NewGameRequest struct {
	Width  int
	Height int
	Food   int
	Snakes []Snake
}

type Snake struct {
	Name string
	Url  string
}

func runGame(args ...string) (string, error) {
	fmt.Println("running game", args)
	client := &http.Client{Timeout: 5 * time.Second}
	data, err := json.Marshal(NewGameRequest{
		Width:  20,
		Height: 20,
		Food:   10,
		Snakes: []Snake{
			{Name: "Snake 1", Url: "https://dsnek.herokuapp.com"},
			{Name: "Snake 2", Url: "https://dsnek.herokuapp.com"},
		},
	})
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer(data)
	resp, err := client.Post(fmt.Sprintf("%s/games", "https://engine.battlesnake.io"), "application/json", buf)
	if err != nil {
		return "", err
	}
	res := map[string]string{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return "", err
	}
	if err := resp.Body.Close(); err != nil {
		return "", err
	}
	id := res["ID"]

	resp, err = client.Post(fmt.Sprintf("%s/games/%s/start", "https://engine.battlesnake.io", id), "application/json", nil)
	if err != nil {
		return "", err
	}
	err = resp.Body.Close()
	if err != nil {
		return "", err
	}
	return id, nil
}
