package main

import (
	"fmt"
	"os"
)

func main() {
	filePath := os.Args[1]
	fmt.Printf("generate src file path =====> %v", filePath)
	General(filePath)
}
