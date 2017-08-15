package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"secret/crypt"
)

func main() {
	name := "./test.txt"
	fi, err := os.Stat(name)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		fmt.Println("directory")
	case mode.IsRegular():
		b, _ := ioutil.ReadFile(name)

		newText := crypt.Encrypt(b, "a very very very very secret key")

		ioutil.WriteFile(name, []byte(newText), 0644)
	}
}
