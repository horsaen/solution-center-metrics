package data

import (
	"bufio"
	"os"
	"strings"
)

func GetUsers() [][]string {
	var users [][]string

	file, _ := os.Open("users")

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		split := strings.Split(line, ",")

		users = append(users, split)
	}

	return users
}
