# Synchro: Timezone-typesafe date and time library for Go

[![test](https://github.com/Code-Hex/synchro/actions/workflows/test.yml/badge.svg)](https://github.com/Code-Hex/synchro/actions/workflows/test.yml) [![codecov](https://codecov.io/gh/Code-Hex/synchro/graph/badge.svg?token=VWPbmNRHw8)](https://codecov.io/gh/Code-Hex/synchro) [![Go Reference](https://pkg.go.dev/badge/github.com/Code-Hex/synchro/.svg)](https://pkg.go.dev/github.com/Code-Hex/synchro/)

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

Please refer to the numerous usage examples on [GoDoc](https://pkg.go.dev/github.com/Code-Hex/synchro/) for reference.

## Utilities

We also have a wide range of very useful utilities!!

If you have a feature request, please open an issue. It would be great if you could provide relevant examples or links that could be helpful.

- [In](https://pkg.go.dev/github.com/Code-Hex/synchro#In)
- [ConvertTz](https://pkg.go.dev/github.com/Code-Hex/synchro#ConvertTz)
- [NowContext](https://pkg.go.dev/github.com/Code-Hex/synchro#NowContext)
- [Quarter](https://pkg.go.dev/github.com/Code-Hex/synchro#Quarter)
- [Semester](https://pkg.go.dev/github.com/Code-Hex/synchro#Semester)
- [StartOfMonth](https://pkg.go.dev/github.com/Code-Hex/synchro#Time.StartOfMonth)
- [EndOfMonth](https://pkg.go.dev/github.com/Code-Hex/synchro#Time.EndOfMonth)
- [StartOfQuarter](https://pkg.go.dev/github.com/Code-Hex/synchro#Time.StartOfQuarter)
- [EndOfQuarter](https://pkg.go.dev/github.com/Code-Hex/synchro#Time.EndOfQuarter)
- [StartOfSemester](https://pkg.go.dev/github.com/Code-Hex/synchro#Time.StartOfSemester)
- [EndOfSemester](https://pkg.go.dev/github.com/Code-Hex/synchro#Time.EndOfSemester)
- [StartOfWeek](https://pkg.go.dev/github.com/Code-Hex/synchro#Time.StartOfWeek)
- [EndOfWeek](https://pkg.go.dev/github.com/Code-Hex/synchro#Time.EndOfWeek)
- [StartOfYear](https://pkg.go.dev/github.com/Code-Hex/synchro#Time.StartOfYear)
- [EndOfYear](https://pkg.go.dev/github.com/Code-Hex/synchro#Time.EndOfYear)
- [IsBetween](https://pkg.go.dev/github.com/Code-Hex/synchro#Time.IsBetween)
- [IsLeapYear](https://pkg.go.dev/github.com/Code-Hex/synchro#Time.IsLeapYear)

## Contributing

Contributions to Synchro are welcome!

To contribute, please fork the repository and submit a pull request.


## License

Synchro is licensed under the MIT License. See LICENSE for more information.
