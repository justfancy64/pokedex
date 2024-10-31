package apireq

import (
  "fmt"
  "io"
  "net/http"
)







func Basicreq(apilink string) ([]byte, error) {
  res, err := http.Get(apilink)
  if err != nil {
    fmt.Println(err)
    return nil, err
  }
  defer res.Body.Close()

  data, err := io.ReadAll(res.Body)
  if err != nil {
    fmt.Println(err)
    return nil, err
  }
  return data, nil
}
