package cli

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	gitignore "github.com/sabhiram/go-gitignore"
)

type Commands struct {
	BasePath        string
	ConfigFile      string
	IgnoreList      gitignore.GitIgnore
	Files           []string
	WhiteListExt    []string
	WhiteListInt    []string
	BlackList       []string
	Timeout         int
	ReguestRepeats  int8
	AllowRedirect   bool
	AllowCodeBlocks bool
	IgnoreExternal  bool
	IgnoreInternal  bool
	Verbose         bool
	FlagsSet        map[string]bool
}

func ParseCommands() Commands {
	basePath := flag.String("base-path", "", "The root source directories used to search for files")
	configFile := flag.String("config-file", "milv.config.yaml", "The config file for bot")
	ignoreFile := flag.String("ignore-file", "/.milvignore", "The .milvignore file")
	whiteListExt := flag.String("white-list-ext", "", "The white list external links")
	whiteListInt := flag.String("white-list-int", "", "The white list internal links")
	blackList := flag.String("black-list", "", "The files black list")
	timeout := flag.Int("timeout", 0, "Timeout for http.get reguest")
	requestRepeats := flag.Int("request-repeats", 0, "Times reguest failuring links")
	allowRedirect := flag.Bool("allow-redirect", false, "Allow redirect")
	allowCodeBlocks := flag.Bool("allow-code-blocks", false, "Allow links in code blocks to check")
	ignoreInternal := flag.Bool("ignore-internal", false, "Ignore internal links")
	ignoreExternal := flag.Bool("ignore-external", false, "Ignore external links")
	verbose := flag.Bool("v", false, "Enable verbose logging")

	flag.Parse()
	files := flag.Args()

	flagset := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		flagset[f.Name] = true
	})

	if *basePath != "" {
		*configFile = fmt.Sprintf("%s/%s", *basePath, *configFile)
	}

	ignoreList, err := readIgnoreFile(*basePath + *ignoreFile)
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()

	if len(files) == 0 {
		files, err = collectFiles(*basePath, *ignoreList)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("collectFiles time %+v\n", time.Since(start))
	// start = time.Now()

	// if len(files) == 0 {
	// 	out := runCmd("find . -name \"*.md\"", true)
	// 	files = strings.Split(string(out), "\n")
	// 	if len(files) > 0 {
	// 		files = files[:len(files)-1]
	// 	}
	// }
	// fmt.Printf("find time %+v\n", time.Since(start))

	// pretty.Println(collected)
	// pretty.Println(files)

	// pretty.Pdiff(log.Default(), collected, files)
	// pretty.Println(pretty.Diff(collected, files))

	return Commands{
		BasePath:        *basePath,
		ConfigFile:      *configFile,
		IgnoreList:      *ignoreList,
		Files:           files,
		WhiteListExt:    strings.Split(*whiteListExt, ","),
		WhiteListInt:    strings.Split(*whiteListInt, ","),
		BlackList:       strings.Split(*blackList, ","),
		Timeout:         *timeout,
		ReguestRepeats:  int8(*requestRepeats),
		AllowRedirect:   *allowRedirect,
		AllowCodeBlocks: *allowCodeBlocks,
		IgnoreExternal:  *ignoreExternal,
		IgnoreInternal:  *ignoreInternal,
		Verbose:         *verbose,
		FlagsSet:        flagset,
	}
}

func readIgnoreFile(path string) (*gitignore.GitIgnore, error) {
	return gitignore.CompileIgnoreFile(path)
}

func collectFiles(basepath string, ignore gitignore.GitIgnore) ([]string, error) {
	files := []string{}
	log.Println(basepath)

	err := filepath.WalkDir(basepath, func(path string, d fs.DirEntry, err error) error {
		if !ignore.MatchesPath(path) {
			if !d.IsDir() {
				i := strings.LastIndexAny(path, ".")
				// fmt.Println(path[i:])
				if i != -1 && path[i:] == ".md" {
					files = append(files, path)
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

// Deprecated
func runCmd(cmd string, shell bool) []byte {
	if shell {
		out, err := exec.Command("/bin/bash", "-c", cmd).Output()
		if err != nil {
			log.Fatal(err)
		}
		return out
	}
	out, err := exec.Command(cmd).Output()
	if err != nil {
		log.Fatal(err)
	}
	return out
}
