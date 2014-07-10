package main

import (
	"../messagebird"
	"fmt"
	"os"
)

var AccessKey = "test_gshuPaZoeEG6ovbc8M79w0QyM"

func main() {
	// Create a MessageBird client with the specified AccessKey.
	mb := &messagebird.Client{AccessKey: AccessKey}

	// Fetch the HLR object.
	hlr, err := mb.CreateHLR("31612345678", "MyReference")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Check for errors returned as JSON.
	if len(hlr.Errors) != 0 {
		for _, error := range hlr.Errors {
			fmt.Println("  code        :", error.Code)
			fmt.Println("  description :", error.Description)
			fmt.Println("  parameter   :", error.Parameter, "\n")
		}
		os.Exit(1)
	}

	// Print the object information.
	fmt.Println("\nThe following information was returned as an HLR object:\n")
	fmt.Println("  id              :", hlr.Id)
	fmt.Println("  href            :", hlr.HRef)
	fmt.Println("  msisdn          :", hlr.MSISDN)
	fmt.Println("  reference       :", hlr.Reference)
	fmt.Println("  status          :", hlr.Status)
	fmt.Println("  createdDatetime :", hlr.CreatedDatetime)
	fmt.Println("  statusDatetime  :", hlr.StatusDatetime, "\n")
}
