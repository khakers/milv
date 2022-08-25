package main

import (
	"fmt"
	"os"
	"time"

	"github.com/khakers/milv/cli"
	milv "github.com/khakers/milv/pkg"
)

func main() {
	start := time.Now()

	cliCommands := cli.ParseCommands()

	// pretty.Println(cliCommands)

	fmt.Printf("Parse time %+v\n", time.Since(start))
	start = time.Now()

	milv.SetBasePath(cliCommands.BasePath, false)

	fmt.Printf("SetBasePath time %+v\n", time.Since(start))
	start = time.Now()

	config, err := milv.NewConfig(cliCommands)
	if err != nil {
		panic(err)
	}

	fmt.Printf("NewConfig time %+v\n", time.Since(start))
	start = time.Now()

	files, _ := milv.NewFiles(cliCommands.Files, config)
	files.Run(cliCommands.Verbose)

	fmt.Printf("NewFiles time %+v\n", time.Since(start))

	if files.Summary() {
		os.Exit(1)
	}

	fmt.Println("NO ISSUES :-)")
}
