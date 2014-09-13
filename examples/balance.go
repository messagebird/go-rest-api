package main

import (
	"../messagebird"
	"fmt"
	"os"
)

var AccessKey = "test_gshuPaZoeEG6ovbc8M79w0QyM"

func main() {
	// Create a MessageBird client with the specified AccessKey.
	mb := messagebird.New(AccessKey)

	// Fetch the Balance object.
	balance, err := mb.Balance()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Check for errors returned as JSON.
	if len(balance.Errors) != 0 {
		for _, error := range balance.Errors {
			fmt.Println("  code        :", error.Code)
			fmt.Println("  description :", error.Description)
			fmt.Println("  parameter   :", error.Parameter, "\n")
		}
		os.Exit(1)
	}

	// Print the object information.
	fmt.Println("\nThe following information was returned as a Balance object:\n")
	fmt.Println("  payment :", balance.Payment)
	fmt.Println("  type    :", balance.Type)
	fmt.Println("  amount  :", balance.Amount, "\n")
}
