# DataEdge Go SDK

Go client for [DataEdge API](https://dataedge.md/).

[![Go Reference](https://pkg.go.dev/badge/github.com/acolev/dataedge-go.svg)](https://pkg.go.dev/github.com/acolev/dataedge-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/acolev/dataedge-go)](https://goreportcard.com/report/github.com/acolev/dataedge-go)

[Русская версия (README_RU.md)](README_RU.md)

## Installation

```bash
go get github.com/acolev/dataedge-go
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/acolev/dataedge-go"
)

func main() {
	// Initialize client
	client := dataedge.NewClient("your-api-token",
		dataedge.WithAlphaNameID(123),  // Default SMS sender ID
		dataedge.WithViberNameID(456),  // Default Viber sender ID
	)

	smsService := dataedge.NewSmsService(client)

	// Send a quick SMS
	res, err := smsService.SendQuickSms(context.Background(), "Hello from Go!", "37360000000")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Result: %v\n", res)
}
```

## Client Configuration

Create a client using `NewClient(token string, opts ...Option)`.

### Options
- `WithBaseURL(url string)`: Set a custom API base URL.
- `WithDebug(debug bool)`: Enable debug logging for requests and responses.
- `WithAlphaNameID(id int)`: Set default sender ID for SMS.
- `WithViberNameID(id int)`: Set default sender ID for Viber.
- `WithTranslit(enabled bool)`: Enable/disable auto-transliteration (default is `true`).
- `WithHTTPClient(client *http.Client)`: Use a custom HTTP client.

---

## SMS Service

Access methods using `dataedge.NewSmsService(client)`.

### Methods

| Method | Description |
|--------|-------------|
| `GetBalance(ctx)` | Returns current balance as `float64`. |
| `GetLimit(ctx)` | Returns current message limit. |
| `GetAlphaNames(ctx)` | Returns list of available sender names. |
| `GetAlphaNameID(ctx, name)` | Returns ID for a specific sender name. |
| `CreateSMSMessage(ctx, text)` | Creates a message template, returns `messageID`. |
| `SendSms(ctx, messageID, phone)` | Sends a previously created message. |
| `SendQuickSms(ctx, text, phone)` | Creates and sends SMS in one call. |
| `CheckSMS(ctx, smsID)` | Checks status of a sent SMS. |
| `GetMessagesList(ctx)` | Returns history of sent messages. |

---

## Viber Service

Access methods using `dataedge.NewViberService(client)`.

### Methods

| Method | Description |
|--------|-------------|
| `CreateViberMessage(ctx, text)` | Creates a Viber template. |
| `SendViberMessage(ctx, phone, messageID)` | Sends a created Viber message. |
| `SendQuickViberMessage(ctx, phone, text)` | Sends a quick Viber message. |
| `SendQuickViberMessageWithImage(...)` | Sends Viber message with image, button, and link. |
| `SendViberMessageList(ctx, name, listID, schedule, extra)` | Sends message to a predefined list. |
| `GetViberMessageList(ctx)` | Returns history of Viber messages. |

#### Sending Viber with Image Example
```go
viberService := dataedge.NewViberService(client)

res, err := viberService.SendQuickViberMessageWithImage(
    ctx,
    "37360000000",
    "Check this out!",
    "IMAGE",           // Type: MESSAGE or IMAGE
    "/path/to/img.jpg",
    "Open Link",       // Button text
    "https://google.com",
)
```

---

## Utilities

### Transliteration
Auto-transliteration is enabled by default for SMS. You can also use it manually:
```go
text := dataedge.Transliterate("Привет мир") // "Privet mir"
```

## Features

- **Full API Coverage**: Supports all methods from the PHP SDK.
- **Context Support**: All methods accept `context.Context` for timeout and cancellation control.
- **Type Safety**: Uses Go structs and error handling.
- **Transliteration**: Built-in support for Cyrillic to Latin transliteration.
- **Multipart Upload**: Support for creating Viber messages with images.

## License

Custom Non-Commercial License (see [LICENSE](LICENSE)). Free for personal and internal use, resale is prohibited.
