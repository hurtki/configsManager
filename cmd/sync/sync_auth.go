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
		Short: "Authorize using OAuth2", 
		Long: `Authorize your ConfigsManager account using OAuth2 (Dropbox). 
		This command sets up authentication
		so that you can push/pull configurations to/from the cloud.`, 
		RunE: syncAuthCmd.run, 
	}


	cmd.Flags().BoolVar(&syncAuthCmd.Dropbox, "dropbox", false, "Use dropbox as cloud sync provider")
	cmd.MarkFlagsOneRequired("dropbox")

	syncAuthCmd.Command = cmd

	return &syncAuthCmd
}
