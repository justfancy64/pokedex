package pokecache


import (
  "time"
  "sync"
  "github.com/justfancy64/pokedexcli/internal/poketypes"
)

type Cache struct {
  Entries       map[string]cacheEntry
  mu            *sync.Mutex
  Interval      time.Duration
  Pokedex       map[string]poketypes.Pokemon


}

type cacheEntry struct {
  CreatedAt time.Time 
  Val       []byte
}




func (c *Cache) Add(key string,val []byte) error {

  t0 := time.Now()

  entry := cacheEntry{
    CreatedAt: t0,
    Val:       val,
  } 

  c.mu.Lock()
  

  c.Entries[key] = entry
  c.mu.Unlock()

  return nil
}

func (c *Cache)AddPokemon(key string, val poketypes.Pokemon) error {
  c.mu.Lock()

  c.Pokedex[key] = val
  c.mu.Unlock()

  return nil

}



func (c *Cache) Get(key string) ([]byte, bool) {

  data, ok := c.Entries[key]
  if !ok {
    return nil, false
  }
  return data.Val, true
  
}
func (c *Cache) GetPokemon(key string) (poketypes.Pokemon, bool) {

  pokemon, ok := c.Pokedex[key]
  if !ok {
    return poketypes.Pokemon{}, false
  }
  return pokemon, true
}

func (c *Cache) reapLoop(interval time.Duration) {

  ticker := time.NewTicker(interval)
  go func() {
    for {
      select{
      case <-ticker.C:

        c.mu.Lock()
        t0 := time.Now()
        for key, entry := range c.Entries {
          if t0.Sub(entry.CreatedAt) > (interval) {
            delete(c.Entries, key )
          }
        c.mu.Unlock()
        }
      }
    }


  }()
}


func NewCache(interval time.Duration) *Cache {

  maps := make(map[string]cacheEntry)
  pokes := make(map[string]poketypes.Pokemon)

  return &Cache{
    Entries:  maps,
    Interval: interval,
    mu:       &sync.Mutex{},
    Pokedex:  pokes,

  }

}

