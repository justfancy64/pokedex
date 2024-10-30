package main

import (
  "fmt"
  "bufio"
  "os"
  "net/http"
  "encoding/json"
  "io"
  "time"
  "github.com/justfancy64/pokedexcli/internal/pokecache"
)

func main() {
  var cfg config
  cache := pokecache.NewCache(5 * time.Second) // this is a pointer


  cfg.Next = "https://pokeapi.co/api/v2/location-area/"
  UI := createUI(&cfg, cache)
  fmt.Printf("pokedex> ") 

  scanner := bufio.NewScanner(os.Stdin)
  for scanner.Scan() {
    err := UI[scanner.Text()].callback(&cfg, cache) 
    fmt.Sprintf("pokedex> ")
    if err != nil {
      fmt.Sprintf("Error: %v", err)
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
  callback     func(cfg *config, c *pokecache.Cache) error
  
}
func commandHelp(cfg *config,c *pokecache.Cache) error {
  fmt.Println("")
  fmt.Println("Welcome to the pokedex\n")
  fmt.Printf("Usage:\n\nhelp: Displays the help menu\nexit: exits the pokedex\n")
  return nil
}

func commandExit(cfg *config,c *pokecache.Cache) error {
  fmt.Println("this command exits the program bai bai")
  os.Exit(0)
  return nil
}

func commandMap(cfg *config, c *pokecache.Cache) error {
  //api := "https://pokeapi.co/api/v2/location-area/" // poke api endpoint

  res, err := http.Get(cfg.Next)
  if err != nil {
    fmt.Println(err)
    return err
  }
  fmt.Println("request succesfull")
  defer res.Body.Close()


  dat,  err  := io.ReadAll(res.Body) // dat is []byte
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

func commandMapb(cfg *config,c *pokecache.Cache) error {
  var data []byte
  if cfg.Previous == "" {
    fmt.Println("no previous locations to display")
    return nil
  }
  chaceEntry, ok := c.Get(cfg.Previous)
  if !ok {
    res, err := http.Get(cfg.Previous)
    if err != nil {
      fmt.Println(err)
      return err
    }
    defer res.Body.Close()

    data, err = io.ReadAll(res.Body)
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
/*
func (c config)  updateNextPrev(next, previous string) {
  c.Next = next
  c.Previous = Previous
}
*/

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



func createUI(cfg *config, c *pokecache.Cache) map[string]cliCommand {
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
    }
  }
