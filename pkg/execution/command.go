package execution

type command struct {
	name     string
	action   func(*Session, string) error
	complete func(*Session, string) []string
	arg      string
	document string
}

var commands []command

func init() {
	commands = []command{
		{
			name:     "info",
			action:   info,
			document: "show the machine information",
		},
		{
			name:     "help",
			action:   actionHelp,
			document: "show this help",
		},
		{
			name:     "quit",
			action:   actionQuit,
			document: "quit the session",
		},
	}
}
