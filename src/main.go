package main

import (
	"fmt"
)

func main() {
	trends := GetTrends()
	fmt.Println(PrettyPrint(trends))
}
