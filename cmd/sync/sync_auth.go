package sync_cmd

import (
	"fmt"

	sync_services "github.com/hurtki/configsManager/services/sync"
	"github.com/spf13/cobra"
)

type SyncAuthCmd struct {
	syncService sync_services.SyncService
	Command     *cobra.Command
	Dropbox     bool
}

func (c *SyncAuthCmd) run(cmd *cobra.Command, args []string) error {
	if c.Dropbox {
		if err := c.syncService.Auth("dropbox"); err != nil {
			return err
		}
		fmt.Println("Authorized in with dropbox!")
	}
	return nil
}

func NewSyncAuthCmd(syncService sync_services.SyncService) *SyncAuthCmd {
	syncAuthCmd := SyncAuthCmd{
		syncService: syncService,
	}

	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authorize using OAuth2 or Access token",
		Long:  ``,
		RunE:  syncAuthCmd.run,
	}

	cmd.Flags().BoolVar(&syncAuthCmd.Dropbox, "dropbox", false, "Use dropbox as cloud sync provider")
	cmd.Flags().StringP("token", "t", "token_not_given", "To authorize with access token, not refresh")
	cmd.MarkFlagsOneRequired("dropbox")

	syncAuthCmd.Command = cmd

	return &syncAuthCmd
}
