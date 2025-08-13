package sync_cmd

import (
	"fmt"

	sync_services "github.com/hurtki/configsManager/services/sync"
	"github.com/spf13/cobra"
)

type SyncLogoutCmd struct {
	SyncDeps *sync_services.Deps
	Command  *cobra.Command
	Dropbox  bool
}

func (c *SyncLogoutCmd) run(cmd *cobra.Command, args []string) error {
	dropboxFlag := c.Dropbox

	fmt.Printf("Flag dropbox: %t\n", dropboxFlag)
	return nil
}

func NewSyncLogoutCmd(d *sync_services.Deps) *SyncLogoutCmd {
	syncLogoutCmd := SyncLogoutCmd{
		SyncDeps: d,
	}

	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Logout from your cloud accounts",
		Long:  ``,
		RunE:  syncLogoutCmd.run,
	}

	cmd.Flags().BoolVar(&syncLogoutCmd.Dropbox, "dropbox", false, "Logout from dropbox sync")

	syncLogoutCmd.Command = cmd

	return &syncLogoutCmd
}
