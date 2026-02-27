package test

import (
	"encoding/json"
	"fmt"
)

type User struct {
	Name string
	Age  int
}

func main() {
	u := User{"John", 30}
	data, _ := json.Marshal(u)
	fmt.Println(string(data))
}
