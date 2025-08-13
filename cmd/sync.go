package cmd

import (
	"github.com/hurtki/configsManager/cmd/sync"
	"github.com/hurtki/configsManager/services"
	syncServices "github.com/hurtki/configsManager/services/sync"
	"github.com/spf13/cobra"
)

type SyncCmd struct {
	Command *cobra.Command
}

func NewSyncCmd(AppConfigService services.AppConfigService,
	ConfigsListService services.ConfigsListService,
	OsService services.OsService,
) *SyncCmd {
	syncCmd := SyncCmd{}

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "sync [push/pull/auth/logout]",
		Long:  ``,
	}

	SyncDeps := syncServices.Deps{
		AppConfigService:   AppConfigService,
		ConfigsListService: ConfigsListService,
		OsService:          OsService,
	}

	cmd.AddCommand(sync_cmd.NewSyncAuthCmd(&SyncDeps).Command)
	cmd.AddCommand(sync_cmd.NewSyncLogoutCmd(&SyncDeps).Command)
	syncCmd.Command = cmd

	return &syncCmd
}
