package sdk

import "fmt"

func ExampleNewClient() {
	client := NewClient("https://api.example.com", "api-key")
	fmt.Println(client != nil)
	// Output: true
}
