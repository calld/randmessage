package main

import (
  "os"
  "log"
  "fmt"
  "math/rand"
  "encoding/json"
  "regexp"
  "strings"
  "time"
)

var dicfilename string

func init() {
  if len(os.Args) < 2 {
    log.Fatalln("missing dictionary filename argument")
  }
  dicfilename = os.Args[1]
}

func main() {
  dicfile, err := os.Open(dicfilename)
  if err != nil {
    log.Fatal(err.Error())
  }
  dec := json.NewDecoder(dicfile)
  gen := new(MessageGenerator)
  err = dec.Decode(gen)
  if err != nil {
    log.Fatal(err.Error())
  }
  fmt.Println(gen.Generate())
}

type MessageGenerator struct {
  start string
  replace map[string][]string
  die *rand.Rand
}

func (m *MessageGenerator) Generate() string {
  re, _ := regexp.Compile(`\{[^}]*\}`)
  current := m.start
  for replaces := re.FindAllStringIndex(current, -1); len(replaces) > 0; replaces = re.FindAllStringIndex(current, -1) {
    newc := new(strings.Builder)
    oldEnd := 0
    for _, pair := range replaces {
      newc.WriteString(current[oldEnd:pair[0]])
      newc.WriteString(m.randphrase(current[pair[0]+1:pair[1]-1]))
      oldEnd = pair[1]
    }
    newc.WriteString(current[oldEnd:])
    current = newc.String()
  }
  return current
}

func (m *MessageGenerator) randphrase(key string) string {
  return m.replace[key][m.die.Intn(len(m.replace[key]))]
}

func (m *MessageGenerator) UnmarshalJSON(data []byte) (err error) {
  defer func(){
    if r := recover(); r != nil {
      err = fmt.Errorf("%s", r)
    }
  }()
  var temp map[string]interface{}
  err = json.Unmarshal(data, &temp)
  if err != nil {
    return err
  }
  if m.replace == nil {
    m.replace = make(map[string][]string)
  }
  for k, v := range temp {
    if k == "__start__" {
      m.start = v.(string)
    } else {
      m.replace[k] = make([]string, 0)
      for _, item := range v.([]interface{}) {
        m.replace[k] = append(m.replace[k], item.(string))
      }
    }
  }
  m.die = rand.New(rand.NewSource(time.Now().Unix()))
  return nil
}

func (m *MessageGenerator) MarshalJSON() ([]byte, error) {
  temp := make(map[string]interface{})
  for k, v := range m.replace {
    temp[k] = v
  }
  temp["__start__"] = m.start
  return json.MarshalIndent(temp, "", "  ")
}
