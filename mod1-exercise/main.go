package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	content := flag.String("content", "", "Content to write to the file")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s: \n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "%s --content \"Hello World\" filename.txt \n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Error: Filename not specified")
		flag.Usage()
		os.Exit(1)
	}

	filename := args[0]

	if _, err := os.Stat(filename); err == nil {
		log.Fatalf("Error: File %s already exists.", filename)
	} else if !os.IsNotExist(err) {
		log.Fatalf("Error checking file: %v", err)
	}

	file, err := os.Create(filename)
	if err!= nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(*content); err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}

	fmt.Printf("File %s created successfully.\n", filename)
}
