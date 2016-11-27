# icloudgo
Fetch contacts from iCloud using Go.
## Usage
```go
package main

import (
	"github.com/oziasg/icloudgo"
	"fmt"
)

func main() {
	username := "example@example.com"
	password := "example"
	icloudgo.Login(username, password)
	contacts := icloudgo.GetContacts()
	fmt.Println(contacts)
}
```
