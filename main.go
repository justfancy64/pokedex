package main

import (
  "fmt"
  "bufio"
  "os"
  "encoding/json"
  "time"
  "strings"
  "github.com/justfancy64/pokedexcli/internal/pokecache"
  "github.com/justfancy64/pokedexcli/internal/apireq"
  "github.com/justfancy64/pokedexcli/internal/poketypes"
  "math/rand"
)

func main() {
  var cfg config
  cache := pokecache.NewCache(5 * time.Second) // this is a pointer


  cfg.Next = "https://pokeapi.co/api/v2/location-area/"
  fmt.Printf("pokedex> ") 

  scanner := bufio.NewScanner(os.Stdin)
  for scanner.Scan() {
    input := scanner.Text()
    args := strings.Fields(input)
    UI := createUI(&cfg, cache, args)
    _,ok := UI[args[0]]
    if !ok {
      fmt.Println("Command not found. type help")
      fmt.Printf("pokedex>")
      continue
    }
    err := UI[args[0]].callback(&cfg, cache, args) 
    fmt.Println("pokedex> ")
    if err != nil {
      fmt.Printf("Error: %v\n", err)
    }

  }
  if err := scanner.Err(); err != nil {
    fmt.Fprintln(os.Stderr, "reading standard input:", err)
  }
  
  return



}


type cliCommand struct {
  name         string
  description  string 
  callback     func(cfg *config, c *pokecache.Cache, args []string) error
 
}
func commandHelp(cfg *config,c *pokecache.Cache, args []string) error {
  fmt.Println("")
  fmt.Println("Welcome to the pokedex\n")
  fmt.Printf("Usage:\n\nhelp: Displays the help menu\nexit: exits the pokedex\n")
  return nil
}

func commandExit(cfg *config,c *pokecache.Cache, args []string ) error {
  fmt.Println("this command exits the program bai bai")
  os.Exit(0)
  return nil
}

func commandMap(cfg *config, c *pokecache.Cache, args []string ) error {
  //api := "https://pokeapi.co/api/v2/location-area/" // poke api endpoint
  
 

  dat, err  := apireq.Basicreq(cfg.Next) // dat is []byte
  if err != nil {
    fmt.Println(err)
  }
  err = c.Add(cfg.Next, dat)
  if err != nil {
    return err
  }

  // call to cacheAdd will be hereadd
  // ?


  var areas locationresponse
  err = json.Unmarshal(dat, &areas)
  if err != nil {
    fmt.Println(err)
  }
  var prev interface{} = areas.Previous

  s, ok := prev.(string)
  if ok {
    cfg.Previous = s
  }

  cfg.Next = areas.Next
  for _, area := range areas.Results {
    fmt.Println(area.Name)
  }


  return nil


}

func commandMapb(cfg *config,c *pokecache.Cache, args []string) error {
  var data []byte
  if cfg.Previous == "" {
    fmt.Println("no previous locations to display")
    return nil
  }
  chaceEntry, ok := c.Get(cfg.Previous)
  if !ok {
    data, err := apireq.Basicreq(cfg.Previous)
    if err != nil {
      fmt.Println(err)
      return err
    }

    
    c.Add(cfg.Previous, data)
  } else { 
    data = chaceEntry
  }
  var areas locationresponse
  err := json.Unmarshal(data, &areas)
  if err != nil {
    fmt.Println(err)
    return err
  }
  var prev interface{} = areas.Previous

  s, ok := prev.(string)
  if ok {
    cfg.Previous = s
  }

  cfg.Next = areas.Next  
  for _, area := range areas.Results {
  fmt.Println(area.Name)
  }
  return nil
}

func commandExplore(cfg *config,c *pokecache.Cache, args []string) error {
  apilink := "https://pokeapi.co/api/v2/location-area/" + args[1] + "/"

  var loc poketypes.SpecifiedLocation
  data, err := apireq.Basicreq(apilink)
  if err != nil {
    return fmt.Errorf("error in explore request: %v", err)
  } 

  err = json.Unmarshal(data, &loc)
  if err != nil {
    return fmt.Errorf("error in explore Unmarshal: %v",err)
  }
  fmt.Println("pokemons found in %s", args[1])
  for _, val := range loc.PokemonEncounters {
    fmt.Println(val.Pokemon.Name)
  }
  
  return nil


}


func commandCatch(cfg *config, c *pokecache.Cache, args []string) error {
  apilink := "https://pokeapi.co/api/v2/pokemon/" + args[1] + "/"
  var data []byte

  entry, ok := c.Get(apilink) // returns pokemon type, bool
  if !ok {
    var err error
    data, err = apireq.Basicreq(apilink)
    if err != nil {
      return err
    } 
    err = c.Add(apilink, data)
    if err != nil {
      return err
    }
    



  } else {
    data = entry
  }

  var pokemon poketypes.Pokemon
  err := json.Unmarshal(data, &pokemon)
    if err != nil {
      return fmt.Errorf("error in catch unmarshal %v", err)
    }


  if caught := catch(pokemon.BaseExperience); caught {
    fmt.Printf("%s was caught and added to the pokedex\n", pokemon.Name)
    err = c.AddPokemon(args[1], pokemon)
    if err != nil {
      return fmt.Errorf("error adding pokemon to cache: %v", err)
    }
  } else {
    fmt.Printf("failed to catch %s ", pokemon.Name)
  }
  return nil
}
 

func catch(baseexp int) bool {
  var catchchance float64


  catchchance = ((300 - float64(baseexp))/300) * 100
  random := rand.Float64()

  if (random * 100) < catchchance{
    return true
  } else {
    return false
  }

}



type config struct {
  Next     string 
  Previous string 
}


type locationresponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}



func createUI(cfg *config, c *pokecache.Cache, args []string) map[string]cliCommand {
    return map[string]cliCommand{
      "help": {
        name:        "help",
        description: "Displays a help message",
        callback:    commandHelp,
      }    ,
      "exit": {
        name:        "exit",
        description: "Exit the Pokedex",
        callback:    commandExit,
      },
      "map": {
        name:        "map",
        description: "Displays 20 map locations",
        callback:    commandMap,
      },
      "mapb": {
        name:        "mapb",
        description: "Displays the 20 previously displayed map locations",
        callback:    commandMapb,
      },
      "explore": {
        name:        "explore",
        description: "lists pokemon of area",
        callback:    commandExplore,
      },
      "catch": {
        name:        "catch",
        description:  "attemps to catch pokemon",
        callback:    commandCatch,
      },
    }
  }
