package sync_cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type SyncLockCmd struct {
	Command *cobra.Command
}

func (c *SyncLockCmd) run(cmd *cobra.Command, args []string) error {
	fmt.Fprintln(os.Stderr, "Session locked. Apply with: eval \"$(cm sync lock)\"")
	fmt.Println("unset CM_SYNC_KEY")
	return nil
}

func NewSyncLockCmd() *SyncLockCmd {
	syncLockCmd := SyncLockCmd{}

	cmd := &cobra.Command{
		Use:   "lock",
		Short: "Print `unset CM_SYNC_KEY` for the current shell",
		Long: `Prints a shell-evaluable line that clears the unlocked session key from
the current shell. After this, cm commands will prompt for the
passphrase again.

Usage:

    eval "$(cm sync lock)"`,
		RunE: syncLockCmd.run,
	}

	syncLockCmd.Command = cmd
	return &syncLockCmd
}
