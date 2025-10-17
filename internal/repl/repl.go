package repl

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/allscorpion/pokedexcli/internal/config"
	"github.com/allscorpion/pokedexcli/internal/pokecache"
	"github.com/allscorpion/pokedexcli/internal/pokedex"
)


type cliCommand struct {
	name        string
	description string
	callback    func(c *config.Config, params ...string) error
}


func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:         "help",
			description: "Displays a help message",
			callback: commandHelp,
		},
		"map": {
			name: "map",
			description: "Get the next page of locations",
			callback: commandMap,
		},
		"mapb": {
			name: "mapb",
			description: "Get the previous page of locations",
			callback: commandMapBack,
		},
		"explore": {
			name: "explore",
			description: "Explore a location area by its name",
			callback: commandExplore,
		},
		"catch": {
			name: "catch",
			description: "Catch a Pokemon by its name",
			callback: commandCatch,
		},
		"inspect": {
			name: "inspect",
			description: "Inspect a Pokemon by its name",
			callback: commandInspect,
		},
		"pokedex": {
			name: "pokdex",
			description: "Check all the pokemon you have caught",
			callback: commandPokedex,
		},
	}
}


func commandExit(c *config.Config, params ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!");
	os.Exit(0);
	return nil;
}

func commandHelp(c *config.Config, params ...string) error {
	fmt.Println("Welcome to the Pokedex!");
	fmt.Print("Usage:\n\n");

	commands := getCommands();

	for cmd := range commands {
		fmt.Printf("%v: %v\n", cmd, commands[cmd].description);
	}

	return nil;
}

func commandMapInner(config *config.Config, url string) error {
	locationArea, err := pokedex.GetLocationAreas(config, url);
	if err != nil {
		return err
	}
	for _, area := range locationArea.Results {
		fmt.Println(area.Name);
	}
	config.Next = locationArea.Next;
	if locationArea.Previous != nil {
		config.Previous = locationArea.Previous.(string);
	}
	return nil;
}

func commandMap(config *config.Config, params ...string) error {
	return commandMapInner(config, config.Next);
}

func commandMapBack(config *config.Config, params ...string) error {
	if config.Previous == "" {
		return fmt.Errorf("no previous page available")
	}
	return commandMapInner(config, config.Previous);
}

func commandExplore(config *config.Config, params ...string) error {
	if len(params) == 0 {
		return fmt.Errorf("please provide a location area ID to explore");
	}

	name := params[0];

	locationArea, err := pokedex.GetLocationArea(config, "https://pokeapi.co/api/v2/location-area/" + name);

	if err != nil {
		return err
	}

	fmt.Printf("Exploring %v...\n", name);
	fmt.Println("Found Pokemon:");

	for _, v := range locationArea.PokemonEncounters {
		fmt.Println(" - " + v.Pokemon.Name);
	}

	return nil;
}

func commandCatch(config *config.Config, params ...string) error {
	if len(params) == 0 {
		return fmt.Errorf("please provide a pokemon name to catch");
	}

	
	name := params[0];

	_, alreadyCaught := config.Pokedex[name];

	if alreadyCaught {
		return fmt.Errorf("%v is already caught", name);
	}

	fmt.Printf("Throwing a Pokeball at %v...\n", name);

	pokemonDetails, err := pokedex.GetPokemonByName(config, name);

	if err != nil {
		return err
	}

	var randomChanceToCatch int = rand.Intn(pokemonDetails.BaseExperience);

	if randomChanceToCatch > 50 {
		fmt.Printf("Oh no! %v escaped the Pokeball!\n", name);
		return nil;
	}

	fmt.Printf("Gotcha! %v was caught successfully!\n", name);
	config.Pokedex[name] = pokemonDetails

	return nil;
}

func commandInspect(config *config.Config, params ...string) error {
	if len(params) == 0 {
		return fmt.Errorf("please provide a pokemon name to inspect");
	}

	name := params[0];

	pokemonDetails, exists := config.Pokedex[name];

	if !exists {
		return fmt.Errorf("%v has not been caught", name)
	}

	fmt.Printf("Name: %v\n", pokemonDetails.Name)
	fmt.Printf("Height: %v\n", pokemonDetails.Height)
	fmt.Printf("Weight: %v\n", pokemonDetails.Weight)
	fmt.Printf("Stats: \n")
	for _, stat := range pokemonDetails.Stats {
		fmt.Printf("  -%v: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Printf("Types: \n")
	for _, t := range pokemonDetails.Types {
		fmt.Printf("  - %v\n", t.Type.Name)
	}
	
	
	return nil;
}

func commandPokedex(config *config.Config, params ...string) error {
	if len(config.Pokedex) == 0 {
		return fmt.Errorf("you have not caught any pokemon")
	}

	fmt.Println("Your Pokedex:")
	for _, v := range config.Pokedex {
		fmt.Printf(" - %v\n", v.Name)
	}

	return nil;
}

func StartRepl() {
	reader := bufio.NewScanner(os.Stdin)

	supportedCommands := getCommands();
	cache := pokecache.NewCache(5 * time.Minute);

	config := &config.Config{
		Next: "https://pokeapi.co/api/v2/location-area",
		Previous: "",
		Cache: cache,
		Pokedex: map[string]config.Pokemon{},
	}

	for {
		fmt.Print("Pokedex > ")
		reader.Scan()

		words := cleanInput(reader.Text())
		if len(words) == 0 {
			continue
		}

		firstWord := words[0]
		restOfWords := words[1:]
		cmd, foundCmd := supportedCommands[firstWord];

		if !foundCmd {
			fmt.Println("Unknown command")
			continue;
		}

		err := cmd.callback(config, restOfWords...);
		if err != nil {
			fmt.Printf("%v\n", err);
		}
	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}