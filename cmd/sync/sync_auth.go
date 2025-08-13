package sync_cmd

import (
	"fmt"

	sync_services "github.com/hurtki/configsManager/services/sync"
	"github.com/spf13/cobra"
)

type SyncAuthCmd struct {
	SyncDeps *sync_services.Deps
	Command  *cobra.Command
	Dropbox  bool
}

func (c *SyncAuthCmd) run(cmd *cobra.Command, args []string) error {
	dropboxFlag := c.Dropbox
	token, _ := cmd.Flags().GetString("token")

	if dropboxFlag && token != "token_not_given" {
		return fmt.Errorf("connot use --dropbox and --token together")
	}

	fmt.Printf("Flag dropbox: %t\n", dropboxFlag)
	fmt.Printf("Flag token: %s\n", token)
	return nil
}

func NewSyncAuthCmd(d *sync_services.Deps) *SyncAuthCmd {
	syncAuthCmd := SyncAuthCmd{
		SyncDeps: d,
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

	syncAuthCmd.Command = cmd

	return &syncAuthCmd
}
