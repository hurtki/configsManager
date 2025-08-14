package sync_cmd

import (
	"fmt"

	sync_services "github.com/hurtki/configsManager/services/sync"
	"github.com/spf13/cobra"
)

type SyncPushCmd struct {
	SyncDeps *sync_services.Deps
	Command  *cobra.Command
	Force    bool
}

func (c *SyncPushCmd) run(cmd *cobra.Command, args []string) error {
	ForceFlag := c.Force

	fmt.Printf("Flag force: %t\n", ForceFlag)

	// здесь нужно думаю сгенерировать из ConfigsList ConfigObj-ты и передать их в SyncService
	// вот так вот передать Push(configs []*ConfigObj) map[*ConfigObj]error
	// еще можно туда как раз флаг --force передавать чтобы если были противоречия они либо возвращали соответствующую ошибку либо делали overwrite
	// дальше нам вернули map[*ConfigObj]error и мы просто показываем всем ошибки, красиво выводим
	return nil
}

func NewSyncPushCmd(d *sync_services.Deps) *SyncPushCmd {
	syncPushCmd := SyncPushCmd{
		SyncDeps: d,
	}

	cmd := &cobra.Command{
		Use:   "push",
		Short: "Push saved configs to cloud",
		Long:  ``,
		RunE:  syncPushCmd.run,
	}

	cmd.Flags().BoolVar(&syncPushCmd.Force, "force", false, "Ignore inappropriate configs while pushing")

	syncPushCmd.Command = cmd

	return &syncPushCmd
}
