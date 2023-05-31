package main

import "os"

func main() {
	mode := os.Args[1]
	if mode != "full" && mode != "runtime" {
		panic("Incorrent mode selected, please use 'full' or 'runtime'")
	}

}
