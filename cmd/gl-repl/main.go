package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ymsodev/gl"
)

var runtime *gl.GL

func init() {
	runtime = gl.New()
	runtime.Init()
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	for scan(s) {
		text := s.Text()
		run(text)
	}
	if err := s.Err(); err != nil {
		fmt.Println(err)
	}
}

func scan(s *bufio.Scanner) bool {
	fmt.Print(">>> ")
	return s.Scan()
}

func run(text string) {
	// TODO: REPL commands

	val := runtime.Run(text)
	fmt.Println(val.String())
}
