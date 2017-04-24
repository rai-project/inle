package cmd

import (
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/Unknwon/com"
	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	"github.com/rai-project/inle/pkg/kernel"
	"github.com/spf13/cobra"
)

// notebookCmd represents the notebook command
var notebookCmd = &cobra.Command{
	Use: "notebook",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.Errorf("invalid number of arguments. Expecting 1 argument, but got %d", len(args))
		}
		connectionFile := args[0]
		if !com.IsFile(connectionFile) {
			return errors.Errorf("cannot locate the connection file %v on disk", connectionFile)
		}
		f, err := os.Open(connectionFile)
		if err != nil {
			return errors.Wrapf(err, "unable to open connection file %v", connectionFile)
		}

		id := getID(connectionFile)

		k, err := kernel.New(id, f)
		f.Close()

		if err != nil {
			return err
		}

		// Without signal handling, Go will exit on signal, even if the signal was caught by ZeroMQ
		chSignal := make(chan os.Signal, 1)
		signal.Notify(chSignal, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		go func() {
			<-chSignal
			k.Shutdown()
		}()

		return k.Run()
	},
}

func getID(c string) string {
	base := filepath.Base(c)
	return strings.TrimSuffix(strings.TrimPrefix(base, "kernel-"), ".json")
}

func init() {
	pp.WithLineInfo = true
	RootCmd.AddCommand(notebookCmd)
}
