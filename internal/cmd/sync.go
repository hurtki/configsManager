package cmd

import (
	sync_cmd "github.com/hurtki/configsManager/internal/cmd/sync"
	"github.com/spf13/cobra"
)

type SyncCmd struct {
	Command *cobra.Command
}

func NewSyncCmd(AppConfigService AppConfigService,
	ConfigsListService ConfigsListService,
	OsService OsService,
	SyncService sync_cmd.SyncService,
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
