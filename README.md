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
	err := icloudgo.Login(username, password)
	if err != nil {
		fmt.Println(err)
	} else {
		contacts := icloudgo.GetContacts()
		fmt.Println(contacts)
	}
}
```
