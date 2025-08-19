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
	dropboxFlag := c.Dropbox

	fmt.Printf("Flag dropbox: %t\n", dropboxFlag)

	// здесь нам нужно попросить SyncService удалить токены из памяти, или конкретно какого провайдера
	// LogOut(string) error
	// типо если строка пустая то их всех выходит
	return nil
}

func NewSyncLogoutCmd(syncService sync_services.SyncService) *SyncLogoutCmd {
	syncLogoutCmd := SyncLogoutCmd{
		syncService: syncService,
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
