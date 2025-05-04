package main

import (
	"flag"
	"fmt"
)

func main() {
	help := flag.Bool("help", false, "show usage")
	port := flag.Int("port", 8080, "port number")
	flag.Parse()

	if *help {
		fmt.Println("Usage:\n  marketflow [--port <N>]\n  marketflow --help")
		return
	}

	fmt.Printf("Starting on :%d …\n", *port)
}
