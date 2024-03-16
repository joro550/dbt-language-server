package tests

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"testing"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

func testGetDefinition(t *testing.T) {
	fileContent, err := os.Open("second_dbt_model_2.sql")
	if err != nil {
		t.Fatal(err)
	}

	compiledRegex, err := regexp.Compile("'[\\S]*'")
	if err != nil {
		t.Fatal(err)
	}

	scanner := bufio.NewScanner(fileContent)

	// var postions []position
	line := 0

	for scanner.Scan() {
		indicies := compiledRegex.FindAllIndex(scanner.Bytes(), -1)
		if len(indicies) != 0 {
			currentLineText := scanner.Text()

			for _, val := range indicies {
				fmt.Println(currentLineText[val[0]:val[1]])
			}
		}

		line += 1
	}
}

func TestThing(t *testing.T) {
	position := protocol.Position{Line: 0, Character: 10}
	code := "{{ ref('my_first_dbt_model') }}"

	node := main.Node{}
}

type position struct {
	Line  int
	Index int
}
