package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

type Entry struct {
	date string
	statementType string
	amount string
	statementLine string
	apartment string
}

type Owner struct {
	apartment string
	paymentIds []string
}

func main() {
	var src = "data.csv"

	fmt.Printf("File to process: [%s]\n", src)

	onwerFile, err := os.Open("apartment.csv")
	if err != nil {
		fmt.Printf("Unable to open apartments.csv")

		os.Exit(1)
	}

	var owners []Owner
	apartment, err := csv.NewReader(onwerFile).ReadAll()
	if err != nil {
		fmt.Printf("err i: %s", err)
	}

	for _, line := range apartment {
		entries := len(line)

		var currentOwner Owner

		for i := 1; i < entries ; i++ {
			currentOwner.paymentIds = append(currentOwner.paymentIds, strings.TrimSpace(line[i]))
		}
		currentOwner.apartment = line[0]

		owners = append(owners, currentOwner)
		fmt.Printf("%s\n", currentOwner)
	}

	f, err := os.Open(src)

	if err != nil {
		fmt.Printf("Unable to open file")

		os.Exit(1)
	}

	var statementLines []Entry

	lines, _ := csv.NewReader(f).ReadAll()

	for _, line := range lines {
		var entry Entry
		entry.date = line[1]
		entry.statementType = line[2]
		entry.amount = line[3][1:len(line[3])]
		entry.statementLine = line[4]

		statementLines = append(statementLines, entry)
	}

	for _, currentStatement := range statementLines {
		for _, currentOwner := range owners {
			for _, currentId := range currentOwner.paymentIds {

				if strings.Contains(currentStatement.statementLine, currentId) {
					currentStatement.apartment = currentOwner.apartment
				}
			}
		}
	}
}
