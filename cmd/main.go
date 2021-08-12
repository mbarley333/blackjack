package main

import (
	"card"
	"fmt"
)

func main() {
	result := card.NewAceOrNothing()
	fmt.Println(result)
}
