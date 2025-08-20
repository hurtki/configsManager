package cmd

import (
	"github.com/hurtki/configsManager/cmd/sync"
	"github.com/hurtki/configsManager/services"
	"github.com/hurtki/configsManager/services/sync"
	"github.com/spf13/cobra"
)

type SyncCmd struct {
	Command *cobra.Command
}

func NewSyncCmd(AppConfigService services.AppConfigService,
	ConfigsListService services.ConfigsListService,
	OsService services.OsService,
	SyncService sync_services.SyncService,
) *SyncCmd {
	syncCmd := SyncCmd{}

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "sync [push/pull/auth/logout]",
		Long:  ``,
	}

	cmd.AddCommand(sync_cmd.NewSyncAuthCmd(SyncService).Command)
	cmd.AddCommand(sync_cmd.NewSyncLogoutCmd(SyncService).Command)
	cmd.AddCommand(sync_cmd.NewSyncPushCmd(SyncService, ConfigsListService, OsService).Command)
	cmd.AddCommand(sync_cmd.NewSyncPullCmd(SyncService, OsService).Command)
	syncCmd.Command = cmd

	return &syncCmd
}
