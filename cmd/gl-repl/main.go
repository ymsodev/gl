package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ymsodev/gl"
)

var runtime *gl.Gl

func init() {
	runtime = gl.New()
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
	val := runtime.Run(text)
	gl.Print(val)
}
