package main

import (
  "os"
  "log"
  "fmt"
  "flag"
  "encoding/json"
  "github.com/calld/randmessage"
)

var dicfilenames []string

var count = flag.Int("size", 1, "number of generated messages")

func init() {
  flag.Parse()
  if len(flag.Args()) < 1 {
    log.Fatalln("missing dictionary filename argument")
  }
  dicfilenames = make([]string, 0, len(flag.Args()))
  for _, name := range flag.Args() {
    dicfilenames = append(dicfilenames, name)
  }
}

func main() {
  gen := new(randmessage.MessageGenerator)
  for _, dicfilename := range dicfilenames {
    dicfile, err := os.Open(dicfilename)
    if err != nil {
      log.Fatal(err.Error())
    }
    dec := json.NewDecoder(dicfile)
    err = dec.Decode(gen)
    if err != nil {
      log.Fatal(err.Error())
    }
    dicfile.Close()
  }
  for i := 0; i < *count; i++ {
    fmt.Println(gen.Generate())
  }
}
