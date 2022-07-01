MessageBird's REST API for Go
=============================
This repository contains the open source Go client for MessageBird's REST API. Documentation can be found at: https://developers.messagebird.com.

[![Build Status](https://travis-ci.org/messagebird/go-rest-api.svg?branch=master)](https://travis-ci.org/messagebird/go-rest-api) [![PkgGoDev](https://pkg.go.dev/badge/github.com/messagebird/go-rest-api/v9)](https://pkg.go.dev/github.com/messagebird/go-rest-api/v9)

Requirements
------------
- [Sign up](https://www.messagebird.com/en/signup) for a free MessageBird account
- Create a new access key in the [dashboard](https://dashboard.messagebird.com/en-us/developers/access).
- An application written in Go to make use of this API

Installation
------------
The easiest way to use the MessageBird API in your Go project is to install it using *go get*:

```
$ go get github.com/messagebird/go-rest-api/v9
```

Examples
--------
Here is a quick example on how to get started. Assuming the **go get** installation worked, you can import the messagebird package like this:

```go
import "github.com/messagebird/go-rest-api/v9"
```

Then, create an instance of **messagebird.Client**. It can be used to access the MessageBird APIs.

```go
// Access keys can be managed through our dashboard.
accessKey := "your-access-key"

// Create a client.
client := messagebird.New(accessKey)

// Request the balance information, returned as a balance.Balance object.
balance, err := balance.Read(client)
if err != nil {
	// Handle error.
	return
}

// Display the results.
fmt.Println("Payment: ", balance.Payment)
fmt.Println("Type:", balance.Type)
fmt.Println("Amount:", balance.Amount)
```

This will give you something like:

```bash
$ go run example.go
Payment: prepaid
Type: credits
Amount: 9
```

Please see the other examples for a complete overview of all the available API calls.

Errors
------
When something goes wrong, our APIs can return more than a single error. They are therefore returned by the client as "error responses" that contain a slice of errors.

It is important to notice that the Voice API returns errors with a format that slightly differs from other APIs.
For this reason, errors returned by the `voice` package are of type `voice.ErrorResponse`. It contains `voice.Error` structs. All other packages return `messagebird.ErrorResponse` structs that contain a slice of `messagebird.Error`.

An example of "simple" error handling is shown in the example above. Let's look how we can gain more in-depth insight in what exactly went wrong:

```go
import "github.com/messagebird/go-rest-api/v9"
import "github.com/messagebird/go-rest-api/v9/sms"

// ...

_, err := sms.Read(client, "some-id")
if err != nil {
	mbErr, ok := err.(messagebird.ErrorResponse)
	if !ok {
		// A non-MessageBird error occurred (no connection, perhaps?) 
		return err
	}
	
	fmt.Println("Code:", mbErr.Errors[0].Code)
	fmt.Println("Description:", mbErr.Errors[0].Description)
	fmt.Println("Parameter:", mbErr.Errors[0].Parameter)
}
```

`voice.ErrorResponse` is very similar, except that it holds `voice.Error` structs - those contain only `Code` and `Message` (not description!) fields:

```go
import "github.com/messagebird/go-rest-api/v9/voice"

// ...

_, err := voice.CallFlowByID(client, "some-id")
if err != nil {
	vErr, ok := err.(voice.ErrorResponse)
	if !ok {
    		// A non-MessageBird (Voice) error occurred (no connection, perhaps?) 
    		return err
    }
	
	fmt.Println("Code:", vErr.Errors[0].Code)
	fmt.Println("Message:", vErr.Errors[0].Message)
}
```

Documentation
-------------
Complete documentation, instructions, and examples are available at:
[https://developers.messagebird.com](https://developers.messagebird.com).

Upgrading
---------
If you're upgrading from older versions, please read the [Messagebird `go-rest-api` upgrading guide](UPGRADING.md).

License
-------
The MessageBird REST Client for Go is licensed under [The BSD 2-Clause License](http://opensource.org/licenses/BSD-2-Clause). Copyright (c) 2022 MessageBird
