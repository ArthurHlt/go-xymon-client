package client

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type colorTest string

const (
	CNan    colorTest = "blue"
	CInf    colorTest = "purple"
	CClear  colorTest = "clear"
	CGreen  colorTest = "green"
	CRed    colorTest = "red"
	CYellow colorTest = "yellow"
)

func (c colorTest) String() string {
	return string(c)
}

func ParseColorString(s string) (colorTest, error) {
	color := colorTest(strings.ToLower(s))
	switch color {
	case CNan, CInf, CClear, CGreen, CRed, CYellow:
		return color, nil
	}
	return colorTest(""), errors.New("Color " + s + " doesn't exists.")
}

type MessageTest struct {
	Color       colorTest // optional when querying
	Host        string
	Name        string
	Text        string // optional when querying
	Group       string // optional
	Lifetime    string // optional, default in minutes (add "h" (hours), "d" (days) or "w" (weeks) immediately after the number to use instead of minute)
	FQDNEnabled bool
}

func (message MessageTest) String() string {
	var msg string
	if message.Lifetime != "" {
		msg += "+" + message.Lifetime
	}
	if message.Group != "" {
		msg += "/group:" + message.Group
	}
	if message.FQDNEnabled {
		msg += " " + strings.Replace(message.Host, ".", ",", -1)
	} else {
		msg += " " + message.Host
	}
	msg += "." + message.Name
	if message.Color != "" {
		msg += " " + string(message.Color)
	}
	if message.Text != "" {
		msg += " " + message.Text
	}
	return msg
}

type EventTest struct {
	Name         string
	Id           string // event id (this is required)
	Color        colorTest
	Host         string
	Activation   time.Time // When activate the event (default: now)
	Ephemeral    bool      // set to true to say that event can be destroyed
	Text         string
	Order        int       // optional, set priority event on others
	Default      bool      // if set to true this will be the default event
	Expiration   time.Time // optional, if set event will be removed at that time
	TimeLocation string    // optional, set timezone location to create time format
	Remove       bool      // optional, if true this will remove event
}

func (t EventTest) String() string {
	activation := t.Activation
	if activation == (time.Time{}) {
		activation = time.Now()
	}
	expiration := t.Expiration
	if t.TimeLocation != "" {
		loc, _ := time.LoadLocation(t.TimeLocation)
		activation = activation.In(loc)
		if expiration != (time.Time{}) {
			expiration = expiration.In(loc)
		}
	}
	if t.Order <= 0 {
		t.Order = 3000
	}
	persistence := "per"
	if t.Ephemeral {
		persistence = "eph"
	}
	msg := fmt.Sprintf(
		"activation: %d\ncolor: %s\nhost: %s\nid: %s\nmessage: %s\norder: %d\npersistence: %s\nservice: %s\n",
		activation.Unix(),
		string(t.Color),
		t.Host,
		t.Id,
		t.Text,
		t.Order,
		persistence,
		t.Name,
	)
	if expiration != (time.Time{}) {
		msg += fmt.Sprintf("expiration: %d\n", expiration.Unix())
	}
	if t.Default {
		msg += "default: YES\n"
	}
	if t.Remove {
		msg += "suppress: YES\n"
	}
	return msg
}
