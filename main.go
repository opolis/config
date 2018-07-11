package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

const version = "v0.4.0"

func main() {
	// Read secret parameter key from CLI
	if len(os.Args) != 2 {
		fmt.Println("error: invalid parameters, must provide exactly one SSM parameter key name")
		os.Exit(1)
	}

	// AWS session
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// SSM client
	client := ssm.New(sess)

	// Get secret value by key name
	key := os.Args[1]
	resp, err := client.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(true),
	})

	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	// print value to stdout, no newline
	fmt.Print(*(resp.Parameter.Value))
}
