# go-xymon-client

A client for [xymon](http://xymon.sourceforge.net/) wrote for golang.

See usage in [/client/client.go](/client/client.go).

This small client have only 3 verbs created:
- ping
- query
- status

No other verbs will be created, if you want more pull requests are accepted.

This client is also followed by a simple cli program.


## Cli

```
NAME:
    xymon-client - A simple cli program to make request on a xymon
 
 USAGE:
    go-xymon-client [global options] command [command options] [arguments...]
 
 VERSION:
    1.0.0
 
 COMMANDS:
      status, q  Send a test (to update or create it) to your xymon
      query, s   Get the current status of a test
      ping, p    Ping your xymon
      help, h    Shows a list of commands or help for one command
 
 GLOBAL OPTIONS:
    --target value, -t value  Target your xymon, e.g: 127.0.0.1:1984 [$XYMON_HOST]
    --no-fqdn                 Do not using fqdn
    --help, -h                show help
    --version, -v             print the version
```

## ping

```
NAME:
   go-xymon-client ping - Ping your xymon

USAGE:
   go-xymon-client ping
```

## query

```
NAME:
   go-xymon-client query - Get the current status of a test

USAGE:
   go-xymon-client query [command options] [arguments...]

OPTIONS:
   --host value, -x value   Host of your test
   --name value, -n value   Name of your test
   --group value, -g value  The associate group to your test (optional)
```

## status

```
NAME:
   go-xymon-client status - Send a test (to update or create it) to your xymon

USAGE:
   go-xymon-client status [command options] [arguments...]

OPTIONS:
   --color value, -c value     Color for your test, can be: clear, green, red, yellow, purple or blue
   --host value, -h value      Host of your test
   --name value, -n value      Name of your test
   --text value, -t value      Message to pass in your test
   --group value, -g value     Associate a group to your test (optional)
   --lifetime value, -l value  Set the expiration time of your test (optional) (add "h" (hours), "d" (days) or "w" (weeks) immediately after the number to use instead of minute)
```