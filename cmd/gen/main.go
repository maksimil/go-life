package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
)

const (
	shaderdir     = "./cmd/life/shaders"
	shaderoutfile = "./cmd/life/gen_shaders.go"
)

func main() {
	shadersfolder, err := os.ReadDir(shaderdir)
	if err != nil {
		panic(err)
	}

	jobs := make(chan struct {
		name string
		data string
	}, 8)

	var wg sync.WaitGroup

	for _, v := range shadersfolder {
		if !v.IsDir() {
			wg.Add(1)
			name := v.Name()
			go func() {
				defer wg.Done()
				fpath := path.Join(shaderdir, name)

				data, err := os.ReadFile(fpath)
				if err != nil {
					panic(err)
				}

				jobs <- struct {
					name string
					data string
				}{
					strings.TrimSuffix(name, ".glsl"),
					string(data),
				}
			}()
		}
	}
	wg.Wait()
	close(jobs)

	gencontents := "package main"

	for job := range jobs {
		gencontents += fmt.Sprintf("\n\nconst %s = `\n%s`+\"\\x00\"", job.name, job.data)
	}

	err = os.WriteFile(shaderoutfile, []byte(gencontents), 0666)
	if err != nil {
		panic(err)
	}
}
