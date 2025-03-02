package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/sniidu/pokedexcli/internal/pokecache"
	"github.com/sniidu/pokedexcli/internal/pokedex"
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

// Create random number between 0 and maxBaseExperience
// If result is greater than experience, caught was successfull
func caught(experience int) bool {
	maxBaseExperience := 300
	return rand.Intn(maxBaseExperience) > experience
}

func Catch(poke string, c *shared.Config, cache *pokecache.Cache, dex map[string]pokedex.Pokemon) error {
	var data []byte
	var err error
	url := c.Next + poke

	data, found := cache.Get(url)

	// Fetch if not in cache
	if !found {
		data, err = fetch(url)
		if err != nil {
			return fmt.Errorf("error while fetching: %w", err)
		}
		cache.Add(url, data)
	}

	// Unmarshal data either from cache or fetch
	pokemonResponse, err := unmarshal[pokedex.Pokemon](data)
	if err != nil {
		return fmt.Errorf("error unmarshaling: %w", err)
	}

	fmt.Println("catch", poke)
	fmt.Printf("Throwing a Pokeball at %s...\n", poke)
	if caught(pokemonResponse.BaseExperience) {
		fmt.Println(poke, "was caught!")
		dex[poke] = pokemonResponse
		fmt.Println("You may now inspect it with the inspect command.")
	} else {
		fmt.Println(poke, "escaped!")
	}

	return nil
}

// Prints Pokemon found from provided area
// Area can be either name or id of said area
func Explore(area string, c *shared.Config, cache *pokecache.Cache) error {
	var data []byte
	var err error
	url := c.Next + area

	data, found := cache.Get(url)

	// Fetch if not in cache
	if !found {
		data, err = fetch(url)
		if err != nil {
			return fmt.Errorf("error while fetching: %w", err)
		}
	}

	// Unmarshal data either from cache or fetch
	singleLocation, err := unmarshal[location](data)
	if err != nil {
		return fmt.Errorf("error unmarshaling: %w", err)
	}

	// Add cache as id and name known latest here
	if !found {
		// Integer to string as plain string(id) would rune it
		cache.Add(url+strconv.Itoa(singleLocation.ID), data)
		cache.Add(url+singleLocation.Name, data)
	}

	fmt.Println("Exploring", singleLocation.Name, "...")
	fmt.Println("Found Pokemon:")
	for _, pokemon := range singleLocation.PokemonEncounters {
		fmt.Println(string('-'), pokemon.Pokemon.Name)
	}

	return nil
}

// Traverses through location areas in Pokemon and prints names of areas
// Keeps track of page and can move backwards and forwards
func Map(c *shared.Config, back bool, cache *pokecache.Cache) error {
	var url string
	if back {
		url = c.Previous
	} else {
		url = c.Next
	}

	if url == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	var data []byte
	var err error

	data, found := cache.Get(url)

	if !found {
		data, err = fetch(url)
		if err != nil {
			return fmt.Errorf("error while fetching: %w", err)
		}
		cache.Add(url, data)
	}

	locations, err := unmarshal[locationBundle](data)
	if err != nil {
		return fmt.Errorf("error unmarshaling: %w", err)
	}

	c.Next = locations.Next
	c.Previous = locations.Previous

	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}

	return nil
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

// Generics introduction
func unmarshal[T any](data []byte) (T, error) {
	var decoded T

	if err := json.Unmarshal(data, &decoded); err != nil {
		return decoded, fmt.Errorf("can't unmarshal result due to %e", err)
	}
	return decoded, nil
}
