# ObjectID

[![GoDoc](https://godoc.org/github.com/pkg-id/objectid?status.svg)](https://godoc.org/github.com/pkg-id/objectid)
[![Go Report Card](https://goreportcard.com/badge/github.com/pkg-id/objectid)](https://goreportcard.com/report/github.com/pkg-id/objectid)

ObjectID is a MongoDB ObjectID implementation for Go. See the [MongoDB ObjectID documentation](https://docs.mongodb.com/manual/reference/method/ObjectId/) for more information.


## Installation

```bash
go get github.com/pkg-id/objectid
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/pkg-id/objectid"
)

func main() {
	id := objectid.New() 
	fmt.Println(id)
}
```
