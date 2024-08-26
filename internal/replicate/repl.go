package replicate

import (
	"fmt"

	"github.com/replicate/replicate-go"
)

var ReplicateClient *replicate.Client

func ReplicateInit() error {
	fmt.Println("Initializing Replicate client...")
	r8, err := replicate.NewClient(replicate.WithTokenFromEnv())
	if err != nil {
		return err
	}

	ReplicateClient = r8

	return nil
}
