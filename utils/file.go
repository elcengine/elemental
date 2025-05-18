package utils

import (
	"log"
	"os"
)

func CreateAndWriteToFile(path string, content string) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatal("Error creating file", err)
	}
	defer f.Close()
	f.WriteString(content)
}
