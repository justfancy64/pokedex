package main


import (
  "fmt"
  "os"
  "encoding/json"
  "github.com/justfancy64/pokedexcli/internal/pokecache"
  "github.com/justfancy64/pokedexcli/internal/apireq"
  "github.com/justfancy64/pokedexcli/internal/poketypes"
  "math/rand"
)

type CliCommand struct {
  name         string
  description  string 
  callback     func(cfg *Config, c *pokecache.Cache, args []string) error
 
}


type Config struct {
  Next     string 
  Previous string 
}


type Locationresponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}


func CommandHelp(cfg *Config,c *pokecache.Cache, args []string) error {
  fmt.Println("")
  fmt.Println("Welcome to the pokedex\n")
  fmt.Printf("Usage:\n\nhelp: Displays the help menu\nexit: exits the pokedex\n")
  return nil
}

func CommandExit(cfg *Config,c *pokecache.Cache, args []string ) error {
  fmt.Println("this command exits the program bai bai")
  os.Exit(0)
  return nil
}

func CommandMap(cfg *Config, c *pokecache.Cache, args []string ) error {
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


  var areas Locationresponse
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

func CommandMapb(cfg *Config,c *pokecache.Cache, args []string) error {
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
  var areas Locationresponse
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
  
func CommandExplore(cfg *Config,c *pokecache.Cache, args []string) error {
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


func CommandCatch(cfg *Config, c *pokecache.Cache, args []string) error {
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
  
  
  if caught := Catch(pokemon.BaseExperience); caught {
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
  
  
  func Catch(baseexp int) bool {
    var catchchance float64
  
  
    catchchance = ((300 - float64(baseexp))/300) * 100
    random := rand.Float64()
  
    if (random * 100) < catchchance{
      return true
    } else {
      return false
    }
  
  }

func CommandInspect(cfg *Config, c *pokecache.Cache, args []string) error {
  poke, ok := c.GetPokemon(args[1])
  if !ok {
    fmt.Printf("you havent caught a %s yet\n", args[1])
    return nil
  }
  fmt.Printf("Name: %s \n", poke.Name)
  fmt.Printf("Height: %d \n", poke.Height)
  fmt.Printf("Weight: %d \n", poke.Weight)
  fmt.Println("Stats:")
  for _, val := range poke.Stats {
    fmt.Printf("  -%s: %d \n", val.Stat.Name, val.BaseStat)
  }
  fmt.Println("Types:")
  for _, val := range poke.Types {
    fmt.Printf("  - %s\n", val.Type.Name)
  }
  return nil
  
}
  
func CommandPokedex(cfg *Config, c *pokecache.Cache, args []string) error {
  fmt.Println("you have caught the following pokemon:")
  for key, _ := range c.Pokedex {
    fmt.Printf("  %s\n", key)
  }
  return nil
}




