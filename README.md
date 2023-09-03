# Synchro: Timezone-typesafe date and time library for Go

This library is inspired by Rust [chrono](https://github.com/chronotope/chrono)

## Features

- Timezone-typesafe date and time handling
- Easy conversion between time zones
- Support for common date and time operations
- Compatible with the standard `time` package

## Installation

To install Synchro, use `go get`:

    go get github.com/Code-Hex/synchro

## Synopsis

To use Synchro, import it in your Go code:

```go
package main

import (
    "fmt"

    "github.com/Code-Hex/synchro"
    "github.com/Code-Hex/synchro/tz"
)

func main() {
    // The current UTC time is fixed to `2023-09-02 14:00:00`.
    utcNow := synchro.Now[tz.UTC]()
    fmt.Println(utcNow)

    jstNow := synchro.Now[tz.AsiaTokyo]()
    fmt.Println(jstNow)
    // Output:
    // 2009-11-10 23:00:00 +0000 UTC
    // 2009-11-11 08:00:00 +0900 JST
}
```

https://go.dev/play/p/Ql3CM7NLfj0

## Contributing

Contributions to Synchro are welcome!

To contribute, please fork the repository and submit a pull request.


## License

Synchro is licensed under the MIT License. See LICENSE for more information.
