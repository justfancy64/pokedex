package main


import (
  "github.com/justfancy64/pokedexcli/internal/pokecache"
)



func CreateUI(cfg *Config, c *pokecache.Cache, args []string) map[string]CliCommand {
    return map[string]CliCommand{
      "help": {
        name:        "help",
        description: "Displays a help message",
        callback:    CommandHelp,
      }    ,
      "exit": {
        name:        "exit",
        description: "Exit the Pokedex",
        callback:    CommandExit,
      },
      "map": {
        name:        "map",
        description: "Displays 20 map locations",
        callback:    CommandMap,
      },
      "mapb": {
        name:        "mapb",
        description: "Displays the 20 previously displayed map locations",
        callback:    CommandMapb,
      },
      "explore": {
        name:        "explore",
        description: "lists pokemon of area",
        callback:    CommandExplore,
      },
      "catch": {
        name:        "catch",
        description: "attemps to catch pokemon",
        callback:    CommandCatch,
      },
      "inspect": {
        name:        "inspect",
        description: "lists information about pokemon in pokedex",
        callback:    CommandInspect,
      },
      "pokedex": {
        name:        "pokedex",
        description: "lists caught pokemon",
        callback:    CommandPokedex,
      },
    }
  }



