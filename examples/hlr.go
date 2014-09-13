package main

import (
	"../messagebird"
	"fmt"
	"os"
)

var AccessKey = ""
var HLRId = ""

func main() {
	if len(AccessKey) == 0 || len(HLRId) == 0 {
		fmt.Println("You need to set an AccessKey and HLRId in this file")
		os.Exit(1)
	}

	// Create a MessageBird client with the specified AccessKey.
	mb := messagebird.New(AccessKey)

	// Fetch the HLR object.
	hlr, err := mb.HLR("d26c94c0353bd8e171a3979h97860638")
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
