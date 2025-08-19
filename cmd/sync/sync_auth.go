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
	dropboxFlag := c.Dropbox
	token, _ := cmd.Flags().GetString("token")

	if (!dropboxFlag) && token != "token_not_given" {
		return ErrAuthTokeWithoutProvider
	}

	fmt.Printf("Flag dropbox: %t\n", dropboxFlag)
	fmt.Printf("Flag token: %s\n", token)
	if dropboxFlag && (token == "token_not_given") {
		if err := c.syncService.Auth("dropbox", ""); err != nil {
			return err
		}
	} else if dropboxFlag && (token != "token_not_given") {
		if err := c.syncService.Auth("dropbox", token); err != nil {
			return err
		}
	}

	// здесь обращаемся к SyncService с методом Auth(provider string, token string)
	// там уже понимают если токен пустой то надо по факту авторизировать
	// ну и собственно там и запоминают токен в память в связку ключей
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

	syncAuthCmd.Command = cmd

	return &syncAuthCmd
}
