package main

import (
	"../messagebird"
	"fmt"
	"net/url"
	"os"
)

var AccessKey = "test_gshuPaZoeEG6ovbc8M79w0QyM"

func main() {
	// Create a MessageBird client with the specified AccessKey.
	mb := messagebird.New(AccessKey)

	// The optional parameters.
	params := &url.Values{"reference": {"MyReference"}}

	// Fetch the VoiceMessage object.
	vmsg, err := mb.CreateVoiceMessage([]string{"31612345678"}, "Hello World", params)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Check for errors returned as JSON.
	if len(vmsg.Errors) != 0 {
		for _, error := range vmsg.Errors {
			fmt.Println("  code        :", error.Code)
			fmt.Println("  description :", error.Description)
			fmt.Println("  parameter   :", error.Parameter, "\n")
		}
		os.Exit(1)
	}

	// Print the object information.
	fmt.Println("\nThe following information was returned as a VoiceMessage object:\n")
	fmt.Println("  id                :", vmsg.Id)
	fmt.Println("  href              :", vmsg.HRef)
	fmt.Println("  body              :", vmsg.Body)
	fmt.Println("  reference         :", vmsg.Reference)
	fmt.Println("  language          :", vmsg.Language)
	fmt.Println("  voice             :", vmsg.Voice)
	fmt.Println("  repeat            :", vmsg.Repeat)
	fmt.Println("  ifMachine         :", vmsg.IfMachine)
	fmt.Println("  scheduledDatetime :", vmsg.ScheduledDatetime)
	fmt.Println("  createdDatetime   :", vmsg.CreatedDatetime)
	fmt.Println("  recipients")
	fmt.Println("    totalCount               :", vmsg.Recipients.TotalCount)
	fmt.Println("    totalSentCount           :", vmsg.Recipients.TotalSentCount)
	fmt.Println("    totalDeliveredCount      :", vmsg.Recipients.TotalDeliveredCount)
	fmt.Println("    TotalDeliveryFailedCount :", vmsg.Recipients.TotalDeliveryFailedCount)
	fmt.Println("    items")

	for _, recipient := range vmsg.Recipients.Items {
		fmt.Println("      recipient      :", recipient.Recipient)
		fmt.Println("      status         :", recipient.Status)
		fmt.Println("      statusDatetime :", recipient.StatusDatetime, "\n")
	}
}
