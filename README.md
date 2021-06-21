Go Coinbase Pro [![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/preichenberger/go-coinbasepro) [![Build Status](https://travis-ci.org/preichenberger/go-coinbasepro.svg?branch=master)](https://travis-ci.org/preichenberger/go-coinbasepro)
========

## Summary

Go client for [CoinBase Pro](https://pro.coinbase.com) formerly known as gdax

## Installation
If Go modules are used (Go version >= 1.11) simply import as given beneath.
```sh
go mod init github.com/yourusername/yourprojectname
go get github.com/preichenberger/go-coinbasepro/v2
```

### Older Go versions
```sh
go get github.com/preichenberger/go-coinbasepro/v2
```

### Significant releases
Use [dep](https://github.com/golang/dep) to install previous releases
```sh
dep ensure --add github.com/preichenberger/go-gdax@0.5.7
```

- 0.5.7, last release before rename package to: coinbasepro
- 0.5, as of 0.5 this library uses strings and is not backwards compatible

## Documentation
For full details on functionality, see [GoDoc](http://godoc.org/github.com/preichenberger/go-coinbasepro) documentation.

### Setup
Client will respect environment variables: COINBASE_PRO_BASEURL, COINBASE_PRO_PASSPHRASE, COINBASE_PRO_KEY, COINBASE_PRO_SECRET by default

```go
import (
  coinbasepro "github.com/preichenberger/go-coinbasepro/v2"
)

client := coinbasepro.NewClient()

// optional, configuration can be updated with ClientConfig
client.UpdateConfig(&coinbasepro.ClientConfig{
  BaseURL: "https://api.pro.coinbase.com",
  Key: "coinbase pro key",
  Passphrase: "coinbase pro passphrase",
  Secret: "coinbase pro secret",
})
```

### Sandbox
You can switch to the sandbox URL by setting environment variable: COINBASE_PRO_SANDBOX

Enable sandbox
```sh
export COINBASE_PRO_SANDBOX=1
```

Disable sandbox
```sh
export COINBASE_PRO_SANDBOX=0
```


### HTTP Settings
```go
import (
  "net/http"
  "time"
)

client.HTTPClient = &http.Client {
  Timeout: 15 * time.Second,
}
```

### Decimals
To manage precision correctly, this library sends all price values as strings. It is recommended to use a decimal library
like Spring's [Decimal](https://github.com/shopspring/decimal) if you are doing any manipulation of prices.

Example:
```go
import (
  "github.com/shopspring/decimal"
)

book, err := client.GetBook("BTC-USD", 1)
if err != nil {
    println(err.Error())  
}

lastPrice, err := decimal.NewFromString(book.Bids[0].Price)
if err != nil {
    println(err.Error())  
}

order := coinbasepro.Order{
  Price: lastPrice.Add(decimal.NewFromFloat(1.00)).String(),
  Size: "2.00",
  Side: "buy",
  ProductID: "BTC-USD",
}

savedOrder, err := client.CreateOrder(&order)
if err != nil {
  println(err.Error())
}

println(savedOrder.ID)
```

### Retry
You can set a retry count which uses exponential backoff: (2^(retry_attempt) - 1) / 2 * 1000 * milliseconds
```
client.RetryCount = 3 # 500ms, 1500ms, 3500ms
```

### Cursor
This library uses a cursor pattern so you don't have to keep track of pagination.

```go
var orders []coinbasepro.Order
cursor = client.ListOrders()

for cursor.HasMore {
  if err := cursor.NextPage(&orders); err != nil {
    println(err.Error())
    return
  }

  for _, o := range orders {
    println(o.ID)
  }
}

```

### Websockets
Listen for websocket messages

```go
  import(
    ws "github.com/gorilla/websocket"
  )

  var wsDialer ws.Dialer
  wsConn, _, err := wsDialer.Dial("wss://ws-feed.pro.coinbase.com", nil)
  if err != nil {
    println(err.Error())
  }

  subscribe := coinbasepro.Message{
    Type:      "subscribe",
    Channels: []coinbasepro.MessageChannel{
      coinbasepro.MessageChannel{
        Name: "heartbeat",
        ProductIds: []string{
          "BTC-USD",
        },
      },
      coinbasepro.MessageChannel{
        Name: "level2",
        ProductIds: []string{
          "BTC-USD",
        },
      },
    },
  }
  if err := wsConn.WriteJSON(subscribe); err != nil {
    println(err.Error())
  }

  for true {
    message := coinbasepro.Message{}
    if err := wsConn.ReadJSON(&message); err != nil {
      println(err.Error())
      break
    }

    println(message.Type)
  }

```

### Time
Results return coinbase time type which handles different types of time parsing that coinbasepro returns. This wraps the native go time type

```go
  import(
    "time"
    coinbasepro "github.com/preichenberger/go-coinbasepro/v2"
  )

  coinbaseTime := coinbasepro.Time{}
  println(time.Time(coinbaseTime).Day())
```

### Examples
This library supports all public and private endpoints

Get Accounts:
```go
  accounts, err := client.GetAccounts()
  if err != nil {
    println(err.Error())
  }

  for _, a := range accounts {
    println(a.Balance)
  }
```

List Account Ledger:
```go
  var ledgers []coinbasepro.LedgerEntry

  accounts, err := client.GetAccounts()
  if err != nil {
    println(err.Error())
  }

  for _, a := range accounts {
    cursor := client.ListAccountLedger(a.ID)
    for cursor.HasMore {
      if err := cursor.NextPage(&ledgers); err != nil {
        println(err.Error())
      }

      for _, e := range ledgers {
        println(e.Amount)
      }
    }
  }
```

Create an Order:
```go
  order := coinbasepro.Order{
    Price: "1.00",
    Size: "1.00",
    Side: "buy",
    ProductID: "BTC-USD",
  }

  savedOrder, err := client.CreateOrder(&order)
  if err != nil {
    println(err.Error())
  }

  println(savedOrder.ID)
```

Transfer funds:
```go
  transfer := coinbasepro.Transfer {
    Type: "deposit",
    Amount: "1.00",
  }

  savedTransfer, err := client.CreateTransfer(&transfer)
  if err != nil {
    println(err.Error())
  }
```

Get Trade history:
```go
  var trades []coinbasepro.Trade
  cursor := client.ListTrades("BTC-USD")

  for cursor.HasMore {
    if err := cursor.NextPage(&trades); err != nil {
      for _, t := range trades {
        println(trade.CoinbaseID)
      }
    }
  }
```

### Testing
To test with Coinbase's public sandbox set the following environment variables:
```sh
export COINBASE_PRO_KEY="sandbox key"
export COINBASE_PRO_PASSPHRASE="sandbox passphrase"
export COINBASE_PRO_SECRET="sandbox secret"
```

Then run `go test`
```sh
go test
```

Note that your sandbox account will need at least 2,000 USD and 2 BTC to run the tests.
