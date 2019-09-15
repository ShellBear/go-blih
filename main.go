package main

import (
	"context"
	"fmt"

	"github.com/ShellBear/go-blih/pkg/blih"
)

func main() {
	svc := blih.New("", "", context.Background(), blih.DefaultApiBaseURL)

	keys, err := svc.SSHKey.List()
	if err != nil {
		panic(err)
	}

	for key, _ := range *keys {
		fmt.Println(key)
	}
}
