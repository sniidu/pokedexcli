package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sniidu/pokedexcli/internal/shared"
)

type location struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func Map(c *shared.Config) {
	locations, err := fetch(c.Next)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(locations.Count)
	fmt.Println(c.Next)
	c.Next = locations.Next
	c.Previous = locations.Previous
	fmt.Println(c.Next)
}

func (l location) printer() {
	for _, loc := range l.Results {
		fmt.Println(loc.Name)
	}
}

func fetch(url string) (location, error) {
	res, err := http.Get(url)
	if err != nil {
		return location{}, fmt.Errorf("GET failed with %e", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return location{}, fmt.Errorf("can't decode body to bytes as %e", err)
	}

	var locations location

	if err = json.Unmarshal(data, &locations); err != nil {
		return location{}, fmt.Errorf("can't unmarshal result due to %e", err)
	}
	return locations, nil
}
