package tracing

import "fmt"

func Init(serviceName string) (func(), error) {
	fmt.Printf("tracing initialized for %s\n", serviceName)
	return func() {}, nil
}
