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

// main is the entry point of the program.
//
// No parameters.
// No return values.
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
			break
		}
		if len(output) == 0 {
			fmt.Println("Command output is empty")
			continue
		}

		// Parse output
		outputStr := string(output)
		clients := parseOutput(outputStr)

		if clients == nil {
			fmt.Println("Do not found the active client")
			continue
		}

		for _, client := range clients {
			if client.ClientID == "" || client.ConsumedAmount == 0 {
				continue
			}
			// Update database
			updateDatabase(db, client)
		}
	}
}

type client struct {
	ClientID       string
	ConsumedAmount float64
}

// parseOutput parses the output string and returns a slice of clients.
//
// Parameter:
// output string - the output string to parse
// Return:
// []client - a slice of client structs parsed from the output string
func parseOutput(output string) (clients []client) {
	outputArr := strings.Split(output, ";")

	if len(outputArr) == 0 {
		fmt.Println("No clients found")
		return nil
	}

	for _, data := range outputArr {
		parts := strings.Split(data, ",")
		client := client{}
		for _, part := range parts {
			if strings.HasPrefix(part, "client_id:") {
				client.ClientID = strings.TrimPrefix(part, "client_id:")
				continue
			}
			if strings.HasPrefix(part, "consumed_amount:") {
				consumedAmountStr := strings.TrimPrefix(part, "consumed_amount:")
				consumedAmountStr = strings.Trim(consumedAmountStr, ";")

				consumedAmount, err := strconv.ParseFloat(consumedAmountStr, 64)
				if err != nil {
					fmt.Println("Error parsing consumed amount:", err)
					return nil
				}
				client.ConsumedAmount = consumedAmount
				continue
			}
		}
		clients = append(clients, client)
	}

	return clients
}

// updateDatabase updates the credit_card table in the database by deducting the consumed amount
// from the credit_money where the customer_id matches the client's ID.
//
// Parameters:
// - db: *sql.DB - pointer to the database connection
// - cd: client - the client struct containing ConsumedAmount and ClientID
func updateDatabase(db *sql.DB, cd client) {
	// Execute update query
	_, err := db.Exec("UPDATE credit_card SET credit_money = (credit_money - ?) WHERE customer_id = ?", cd.ConsumedAmount, cd.ClientID)
	if err != nil {
		fmt.Println("Error updating database:", err)
		return
	}
	fmt.Println("Database updated successfully for client: " + cd.ClientID)
}
