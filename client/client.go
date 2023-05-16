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
	eventDeleteRequest  typeRequest = "event eventdelete\n"
)

type Client interface {
	Status(MessageTest) (string, error)
	Query(MessageTest) (string, error)
	Ping() (string, error)
	Event(EventTest) (string, error)
	EventDelete(EventTest) (string, error)
}

type XymonClient struct {
	Host            string
	Port            int
	FQDNEnabled     bool
	TimeoutInSecond int
}

func NewClient(hostname string) Client {
	c := &XymonClient{
		FQDNEnabled:     true,
		Host:            hostname,
		Port:            1984,
		TimeoutInSecond: 3,
	}
	urlSplit := strings.Split(hostname, ":")
	if len(urlSplit) == 2 {
		c.Host = urlSplit[1]
		port, _ := strconv.Atoi(urlSplit[2])
		c.Port = port
	}
	return c
}
func NewClientFQDNDisabled(hostname string) Client {
	c := NewClient(hostname)
	c.(*XymonClient).FQDNEnabled = false
	return c
}
func (c XymonClient) Status(message MessageTest) (string, error) {
	message.FQDNEnabled = c.FQDNEnabled
	return c.sendRequest(statusRequest, message)
}
func (c XymonClient) Query(message MessageTest) (string, error) {
	message.FQDNEnabled = c.FQDNEnabled
	return c.sendRequest(queryRequest, c.filterMessageForQuery(message))
}
func (c XymonClient) Ping() (string, error) {
	return c.sendRequest(pingRequest, "")
}
func (c XymonClient) Event(evt EventTest) (string, error) {
	return c.sendRequest(eventRequest, evt)
}
func (c XymonClient) EventDelete(evt EventTest) (string, error) {
	return c.sendRequest(eventDeleteRequest, evt)
}
func (c XymonClient) sendRequest(req typeRequest, data interface{}) (string, error) {
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
func (c XymonClient) filterMessageForQuery(message MessageTest) MessageTest {
	return MessageTest{
		Name:  message.Name,
		Host:  message.Host,
		Group: message.Group,
	}
}
