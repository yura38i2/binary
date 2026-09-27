# binary

Binary serialization for Go structures.

The package provides a simple way to serialize Go values and structures into a binary representation and deserialize them back.

## Installation

```bash
go get github.com/yura38i2/binary
```

## Usage

Import the package:

```go
import "github.com/yura38i2/binary"
```

Example:

```go
package main

import (
	"fmt"
	"log"

	"github.com/yura38i2/binary"
)

type TestStruct struct {
	ID    uint32
	Score int32
	Name  string
}

func main() {
	marshaled := TestStruct{
		ID:    123,
		Score: 100,
		Name:  "Player",
	}

	data, err := binary.Marshal(marshaled)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Encoded: %v\n", data)

	var unmarshaled TestStruct

	err = binary.Unmarshal(data, &unmarshaled)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Decoded: %+v\n", unmarshaled)
}
```

## Features

* Serialization of Go structures
* Primitive field types
* Strings
* Arrays and slices
* Arrays/slices of structures
* Support for `byte` fields
* Deserialization back into Go values

## Requirements

* Go 1.24 or newer

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

