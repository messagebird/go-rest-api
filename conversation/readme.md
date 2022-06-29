# Conversation API

Package conversation provides an interface to the Conversations API. The API's docs are located here:
https://developers.messagebird.com/docs/conversations.

The `api.go` file contains all structs which uses in several others and some wrapping method to provide API requests and
pagination.

# Examples

## Conversations

### Start Conversation
```go
// create a client
client := messagebird.New("your-access-key")

conv, err := conversation.Start(client, &conversation.StartRequest{
	ChannelID: "619747f69cf940a98fb443140ce9aed2",
    To:        "31612345678",
    Content: &conversation.MessageContent{
        Text: "Hello",
    },
    Type: conversation.MessageTypeText,
})
if err != nil {
    // handle error
    return
}

// display the results
fmt.Println("conversationId: ", conv.ID)
```

### Reply Conversation
```go
// create a client
client := messagebird.New("your-access-key")

message, err := conversation.Reply(client, "2e15efafec384e1c82e9842075e87beb", &conversation.ReplyRequest{
    ChannelID: "chid",
    Content: &conversation.MessageContent{
        Text: "Hello world",
    },
    Type: conversation.MessageTypeText,
})
if err != nil {
    // handle error
    return
}

// display the results
fmt.Println("messageId: ", message.ID)
```

### Get Conversation
```go
// create a client
client := messagebird.New("your-access-key")

conv, err := conversation.Read(client, "2e15efafec384e1c82e9842075e87beb")
if err != nil {
    // handle error
    return
}

// display the results
fmt.Println("conversationId: ", conv.ID)
```

### Update Conversation
```go
// create a client
client := messagebird.New("your-access-key")

conv, err := conversation.Update(client, "id", &conversation.UpdateRequest{
    Status: conversation.ConversationStatusArchived,
})
if err != nil {
    // handle error
    return
}

// display the results
fmt.Println("status: ", conv.Status)
```

### List Conversations
```go
// create a client
client := messagebird.New("your-access-key")

convList, err := conversation.List(client, &conversation.ListRequest{conversation.PaginationRequest{Limit: 10, Offset: 20}, "", nil})
if err != nil {
    // handle error
    return
}

// display the results
fmt.Println("total: ", convList.Total)
fmt.Println(convList.Items)
```

### Get Conversations by Contact
```go
// create a client
client := messagebird.New("your-access-key")

convList, err := conversation.ListByContact(client, "ebf6aceed7ae4375b726e247318d3377", &conversation.PaginationRequest{Limit: 20, Offset: 2})
if err != nil {
    // handle error
    return
}

// display the results
fmt.Println(convList.Items)
```

## Messages

### Send Message
```go
// create a client
client := messagebird.New("your-access-key")

message, err := conversation.SendMessage(client, &conversation.SendMessageRequest{
    To:   "+31624971134",
    From: "MessageBird",
    Type: conversation.MessageTypeText,
    Content: &conversation.MessageContent{
        Text: "Hello world",
    },
    ReportUrl: "https://myreport.site",
    Source:    map[string]interface{}{"name": "Valera"},
})
if err != nil {
    // handle error
    return
}

// display the results
fmt.Println("messageId: ", message.ID)
```

### Get Message
```go
// create a client
client := messagebird.New("your-access-key")

message, err := conversation.ReadMessage(client, "5f3437fdb8444583aea093a047ac014b")
if err != nil {
    // handle error
    return
}

// display the results
fmt.Println("messageId: ", message.ID)
```

### List Messages
```go
// create a client
client := messagebird.New("your-access-key")

messageList, err := conversation.ListMessages(client, &conversation.ListMessagesRequest{Ids: "5f3437fdb8444583aea093a047ac014b,4abc37fdb8444583aea093a047ac014c"})
if err != nil {
    // handle error
    return
}

// display the results
fmt.Println(messageList.Items)
```

### Get Messages in Conversation
```go
// create a client
client := messagebird.New("your-access-key")

messageList, err := conversation.ListConversationMessages(
    client,
    conversationId,
    &conversation.ListConversationMessagesRequest{conversation.PaginationRequest{20, 2}, "sms,whatsapp,facebook"},
)
if err != nil {
    // handle error
    return
}

// display the results
fmt.Println(messageList.Items)
```

## Webhooks

### Create a Webhook
```go
// create a client
client := messagebird.New("your-access-key")

webhook, err := conversation.CreateWebhook(client, &conversation.WebhookCreateRequest{
    ChannelID: "chid",
    Events: []conversation.WebhookEvent{
        conversation.WebhookEventConversationCreated,
        conversation.WebhookEventMessageUpdated,
    },
    URL: "https://example.com/webhooks",
})
if err != nil {
    // handle error
    return
}

// display the results
fmt.Println("webhookId: ", webhook.ID)
```

### List Webhooks
```go
// create a client
client := messagebird.New("your-access-key")

webhookList, err := conversation.ListWebhooks(client, &conversation.PaginationRequest{Limit: 20, Offset: 2})
if err != nil {
    // handle error
    return
}

// display the results
fmt.Println("total: ", webhookList.TotalCount)
fmt.Println(webhookList.Items)
```

### Get Webhook
```go
// create a client
client := messagebird.New("your-access-key")

webhook, err := conversation.ReadWebhook(client, "985ae50937a94c64b392531ea87a0263")
if err != nil {
    // handle error
    return
}

// display the results
fmt.Println("webhookId: ", webhook.ID)
```

### Update Webhook
```go
// create a client
client := messagebird.New("your-access-key")

webhook, err := conversation.UpdateWebhook(client, "985ae50937a94c64b392531ea87a0263", &conversation.WebhookUpdateRequest{
    Events: []conversation.WebhookEvent{
        conversation.WebhookEventConversationUpdated,
    },
    URL:    "https://example.com/mynewwebhookurl",
    Status: conversation.WebhookStatusDisabled,
})
if err != nil {
    // handle error
    return
}

// display the results
fmt.Println("webhookStatus: ", webhook.Status)
```

### Delete Webhook
```go
// create a client
client := messagebird.New("your-access-key")

err := conversation.DeleteWebhook(client, "985ae50937a94c64b392531ea87a0263")
if err != nil {
    // handle error
    return
}
```
