package main

import (
	"../messagebird"
	"fmt"
	"os"
)

var AccessKey = ""
var MessageId = ""

func main() {
	if len(AccessKey) == 0 || len(MessageId) == 0 {
		fmt.Println("You need to set an AccessKey and MessageId in this file")
		os.Exit(1)
	}

	// Create a MessageBird client with the specified AccessKey.
	mb := messagebird.New(AccessKey)

	// Fetch the Message object.
	message, err := mb.Message(MessageId)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Check for errors returned as JSON.
	if len(message.Errors) != 0 {
		for _, error := range message.Errors {
			fmt.Println("  code        :", error.Code)
			fmt.Println("  description :", error.Description)
			fmt.Println("  parameter   :", error.Parameter, "\n")
		}
		os.Exit(1)
	}

	// Print the object information.
	fmt.Println("\nThe following information was returned as a Message object:\n")
	fmt.Println("  id                :", message.Id)
	fmt.Println("  href              :", message.HRef)
	fmt.Println("  direction         :", message.Direction)
	fmt.Println("  type              :", message.Type)
	fmt.Println("  originator        :", message.Originator)
	fmt.Println("  body              :", message.Body)
	fmt.Println("  reference         :", message.Reference)
	fmt.Println("  validity          :", message.Validity)
	fmt.Println("  gateway           :", message.Gateway)

	if len(message.TypeDetails) > 0 {
		fmt.Println("  typeDetails")
		for k, v := range message.TypeDetails {
			fmt.Println("    ", k, " : ", v)
		}
	}

	fmt.Println("  datacoding        :", message.DataCoding)
	fmt.Println("  mclass            :", message.MClass)
	fmt.Println("  scheduledDatetime :", message.ScheduledDatetime)
	fmt.Println("  createdDatetime   :", message.CreatedDatetime)
	fmt.Println("  recipients")
	fmt.Println("    totalCount               :", message.Recipients.TotalCount)
	fmt.Println("    totalSentCount           :", message.Recipients.TotalSentCount)
	fmt.Println("    totalDeliveredCount      :", message.Recipients.TotalDeliveredCount)
	fmt.Println("    TotalDeliveryFailedCount :", message.Recipients.TotalDeliveryFailedCount)
	fmt.Println("    items")

	for _, recipient := range message.Recipients.Items {
		fmt.Println("      recipient      :", recipient.Recipient)
		fmt.Println("      status         :", recipient.Status)
		fmt.Println("      statusDatetime :", recipient.StatusDatetime, "\n")
	}
}
