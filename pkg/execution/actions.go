package execution

import (
	"os"
	"text/tabwriter"
)

func actionInfo(s *Session, _ string) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 4, ' ', 0)
	for _, command := range commands {
		cmd := ":" + command.name
		if command.arg != "" {
			cmd = cmd + " " + command.arg
		}
		w.Write([]byte("    " + cmd + "\t" + command.document + "\n"))
	}
	w.Flush()

	return nil
}

func actionHelp(s *Session, _ string) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 4, ' ', 0)
	for _, command := range commands {
		cmd := ":" + command.name
		if command.arg != "" {
			cmd = cmd + " " + command.arg
		}
		w.Write([]byte("    " + cmd + "\t" + command.document + "\n"))
	}
	w.Flush()

	return nil
}

func actionQuit(s *Session, _ string) error {
	return ErrQuit
}
