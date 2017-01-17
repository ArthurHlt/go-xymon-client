package client

type ColorTest string

const (
	CNan ColorTest = "blue"
	CInf ColorTest = "purple"
	CClear ColorTest = "clear"
	CGreen ColorTest = "green"
	CRed ColorTest = "red"
	CYellow ColorTest = "yellow"
)

type MessageTest struct {
	Color    ColorTest // optional when querying
	Host     string
	Name     string
	Text     string    // optional when querying
	Group    string    // optional
	Lifetime string    // optional, default in minutes (add "h" (hours), "d" (days) or "w" (weeks) immediately after the number to use instead of minute)
}