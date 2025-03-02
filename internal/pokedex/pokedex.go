package pokedex

import "fmt"

type Pokemon struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}

// Inspect main information of Pokemon
func Inspect(poke string, dex map[string]Pokemon) error {
	if inspectee, found := dex[poke]; found {
		fmt.Println("Name:", inspectee.Name)
		fmt.Println("Height:", inspectee.Height)
		fmt.Println("Weight:", inspectee.Weight)
		fmt.Println("Stats:")
		for _, val := range inspectee.Stats {
			fmt.Printf("  -%s: %d\n", val.Stat.Name, val.BaseStat)
		}
		fmt.Println("Types:")
		for _, val := range inspectee.Types {
			fmt.Println("  -", val.Type.Name)
		}
	} else {
		fmt.Println("you have not caught that pokemon")
	}
	return nil
}

// List all Pokemon caught
func List(dex map[string]Pokemon) error {
	if len(dex) == 0 {
		fmt.Println("Go catch some Pokemon!")
	} else {
		fmt.Println("Your Pokedex:")
		for _, poke := range dex {
			fmt.Printf("- %s\n", poke.Name)
		}
	}
	return nil
}
