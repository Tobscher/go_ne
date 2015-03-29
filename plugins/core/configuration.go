package plugin

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
)

func LoadConfig(reader io.Reader, v interface{}) {
	bio := bufio.NewReader(reader)
	bytes, _, err := bio.ReadLine()
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(bytes, v)
	if err != nil {
		log.Fatal(err)
	}
}
