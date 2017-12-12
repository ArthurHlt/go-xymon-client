package main

import (
	"fmt"
	xymclient "github.com/ArthurHlt/go-xymon-client/client"
	"github.com/mitchellh/colorstring"
	"github.com/urfave/cli"
	"net"
	"os"
	"time"
)

var version_major int = 1
var version_minor int = 1
var version_build int = 0

func main() {
	app := cli.NewApp()
	app.Name = "go-xymon-client"
	app.Usage = "A simple cli program to make request on a xymon"
	app.Version = fmt.Sprintf("%d.%d.%d", version_major, version_minor, version_build)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "target, t",
			EnvVar: "XYMON_HOST",
			Usage:  "Target your xymon, e.g: 127.0.0.1:1984",
		},
		cli.BoolFlag{
			Name:  "no-fqdn",
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
					Name:  "color, c",
					Usage: "Color for your test, can be: clear, green, red, yellow, purple or blue",
				},
				cli.StringFlag{
					Name:  "host, h",
					Usage: "Host of your test",
				},
				cli.StringFlag{
					Name:  "name, n",
					Usage: "Name of your test",
				},
				cli.StringFlag{
					Name:  "text, t",
					Usage: "Message to pass in your test",
				},
				cli.StringFlag{
					Name:  "group, g",
					Value: "",
					Usage: "Associate a group to your test (optional)",
				},
				cli.StringFlag{
					Name:  "lifetime, l",
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
					Name:  "host, x",
					Usage: "Host of your test",
				},
				cli.StringFlag{
					Name:  "name, n",
					Usage: "Name of your test",
				},
				cli.StringFlag{
					Name:  "group, g",
					Value: "",
					Usage: "The associate group to your test (optional)",
				},
			},
			Action: query,
		},
		{
			Name:    "ping",
			Aliases: []string{"p"},
			Usage:   "Ping your xymon",
			Action:  ping,
		},
		{
			Name:    "event",
			Aliases: []string{"e"},
			Usage:   "Send test of type event (may not be available in your xymon)",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "host, x",
					Usage: "Host of your test",
				},
				cli.StringFlag{
					Name:  "name, n",
					Usage: "Name of your test",
				},
				cli.StringFlag{
					Name:  "text, t",
					Usage: "Message to pass in your test",
				},
				cli.StringFlag{
					Name:  "id, i",
					Usage: "Id of your event",
				},
				cli.StringFlag{
					Name:  "color, c",
					Usage: "Color for your event, can be: clear, green, red, yellow, purple or blue",
				},
				cli.DurationFlag{
					Name:  "activation, a",
					Usage: "When event will be activated",
				},
				cli.BoolFlag{
					Name:  "ephemeral, e",
					Usage: "set to true to say that event can be destroyed",
				},
				cli.IntFlag{
					Name:  "order, o",
					Usage: "optional, set priority event on others",
					Value: 0,
				},
				cli.BoolFlag{
					Name:  "default, d",
					Usage: "if set to true this will be the default event",
				},
				cli.BoolFlag{
					Name:  "remove, rm",
					Usage: "if set to true this will remove the event",
				},
				cli.DurationFlag{
					Name:  "expiration, exp",
					Usage: "optional, if set event will be removed at that time",
				},
				cli.StringFlag{
					Name:  "time-location, l",
					Usage: "optional, set timezone location to create time format",
				},
			},
			Action: event,
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
	m, err := flagsToMessageTest(c)
	if err != nil {
		return err
	}
	client := getClient(c)
	resp, err := client.Status(m)
	return showResponse(resp, err)
}
func event(c *cli.Context) error {
	_, err := pingWithoutPrint(c)
	if err != nil {
		return err
	}
	if c.String("id") == "" ||
		c.String("host") == "" {
		return cli.NewExitError("ERROR: You must set host and name.", 1)
	}
	if (c.String("name") == "" || c.String("text") == "" || c.String("color") == "") && !c.Bool("remove") {
		return cli.NewExitError("ERROR: You must set text/id/color or mark as remove.", 1)
	}
	t, err := flagsToEventTest(c)
	if err != nil {
		return err
	}
	fmt.Println(t)
	client := getClient(c)
	resp, err := client.Event(t)
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
	m, err := flagsToMessageTest(c)
	if err != nil {
		return err
	}
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
		return "", cli.NewExitError("ERROR: "+err.Error(), 1)
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
func flagsToMessageTest(c *cli.Context) (xymclient.MessageTest, error) {
	color, err := xymclient.ParseColorString(c.String("color"))
	if err != nil {
		return xymclient.MessageTest{}, err
	}
	return xymclient.MessageTest{
		Color:    color,
		Host:     c.String("host"),
		Name:     c.String("name"),
		Text:     c.String("text"),
		Group:    c.String("group"),
		Lifetime: c.String("lifetime"),
	}, nil
}

func flagsToEventTest(c *cli.Context) (xymclient.EventTest, error) {
	color, err := xymclient.ParseColorString(c.String("color"))
	if err != nil {
		return xymclient.EventTest{}, err
	}
	evt := xymclient.EventTest{
		Color:        color,
		Host:         c.String("host"),
		Name:         c.String("name"),
		Text:         c.String("text"),
		Id:           c.String("id"),
		Activation:   time.Now().Add(c.Duration("activation")),
		Ephemeral:    c.Bool("ephemeral"),
		Order:        c.Int("order"),
		Default:      c.Bool("default"),
		TimeLocation: c.String("time-location"),
		Remove:       c.Bool("remove"),
	}
	if c.Duration("expiration") != time.Duration(0) {
		evt.Expiration = time.Now().Add(c.Duration("expiration"))
	}
	return evt, nil
}
