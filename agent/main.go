package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/tobscher/kiss/configuration"
)

func main() {
	bio := bufio.NewReader(os.Stdin)
	bytes, _, err := bio.ReadLine()
	if err != nil {
		log.Fatal(err)
	}

	var task configuration.Task
	err = json.Unmarshal(bytes, &task)

	log.Printf("%+v\n", task)

	err = RunTask(&task)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
