package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/urfave/cli"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	app := cli.NewApp()
	app.Name = "config"
	app.Usage = "set and get encrypted values from AWS SSM"
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "profile",
			Usage: "AWS `PROFILE` in ~/.aws/credentials",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "set",
			Usage: "set a secret value in SSM - set <name> <kms key id>",
			Action: enforceSession(func(c *cli.Context, client *ssm.SSM) error {
				if !isEven(c.NArg()) {
					return cli.NewExitError("arguments must be pairs of <name> <key id>", 1)
				}

				// grab each pair of <name> <key id> and prompt for value
				for i := 0; i < (c.NArg() - 1); i += 2 {
					name := c.Args().Get(i)
					keyId := c.Args().Get(i + 1)

					// prompt for value
					fmt.Printf("%s: ", name)
					bytes, _ := terminal.ReadPassword(0)
					fmt.Println()

					// encrypt to SSM
					_, err := client.PutParameter(&ssm.PutParameterInput{
						KeyId: aws.String(keyId),
						Type:  aws.String("SecureString"),
						Name:  aws.String(name),
						Value: aws.String(string(bytes)),
					})

					if err != nil {
						return cli.NewExitError(
							fmt.Sprintf("could not set parameter %s - %s", name, err.Error()), 1)
					}
				}

				return nil
			}),
		},
		{
			Name:  "get",
			Usage: "get a secret value in SSM - get <name1> <name2> ... <nameN>",
			Action: enforceSession(func(c *cli.Context, client *ssm.SSM) error {
				if c.NArg() == 0 {
					return cli.NewExitError("get needs at least one parameter name to fetch", 1)
				}

				for _, name := range c.Args() {
					resp, err := client.GetParameter(&ssm.GetParameterInput{
						Name:           aws.String(name),
						WithDecryption: aws.Bool(true),
					})

					if err != nil {
						return cli.NewExitError(err.Error(), 1)
					}

					fmt.Println(*(resp.Parameter.Value))
				}

				return nil
			}),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func enforceSession(action func(*cli.Context, *ssm.SSM) error) cli.ActionFunc {
	return func(c *cli.Context) error {
		defer func() {
			if r := recover(); r != nil {
				log.Fatal(r)
			}
		}()

		// AWS session
		var sess *session.Session
		profile := c.GlobalString("profile")

		if profile != "" {
			// look for profile by name
			sess = session.Must(session.NewSessionWithOptions(session.Options{
				Profile:           profile,
				SharedConfigState: session.SharedConfigEnable,
			}))
		} else {
			// no profile provided, lookup by env vars
			sess = session.Must(session.NewSession())
		}

		// SSM client
		client := ssm.New(sess)

		return action(c, client)
	}
}

func isEven(v int) bool {
	return v%2 == 0
}
