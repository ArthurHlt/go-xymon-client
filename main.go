package main

import (
	"fmt"
	"github.com/urfave/cli"
	"github.com/mitchellh/colorstring"
	"os"
	xymclient "github.com/ArthurHlt/go-xymon-client/client"
	"net"
)

var version_major int = 1
var version_minor int = 0
var version_build int = 0

func main() {
	app := cli.NewApp()
	app.Name = "xymon-client"
	app.Usage = "A simple cli program to make request on a xymon"
	app.Version = fmt.Sprintf("%d.%d.%d", version_major, version_minor, version_build)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "target, t",
			EnvVar: "XYMON_HOST",
			Usage: "Target your xymon, e.g: 127.0.0.1:1984",
		},
		cli.BoolFlag{
			Name: "no-fqdn",
			Usage: "Do not using fqdn",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "status",
			Aliases: []string{"q"},
			Usage:   "Send a test (to update or create it) to your xymon",
			Action:  status,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "color, c",
					Usage: "Color for your test, can be: clear, green, red, yellow, purple or blue",
				},
				cli.StringFlag{
					Name: "host, h",
					Usage: "Host of your test",
				},
				cli.StringFlag{
					Name: "name, n",
					Usage: "Name of your test",
				},
				cli.StringFlag{
					Name: "text, t",
					Usage: "Message to pass in your test",
				},
				cli.StringFlag{
					Name: "group, g",
					Value: "",
					Usage: "Associate a group to your test (optional)",
				},
				cli.StringFlag{
					Name: "lifetime, l",
					Value: "",
					Usage: `Set the expiration time of your test (optional) (add "h" (hours), "d" (days) or "w" (weeks) immediately after the number to use instead of minute)`,
				},
			},
		},
		{
			Name:    "query",
			Aliases: []string{"s"},
			Usage:   "Get the current status of a test",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "host, x",
					Usage: "Host of your test",
				},
				cli.StringFlag{
					Name: "name, n",
					Usage: "Name of your test",
				},
				cli.StringFlag{
					Name: "group, g",
					Value: "",
					Usage: "The associate group to your test (optional)",
				},
			},
			Action:  query,
		},
		{
			Name:    "ping",
			Aliases: []string{"p"},
			Usage:   "Ping your xymon",
			Action:  ping,
		},
	}

	app.Run(os.Args)
}
func getClient(c *cli.Context) *xymclient.Client {
	var client *xymclient.Client
	if c.GlobalBool("no-fqdn") {
		client = xymclient.NewClientFQDNDisabled(c.GlobalString("target"))
	} else {
		client = xymclient.NewClient(c.GlobalString("target"))
	}
	return client
}
func status(c *cli.Context) error {
	_, err := pingWithoutPrint(c)
	if err != nil {
		return err
	}
	if c.String("host") == "" || c.String("name") == "" || c.String("text") == "" || c.String("color") == "" {
		return cli.NewExitError("ERROR: You must set host, name, text and color.", 1)
	}
	m := flagsToMessageTest(c)
	client := getClient(c)
	resp, err := client.Status(m)
	return showResponse(resp, err)
}
func query(c *cli.Context) error {
	_, err := pingWithoutPrint(c)
	if err != nil {
		return err
	}
	if c.String("host") == "" || c.String("name") == "" {
		return cli.NewExitError("ERROR: You must set host and name.", 1)
	}
	m := flagsToMessageTest(c)
	client := getClient(c)
	resp, err := client.Query(m)
	return showResponse(resp, err)
}
func pingWithoutPrint(c *cli.Context) (string, error) {
	client := getClient(c)
	resp, err := client.Ping()
	if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
		fmt.Println("INFO: We had a timeout during getting response. No response will be provided.")
		return "", nil
	}
	if _, ok := err.(net.Error); !ok && err != nil {
		return "", cli.NewExitError("ERROR: " + err.Error(), 1)
	}
	if resp == "" {
		return "", cli.NewExitError("ERROR: The server is not answering. Probably connection has been shut down by tier (firewall can block access).", 1)
	}
	return formatRepsponse(resp), nil
}
func ping(c *cli.Context) error {
	resp, err := pingWithoutPrint(c)
	if err != nil {
		return err
	}
	fmt.Print(resp)
	return nil
}
func formatRepsponse(resp string) string {
	return colorstring.Color("\n[blue]RESPONSE: [reset]" + resp)
}
func showResponse(resp string, err error) error {
	if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
		fmt.Println("We had a timeout during getting response. No response will be provided.")
		return nil
	}
	if _, ok := err.(net.Error); !ok && err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Print(formatRepsponse(resp))
	return nil
}
func flagsToMessageTest(c *cli.Context) xymclient.MessageTest {
	return xymclient.MessageTest{
		Color: xymclient.ColorTest(c.String("color")),
		Host: c.String("host"),
		Name: c.String("name"),
		Text: c.String("text"),
		Group: c.String("group"),
		Lifetime: c.String("lifetime"),
	}
}