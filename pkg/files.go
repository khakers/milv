package pkg

import (
	"log"
	"sync"

	"github.com/panjf2000/ants/v2"
	"github.com/schollz/progressbar/v3"
)

type Files []*File

type taskFunc func()

func NewFiles(filePaths []string, config *Config) (Files, error) {
	// var files Files

	filePaths = removeBlackList(filePaths, config.BlackList)

	bar := progressbar.NewOptions(len(filePaths),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionSetElapsedTime(true),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts())

	destFiles := make(Files, len(filePaths))

	pool, err := ants.NewPool(16)
	if err != nil {
		return nil, err
	}

	defer pool.Release()

	var wg sync.WaitGroup

	for i, path := range filePaths {
		err := pool.Submit(newFileWrapper(&wg, i, path, bar, config, &destFiles))
		if err != nil {
			return nil, err
		}
	}
	wg.Wait()

	return destFiles, nil
}

func newFileWrapper(wg *sync.WaitGroup, idx int, path string, bar *progressbar.ProgressBar, config *Config, destFiles *Files) taskFunc {
	return func() {
		wg.Add(1)
		defer wg.Done()
		defer bar.Add(1)

		file, err := NewFile(path, NewLinks(path, config), NewFileConfig(path, config))
		if err != nil {
			log.Println(err)
			return
		}
		files := *destFiles
		files[idx] = file
	}
}

func (f Files) Run(verbose bool) error {
	bar := progressbar.NewOptions(len(f),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionSetElapsedTime(true),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts())

	pool, err := ants.NewPool(16)
	if err != nil {
		return err
	}

	defer pool.Release()

	var wg sync.WaitGroup

	for _, file := range f {

		err := pool.Submit(fileRunWrapper(file, verbose, &wg, bar))
		if err != nil {
			return err
		}
	}

	wg.Wait()
	return nil
}

func fileRunWrapper(file *File, verbose bool, wg *sync.WaitGroup, bar *progressbar.ProgressBar) taskFunc {
	return func() {
		wg.Add(1)
		defer wg.Done()
		file.Run()
		if verbose {
			file.WriteStats()
		}
		bar.Add(1)
	}
}

func (f Files) Summary() bool {
	return summaryOfFiles(f)
}
