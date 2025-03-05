package main

import (
	"flag"
	"log"

	"github.com/pocketix/pocketix-go/src/parser"
	"github.com/pocketix/pocketix-go/src/services"
)

func main() {
	path := flag.String("path", "programs/basic/empty_block.json", "path to the program file")
	flag.Parse()

	data := services.OpenFile("programs/" + *path)
	if _, err := parser.Parse(data); err != nil {
		log.Fatalln(err)
	}
	// else {
	// 	fmt.Println(program)
	// }
}
