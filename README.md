# DataEdge Go SDK

This is a Go client library for the [DataEdge](https://dataedge.md) API, providing a full-featured and type-safe interface for SMS and Viber messaging services. It is a port of the official [PHP SDK](https://gitlab.com/dataedgemd/php-sdk).

## Installation

```bash
go get github.com/acolev/dataedge-go
```

## Usage

### Initialization

```go
package main

import (
    "context"
    "fmt"
    "github.com/acolev/dataedge-go"
)

func main() {
    // Initialize client with defaults
    client := dataedge.NewClient(
        "YOUR_API_TOKEN",
        dataedge.WithDebug(true),
        dataedge.WithDefaultAlphaNameID(123), // Set default sender ID
    )
    
    // Create service instances
    smsService := dataedge.NewSmsService(client)
    viberService := dataedge.NewViberService(client)
    
    // Now you can call SendQuickSms with 0 as alphaNameID to use the default
    result, err := smsService.SendQuickSms(context.Background(), "Hello World", "373xxxxxxxxx", 0, true)
    balance, err := smsService.GetBalance(context.Background())
    if err != nil {
        panic(err)
    }
    fmt.Printf("Balance: %.2f\n", balance)
}
```

### Sending SMS

```go
// Send a quick SMS (without creating it first)
result, err := smsService.SendQuickSms(context.Background(), "Hello World", "373xxxxxxxxx", 0, true)
if err != nil {
    panic(err)
}
fmt.Println("SMS Sent:", result)

// Create and Send
msgID, err := smsService.CreateSMSMessage(context.Background(), "Hello World", 0, true)
if err != nil {
    panic(err)
}
res, err := smsService.SendSms(context.Background(), msgID, "373xxxxxxxxx")
```

### Sending Viber Messages

```go
// Send quick Viber message
res, err := viberService.SendQuickViberMessage(context.Background(), "373xxxxxxxxx", 123, "Hello Viber")

// Send Viber message with image
res, err = viberService.SendQuickViberMessageWithImage(
    context.Background(),
    "373xxxxxxxxx",
    123,
    "Hello with Image",
    "IMAGE",
    "/path/to/image.jpg",
    "Button Text",
    "https://example.com",
)
```

## Features

- **Full API Coverage**: Supports all methods from the PHP SDK.
- **Context Support**: All methods accept `context.Context` for timeout and cancellation control.
- **Type Safety**: Uses Go structs and error handling.
- **Transliteration**: Built-in support for Cyrillic to Latin transliteration.
- **Multipart Upload**: Support for creating Viber messages with images.

## License

MIT
