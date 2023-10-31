package main

import (
	"context"
	"fmt"
	"log"
	"time"

	llamaservergo "github.com/a-h/llama-server-go"
)

func main() {
	ls, err := llamaservergo.New("http://localhost:8080", time.Duration(0))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	resp, err := ls.EmbeddingPost(ctx, llamaservergo.EmbeddingPostRequest{
		Content: "Hello world!",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Response: %+v\n", resp)
}
