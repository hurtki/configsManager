package sync_cmd

import (
	"fmt"

	sync_services "github.com/hurtki/configsManager/services/sync"
	"github.com/spf13/cobra"
)

type SyncLogoutCmd struct {
	syncService sync_services.SyncService
	Command     *cobra.Command
	Dropbox     bool
}

func (c *SyncLogoutCmd) run(cmd *cobra.Command, args []string) error {
	if c.Dropbox {
		if err := c.syncService.Logout("dropbox"); err != nil {
			return err
		}
		fmt.Println("Deleted dropbox access and refresh tokens from system successfully!")
	} else {
		if err := c.syncService.Logout(""); err != nil {
			return err
		}
		fmt.Println("Deleted all access and refresh tokens from system successfully!")
	}

	return nil
}

func NewSyncLogoutCmd(syncService sync_services.SyncService) *SyncLogoutCmd {
	syncLogoutCmd := SyncLogoutCmd{
		syncService: syncService,
	}

	cmd := &cobra.Command{
        Use:   "logout",
        Short: "Logout from your cloud accounts",
        Long: `Logout from cloud services configured with ConfigsManager.
Supports logging out from specific providers (e.g., Dropbox) or all at once.`,
        RunE: syncLogoutCmd.run,
    }

	cmd.Flags().BoolVar(&syncLogoutCmd.Dropbox, "dropbox", false, "Logout from dropbox sync")

	syncLogoutCmd.Command = cmd

	return &syncLogoutCmd
}
