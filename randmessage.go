package randmessage

import (
  "fmt"
  "math/rand"
  "encoding/json"
  "regexp"
  "strings"
  "time"
)

type MessageGenerator struct {
  start string
  replace map[string][]string
  dice *rand.Rand
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
  defer func(){
    if r := recover(); r != nil {
      panic(fmt.Errorf("%s when attempting to get %s", r, key))
    }
  }()
  return m.replace[key][m.dice.Intn(len(m.replace[key]))]
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
      if len(m.start) > 0 {
        panic("multiple start strings, only one dictionary should have a non-empty __start__ field.")
      }
      m.start = v.(string)
    } else {
      if m.replace[k] == nil {
        m.replace[k] = make([]string, 0)
      }
      for _, item := range v.([]interface{}) {
        m.replace[k] = append(m.replace[k], item.(string))
      }
    }
  }
  if m.dice == nil {
    m.dice = rand.New(rand.NewSource(time.Now().Unix()))
  }
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
