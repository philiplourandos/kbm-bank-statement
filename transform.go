package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

type Entry struct {
	date          time.Time
	statementType string
	amount        string
	statementLine string
	apartment     string
}

type Owner struct {
	apartment  string
	paymentIds []string
}

func main() {
	var src = "data.csv"

	fmt.Printf("File to process: [%s]\n", src)

	ownerFile, err := os.Open("apartment.csv")
	if err != nil {
		fmt.Printf("Unable to open apartments.csv")

		os.Exit(1)
	}

	var owners []Owner
	apartment, err := csv.NewReader(ownerFile).ReadAll()
	if err != nil {
		fmt.Printf("err i: %s", err)
	}

	for _, line := range apartment {
		entries := len(line)

		var currentOwner Owner

		for i := 1; i < entries; i++ {
			if line[i] != "" {
				currentOwner.paymentIds = append(currentOwner.paymentIds, strings.TrimSpace(line[i]))
			}
		}
		currentOwner.apartment = line[0]

		owners = append(owners, currentOwner)
	}

	ownerFile.Close()

	f, err := os.Open(src)

	if err != nil {
		fmt.Printf("Unable to open file")

		os.Exit(1)
	}

	var statementLines []Entry

	lines, _ := csv.NewReader(f).ReadAll()

	rx, _ := regexp.Compile("^*?([1-9][0-9]*?\\.[0-9]{2})")

	for _, line := range lines {
		var entry Entry

		trimmedDate := []rune(line[1])[1:]

		date, dateErr := time.Parse("20060102", string(trimmedDate))
		if dateErr != nil {
			fmt.Println(dateErr)
			os.Exit(1)
		}

		entry.date = date
		entry.statementType = line[2]

		amount := string(line[3][1:])
		vals := rx.FindAllStringSubmatch(amount, 1)
		entry.amount = vals[0][1]

		entry.statementLine = line[4]

		statementLines = append(statementLines, entry)
	}

	f.Close()

	for statementIndex := 0; statementIndex < len(statementLines); statementIndex++ {
		currentStatement := &statementLines[statementIndex]

		for _, currentOwner := range owners {
			for _, currentId := range currentOwner.paymentIds {
				if strings.Contains(currentStatement.statementLine, currentId) {
					currentStatement.apartment = currentOwner.apartment
				}
			}
		}
	}

	fmt.Println("Writing file")

	spreadsheet, _ := os.Create("kbm-payments.csv")
	spreadsheet.WriteString("Date,Apartment,Amount\n")

	for _, currentStatement := range statementLines {
		if currentStatement.apartment != "" {
			line := fmt.Sprintf("%s,%s,%s\n", currentStatement.date.Format("2006-01-02"),
				currentStatement.apartment, currentStatement.amount)
			spreadsheet.WriteString(line)

		}
	}

	spreadsheet.Close()

	for _, currentOwner := range owners {

	}
}
