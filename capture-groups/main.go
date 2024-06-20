package main

import (
	"fmt"
	"regexp"
)

func main() {
	demo()

}

func demo() {
	re := regexp.MustCompile(`(\w+):(\d+)`)

	matches := re.FindSubmatch([]byte("bob:25"))

	name := matches[1]
	age := matches[2]
	fmt.Printf("%v is %v years old\n", string(name), string(age))
}

// func demo() {
// 	re := regexp.MustCompile(`(?P<name>\w+):(?P<age>\d+)`)

// 	matches := re.FindSubmatch([]byte("bob:25"))

// 	name := matches[re.SubexpIndex("name")]
// 	age := matches[re.SubexpIndex("age")]
// 	fmt.Printf("%v is %v years old\n", string(name), string(age))
// }
