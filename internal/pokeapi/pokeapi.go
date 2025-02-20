package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sniidu/pokedexcli/internal/pokecache"
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

func Map(c *shared.Config, back bool, cache *pokecache.Cache) {
	var url string
	if back {
		url = c.Previous
	} else {
		url = c.Next
	}

	if url == "" {
		fmt.Println("you're on the first page")
		return
	}

	var data []byte
	var err error

	data, found := cache.Get(url)

	if !found {
		fmt.Println("Cache miss!")
		data, err = fetch(url)
		if err != nil {
			fmt.Println(err)
			return
		}
		cache.Add(url, data)
	}

	locations, err := unmarsh(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	c.Next = locations.Next
	c.Previous = locations.Previous

	locations.printer()
}

func (l location) printer() {
	for _, loc := range l.Results {
		fmt.Println(loc.Name)
	}
}

func fetch(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET failed with %e", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("can't decode body to bytes as %e", err)
	}

	return data, nil
}

func unmarsh(data []byte) (location, error) {
	var locations location

	if err := json.Unmarshal(data, &locations); err != nil {
		return location{}, fmt.Errorf("can't unmarshal result due to %e", err)
	}
	return locations, nil
}
