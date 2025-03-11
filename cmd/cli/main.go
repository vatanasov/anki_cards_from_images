package main

import (
	"anki_cards_from_images"
	"flag"
	"fmt"
	"os"
)

func main() {
	questionsPath := flag.String("d", "", "path to questions file")
	flag.Parse()

	if *questionsPath == "" {
		fmt.Println("You must specify a path to a questions file")
		os.Exit(1)
	}

	err := anki_cards_from_images.Process(*questionsPath, os.Stdout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
