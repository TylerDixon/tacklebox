package main

import (
	"fmt"
)

func main() {
	configData, err := RetrieveConfig()
	if err != nil {
		return;
	}
	fmt.Printf("%s", configData)
}
