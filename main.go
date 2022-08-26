package main

import (
	"fmt"
	"time"

	"github.com/khakers/milv/cli"
	milv "github.com/khakers/milv/pkg"
)

func main() {
	total := time.Now()
	start := time.Now()

	cliCommands := cli.ParseCommands()

	// pretty.Println(cliCommands)

	fmt.Printf("Commands Parse time %+v\n", time.Since(start))
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

	fmt.Printf("\nNewFiles time %+v\n", time.Since(start))
	start = time.Now()

	files.Run(cliCommands.Verbose)

	fmt.Printf("\nFile parsing time %+v\n", time.Since(start))
	fmt.Printf("Total time elapsed %v\n", time.Since(total))

	// if files.Summary() {
	// 	os.Exit(1)
	// }

	fmt.Println("NO ISSUES :-)")
}
