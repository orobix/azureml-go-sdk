package workspace

import (
	"fmt"
	"log"
	"os"
)

const (
	exampleRespsDir = "assets"
)

func loadExampleResp(fileName string) []byte {
	file, err := os.ReadFile(fmt.Sprintf("%s/%s", exampleRespsDir, fileName))
	if err != nil {
		log.Fatalf("Could not load example resp \"%s\"", fileName)
	}
	return file
}
