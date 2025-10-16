package repl

import (
	"bufio"
	"fmt"
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
	callback    func(c *config.Config) error
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
	}
}


func commandExit(c *config.Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!");
	os.Exit(0);
	return nil;
}

func commandHelp(c *config.Config) error {
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

func commandMap(config *config.Config) error {
	return commandMapInner(config, config.Next);
}

func commandMapBack(config *config.Config) error {
	if config.Previous == "" {
		return fmt.Errorf("no previous page available\n")
	}
	return commandMapInner(config, config.Previous);
}


func StartRepl() {
	reader := bufio.NewScanner(os.Stdin)

	supportedCommands := getCommands();
	cache := pokecache.NewCache(5 * time.Minute);

	config := &config.Config{
		Next: "https://pokeapi.co/api/v2/location-area",
		Previous: "",
		Cache: cache,
	}

	for {
		fmt.Print("Pokedex > ")
		reader.Scan()

		words := cleanInput(reader.Text())
		if len(words) == 0 {
			continue
		}

		cmd, foundCmd := supportedCommands[words[0]];

		if !foundCmd {
			fmt.Println("Unknown command")
			continue;
		}

		err := cmd.callback(config);
		if err != nil {
			fmt.Print(err)
		}
	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}