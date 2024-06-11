package mqttclient

import "fmt"

type mockClient struct {
}

func (c *mockClient) sub() {
	fmt.Println("subscribed mock mqtt client to all topics")
}

func (c *mockClient) Disconnect() {
	fmt.Println("diconnected mock mqtt client")
}

//func ()
