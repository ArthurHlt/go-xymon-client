package client

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type typeRequest string

const (
	statusRequest typeRequest = "status"
	queryRequest  typeRequest = "query"
	pingRequest   typeRequest = "ping"
	eventRequest  typeRequest = "event"
)

type Client struct {
	Host            string
	Port            int
	FQDNEnabled     bool
	TimeoutInSecond int
}

func NewClient(hostname string) *Client {
	client := &Client{
		FQDNEnabled:     true,
		Host:            hostname,
		Port:            1984,
		TimeoutInSecond: 3,
	}
	urlSplit := strings.Split(hostname, ":")
	if len(urlSplit) == 2 {
		client.Host = urlSplit[1]
		port, _ := strconv.Atoi(urlSplit[2])
		client.Port = port
	}
	return client
}
func NewClientFQDNDisabled(hostname string) *Client {
	client := NewClient(hostname)
	client.FQDNEnabled = false
	return client
}
func (c Client) Status(message MessageTest) (string, error) {
	message.FQDNEnabled = c.FQDNEnabled
	return c.sendRequest(statusRequest, message)
}
func (c Client) Query(message MessageTest) (string, error) {
	message.FQDNEnabled = c.FQDNEnabled
	return c.sendRequest(queryRequest, c.filterMessageForQuery(message))
}
func (c Client) Ping() (string, error) {
	return c.sendRequest(pingRequest, "")
}
func (c Client) Event(evt EventTest) (string, error) {
	return c.sendRequest(eventRequest, evt)
}
func (c Client) sendRequest(req typeRequest, data interface{}) (string, error) {
	conn, err := net.DialTimeout("tcp", c.Host+":"+strconv.Itoa(c.Port), time.Duration(c.TimeoutInSecond)*time.Second)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	message := fmt.Sprint(data)
	_, err = conn.Write([]byte(string(req) + message))
	if err != nil {
		return "", err
	}
	// Xymon waiting that write connection has been closed to send response...
	conn.(*net.TCPConn).CloseWrite()

	return bufio.NewReader(conn).ReadString('\n')
}
func (c Client) filterMessageForQuery(message MessageTest) MessageTest {
	return MessageTest{
		Name:  message.Name,
		Host:  message.Host,
		Group: message.Group,
	}
}
