package main

import (
  "fmt"
  "bufio"
  "os"
  "time"
  "strings"
  "github.com/justfancy64/pokedexcli/internal/pokecache"
)

func main() {
  var cfg Config
  cache := pokecache.NewCache(5 * time.Second) // this is a pointer
  // cfg, cache := startup()

  cfg.Next = "https://pokeapi.co/api/v2/location-area/"
  Startup()
  fmt.Printf("pokedex> ") 
  // reaploop()
  scanner := bufio.NewScanner(os.Stdin)
  for scanner.Scan() {
    input := scanner.Text()
    args := strings.Fields(input)
    UI := CreateUI(&cfg, cache, args)
    cfg.Commands = UI
    _,ok := UI[args[0]]
    if !ok {
      fmt.Println("Command not found. type help")
      fmt.Printf("pokedex>")
      continue
    }
    err := UI[args[0]].callback(&cfg, cache, args) 
    fmt.Printf("pokedex> ")
    if err != nil {
      fmt.Printf("Error: %v\n", err)
    }

  }
  if err := scanner.Err(); err != nil {
    fmt.Fprintln(os.Stderr, "reading standard input:", err)
  }
  
  return



}



