package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var cmdName, cmdArg, mysqlUser, mysqlPass, mysqlHost, mysqlPort, mysqlDatabase string

	flag.StringVar(&cmdName, "cmdName", "cat test.txt", "Provide a command name")
	flag.StringVar(&cmdArg, "cmdArg", "", "Provide a command argument")
	flag.StringVar(&mysqlUser, "mysqlUser", "root", "Mysql user")
	flag.StringVar(&mysqlPass, "mysqlPass", "ANSKk08aPEDbFjDO", "Mysql password")
	flag.StringVar(&mysqlHost, "mysqlHost", "127.0.0.1", "Mysql host")
	flag.StringVar(&mysqlPort, "mysqlPort", "3306", "Mysql port")
	flag.StringVar(&mysqlDatabase, "mysqlDatabase", "local", "Mysql database")

	flag.Parse()

	fmt.Println("=================")
	fmt.Println("Command name:", cmdName)
	fmt.Println("Command argument:", cmdArg)
	fmt.Println("Mysql user:", mysqlUser)
	fmt.Println("Mysql password:", mysqlPass)
	fmt.Println("Mysql host:", mysqlHost)
	fmt.Println("Mysql port:", mysqlPort)
	fmt.Println("Mysql database:", mysqlDatabase)
	fmt.Println("=================")

	// Open database connection
	src := mysqlUser + ":" + mysqlPass + "@tcp(" + mysqlHost + ":" + mysqlPort + ")/" + mysqlDatabase
	db, err := sql.Open("mysql", src)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	} else {
		fmt.Println("Connected to database successfully")
	}
	defer db.Close()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	fmt.Println("=================>> Start running command <<=================")
	for range ticker.C {
		// Run command
		cmd := exec.Command(cmdName, cmdArg)
		output, err := cmd.Output()
		if err != nil {
			fmt.Println("Error running command:", err)
			continue
		}

		// Parse output
		outputStr := string(output)
		consumedAmount := parseOutput(outputStr)

		// Update database
		updateDatabase(db, consumedAmount)
	}
}

type clientData struct {
	ClientID       string
	ConsumedAmount float64
}

func parseOutput(output string) clientData {
	parts := strings.Split(output, ",")
	clientData := clientData{}
	for _, part := range parts {
		if strings.HasPrefix(part, "lient_id:") {
			clientData.ClientID = strings.TrimPrefix(part, "lient_id:")
			continue
		}
		if strings.HasPrefix(part, "consumed_amount:") {
			consumedAmountStr := strings.TrimPrefix(part, "consumed_amount:")
			consumedAmountStr = strings.Trim(consumedAmountStr, ";")

			consumedAmount, err := strconv.ParseFloat(consumedAmountStr, 64)
			if err != nil {
				fmt.Println("Error parsing consumed amount:", err)
				return clientData
			}
			clientData.ConsumedAmount = consumedAmount
			continue
		}
	}
	return clientData
}

func updateDatabase(db *sql.DB, cd clientData) {
	// Execute update query
	_, err := db.Exec("UPDATE credit_card SET credit_money = (credit_money - ?) WHERE customer_id = ?", cd.ConsumedAmount, cd.ClientID)
	if err != nil {
		fmt.Println("Error updating database:", err)
		return
	}
	fmt.Println("Database updated successfully")
}
