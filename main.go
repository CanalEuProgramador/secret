package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"secret/crypt"
	"sync"
)

var wg sync.WaitGroup

var decrypt bool
var secret string

func init() {
	flag.BoolVar(&decrypt, "decrypt", false, "-decrypt")
	flag.StringVar(&secret, "secret", "", "-secret=your_secret")
	flag.Parse()
}

func exec(name string) {
	defer wg.Done()
	b, _ := ioutil.ReadFile(name)

	var newText string
	if decrypt {
		newText = crypt.Decrypt(b, secret)
	} else {
		newText = crypt.Encrypt(b, secret)
	}

	ioutil.WriteFile(name, []byte(newText), 0644)
}

func interate(name string) {
	defer wg.Done()
	fi, err := os.Stat(name)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		var i = 0
		filepath.Walk(name, func(path string, f os.FileInfo, err error) error {
			if i > 0 {
				wg.Add(1)
				go interate(path)
			}
			i++
			return nil
		})
	case mode.IsRegular():
		fmt.Println("Running at:", name)
		wg.Add(1)
		go exec(name)
	}
}

func main() {
	name := os.Args[len(os.Args)-1]
	if secret == "" {
		panic("Invalid secret, try: -secret=your_secret")
	}

	wg.Add(1)
	go interate(name)
	wg.Wait()
}
