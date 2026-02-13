# DataEdge Go SDK

Go-клиент для [DataEdge API](https://dataedge.md/).

[![Go Reference](https://pkg.go.dev/badge/github.com/acolev/dataedge-go.svg)](https://pkg.go.dev/github.com/acolev/dataedge-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/acolev/dataedge-go)](https://goreportcard.com/report/github.com/acolev/dataedge-go)

[English version (README.md)](README.md)

## Установка

```bash
go get github.com/acolev/dataedge-go
```

## Быстрый старт

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/acolev/dataedge-go"
)

func main() {
	// Инициализация клиента
	client := dataedge.NewClient("your-api-token",
		dataedge.WithAlphaNameID(123),  // ID отправителя SMS по умолчанию
		dataedge.WithViberNameID(456),  // ID отправителя Viber по умолчанию
	)

	smsService := dataedge.NewSmsService(client)

	// Быстрая отправка SMS
	res, err := smsService.SendQuickSms(context.Background(), "Привет из Go!", "37360000000")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Результат: %v\n", res)
}
```

## Конфигурация клиента

Создайте клиента с помощью `NewClient(token string, opts ...Option)`.

### Опции
- `WithBaseURL(url string)`: Установить пользовательский базовый URL API.
- `WithDebug(debug bool)`: Включить логирование запросов и ответов для отладки.
- `WithAlphaNameID(id int)`: Установить ID отправителя по умолчанию для SMS.
- `WithViberNameID(id int)`: Установить ID отправителя по умолчанию для Viber.
- `WithTranslit(enabled bool)`: Включить/выключить автоматическую транслитерацию (по умолчанию `true`).
- `WithHTTPClient(client *http.Client)`: Использовать пользовательский HTTP-клиент.

---

## SMS Сервис

Доступ к методам через `dataedge.NewSmsService(client)`.

### Методы

| Метод | Описание |
|--------|-------------|
| `GetBalance(ctx)` | Возвращает текущий баланс как `float64`. |
| `GetLimit(ctx)` | Возвращает текущий лимит сообщений. |
| `GetAlphaNames(ctx)` | Возвращает список доступных имен отправителей. |
| `GetAlphaNameID(ctx, name)` | Возвращает ID для конкретного имени отправителя. |
| `CreateSMSMessage(ctx, text)` | Создает шаблон сообщения, возвращает `messageID`. |
| `SendSms(ctx, messageID, phone)` | Отправляет ранее созданное сообщение. |
| `SendQuickSms(ctx, text, phone)` | Создает и отправляет SMS за один вызов. |
| `CheckSMS(ctx, smsID)` | Проверяет статус отправленного SMS. |
| `GetMessagesList(ctx)` | Возвращает историю отправленных сообщений. |

---

## Viber Сервис

Доступ к методам через `dataedge.NewViberService(client)`.

### Методы

| Метод | Описание |
|--------|-------------|
| `CreateViberMessage(ctx, text)` | Создает шаблон Viber сообщения. |
| `SendViberMessage(ctx, phone, messageID)` | Отправляет созданное Viber сообщение. |
| `SendQuickViberMessage(ctx, phone, text)` | Быстрая отправка Viber сообщения. |
| `SendQuickViberMessageWithImage(...)` | Отправка Viber сообщения с изображением, кнопкой и ссылкой. |
| `SendViberMessageList(ctx, name, listID, schedule, extra)` | Отправка сообщения по предустановленному списку. |
| `GetViberMessageList(ctx)` | Возвращает историю Viber сообщений. |

#### Пример отправки Viber с изображением
```go
viberService := dataedge.NewViberService(client)

res, err := viberService.SendQuickViberMessageWithImage(
    ctx,
    "37360000000",
    "Посмотрите это!",
    "IMAGE",           // Тип: MESSAGE или IMAGE
    "/path/to/img.jpg",
    "Открыть ссылку",   // Текст кнопки
    "https://google.com",
)
```

---

## Утилиты

### Транслитерация
Транслитерация включена по умолчанию для SMS. Вы также можете использовать ее вручную:
```go
text := dataedge.Transliterate("Привет мир") // "Privet mir"
```

## Особенности

- **Полное покрытие API**: Поддерживает все методы из PHP SDK.
- **Поддержка Context**: Все методы принимают `context.Context` для контроля таймаутов и отмены.
- **Типобезопасность**: Используются структуры Go и стандартная обработка ошибок.
- **Транслитерация**: Встроенная поддержка транслитерации кириллицы в латиницу.
- **Загрузка файлов (Multipart)**: Поддержка отправки Viber сообщений с изображениями.

## Лицензия

Пользовательская некоммерческая лицензия (см. [LICENSE](LICENSE)). Свободное использование в личных и внутренних целях, перепродажа запрещена.
