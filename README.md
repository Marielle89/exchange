# exchange

[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](https://github.com/peterhellberg/fixer#license-mit)

Go client for [v6.exchangerate-api.co](https://v6.exchangerate-api.com/) (Foreign exchange rates and currency conversion API)

> You need to [register](https://app.exchangerate-api.com/sign-up) for a free trial api key.

The client loads its api key from the environment variable `EXCHANGE_RATE_API_KEY`.
The function ApiKey allows set new api key.

## Installation

    go get -u github.com/Marielle89/exchange

## Usage examples

**Convert United States Dollar To Euro**

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/Marielle89/exchange"
)

func main() {
	client := exchange.NewClient()

	rate, errRate := client.Rate(context.Background(),
		exchange.USD,
		exchange.UAH,
	)
	if errRate == nil {
		fmt.Println(rate)
	}

	amount, errAmount := client.Amount(context.Background(),
		exchange.USD,
		exchange.UAH,
		155.55,
	)
	
	if errAmount == nil {
		fmt.Println(amount)
	}
}
```