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
	Token       string
}

func (c *SyncAuthCmd) run(cmd *cobra.Command, args []string) error {

	if (!c.Dropbox) && c.Token != "token_not_given" {
		return ErrAuthTokeWithoutProvider
	}

	if c.Dropbox && (c.Token == "token_not_given") {
		if err := c.syncService.Auth("dropbox", ""); err != nil {
			return err
		}
		fmt.Println("Authorized in with dropbox!")
	} else if c.Dropbox && (c.Token != "token_not_given") {
		if err := c.syncService.Auth("dropbox", c.Token); err != nil {
			return err
		}
		fmt.Println("Authorized in with dropbox token!")
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
	cmd.MarkFlagsOneRequired("token", "dropbox")
	syncAuthCmd.Token, _ = cmd.Flags().GetString("token")

	syncAuthCmd.Command = cmd

	return &syncAuthCmd
}
