# go-blih

**go-blih** is a Go client library for Blih (the Epitech bocal API).

## Installation

You can use `go get` to install the latest version of the library. This command will install the the library and its dependencies:

```go
go get -u github.com/ShellBear/go-blih/blih
```

## Examples

#### Generate new Service
```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ShellBear/go-blih/blih"
)

func main() {
	svc := blih.New("EPITECH_LOGIN", "EPITECH_PASSWORD", context.Background())

	response, err := svc.Utils.WhoAmI()
	if err != nil {
		fmt.Println("An error ocurred. Err:", err)
		os.Exit(1)
	}

	fmt.Printf("I am %s\n", response.Message)

}
```

```bash
> go build -o go-blih .
> ./go-blih
I am YOUR_EPITECH_LOGIN
```

#### Create repository

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ShellBear/go-blih/blih"
)

func main() {
	svc := blih.New("EPITECH_LOGIN", "EPITECH_PASSWORD", context.Background())

	repo := &blih.Repository{Name: "example"}

	response, err := svc.Repository.Create(repo)
	if err != nil {
		fmt.Println("An error ocurred. Err:", err)
		os.Exit(1)
	}

	fmt.Println(response.Message)
}

```

#### List repositories

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ShellBear/go-blih/blih"
)

func main() {
	svc := blih.New("EPITECH_LOGIN", "EPITECH_PASSWORD", context.Background())

	response, err := svc.Repository.List()
	if err != nil {
		fmt.Println("An error ocurred. Err:", err)
		os.Exit(1)
	}

	for repoName, repo := range response.Repositories {
		fmt.Printf("%s (%s)\n", repoName, repo.URL)
	}
}
```

## Documentation

A generated documentation is available on [GoDoc](https://godoc.org/github.com/ShellBear/go-blih).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details