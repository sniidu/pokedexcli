package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/sniidu/pokedexcli/internal/pokecache"
	"github.com/sniidu/pokedexcli/internal/shared"
)

type locationBundle struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type location struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func Explore(area string, c *shared.Config, cache *pokecache.Cache) {
	var data []byte
	var err error
	url := c.Next + area

	data, found := cache.Get(url)

	if !found {
		fmt.Println("Cache miss!")
		data, err = fetch(url)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	singleLocation, err := unmarshalLocation(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !found {
		cache.Add(url+strconv.Itoa(singleLocation.ID), data)
		cache.Add(url+singleLocation.Name, data)
	}

	fmt.Println("Exploring", singleLocation.Name, "...")
	fmt.Println("Found Pokemon:")
	for _, pokemon := range singleLocation.PokemonEncounters {
		fmt.Println(string('-'), pokemon.Pokemon.Name)
	}
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

	locations, err := unmarshalLocationBundle(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	c.Next = locations.Next
	c.Previous = locations.Previous

	for _, loc := range locations.Results {
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

func unmarshalLocation(data []byte) (location, error) {
	var singleLocation location

	if err := json.Unmarshal(data, &singleLocation); err != nil {
		return location{}, fmt.Errorf("can't unmarshal result due to %e", err)
	}
	return singleLocation, nil
}

func unmarshalLocationBundle(data []byte) (locationBundle, error) {
	var locations locationBundle

	if err := json.Unmarshal(data, &locations); err != nil {
		return locationBundle{}, fmt.Errorf("can't unmarshal result due to %e", err)
	}
	return locations, nil
}
