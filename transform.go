package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type Entry struct {
	date          time.Time
	statementType string
	amount        string
	statementLine string
	apartment     string
}

type Owners struct {
	Owner []struct {
		Apartment  string `yaml:"apartment"`
		PaymentIds []string `yaml:"paymentIds"`
	} `yaml:"owner"`
}

func main() {
	// Load apartment owner metadata
	content, _ := ioutil.ReadFile("apartments.yml")

	owners := Owners{}
	err := yaml.Unmarshal(content, &owners)
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	// Load bank statement data
	var src = "data.csv"
	fmt.Printf("File to process: [%s]\n", src)

	f, err := os.Open(src)

	if err != nil {
		fmt.Printf("Unable to open file")

		os.Exit(1)
	}

	var statementLines []Entry

	lines, csvErr := csv.NewReader(f).ReadAll()

	if csvErr != nil {
		fmt.Print(csvErr)

		os.Exit(1)
	}

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

		for _, currentOwner := range owners.Owner {
			for _, currentId := range currentOwner.PaymentIds {
				if strings.Contains(currentStatement.statementLine, currentId) {
					currentStatement.apartment = currentOwner.Apartment
				}
			}
		}
	}

	fmt.Println("Writing owner payments to file")

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

	// Find outstanding owners who have not paid in the last 40 days
	fmt.Println("finding late payment")
	var latePaymentThreshold = time.Now().AddDate(0, 0, -40)
	var lastPayments []Entry

	for _, currentOwner := range owners.Owner {
		var found Entry

		for _, currentStatement := range statementLines {
			if currentStatement.apartment != "" && currentOwner.Apartment == currentStatement.apartment {
				found = currentStatement
			}
		}

		if found.date.Before(latePaymentThreshold) {
			lastPayments = append(lastPayments, found)
			fmt.Printf("Apartment: %s, last paid on: %s\n", found.apartment, found.date)
		}
	}
}
