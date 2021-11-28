package workspace

import (
	"fmt"
	"io/ioutil"
	"log"
)

const (
	exampleRespsDir = "assets"
)

func loadExampleResp(fileName string) []byte {
	file, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", exampleRespsDir, fileName))
	if err != nil {
		log.Fatalf("Could not load example resp \"%s\"", fileName)
	}
	return file
}
