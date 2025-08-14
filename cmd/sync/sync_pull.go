package sync_cmd

import (
	"fmt"

	sync_services "github.com/hurtki/configsManager/services/sync"
	"github.com/spf13/cobra"
)

type SyncPullCmd struct {
	SyncDeps  *sync_services.Deps
	Command   *cobra.Command
	All       bool
	SamePlace bool
}

func (c *SyncPullCmd) run(cmd *cobra.Command, args []string) error {
	AllFlag := c.All
	SpFlag := c.SamePlace

	if len(args) == 0 {
		if !(AllFlag && SpFlag) {
			return ErrPullBothFlagsRequired
		}
		// здесь нужно будет получить от SyncService все конфиги тоесть []ConfigObj
		// дальше раскинуть их по папками с паралельным сохранением пути в локальный ConfigsList
	} else if len(args) == 1 {
		if AllFlag {
			return ErrPullAllFlagNotSupported
		}
		// здесь надо получить от SyncService один ConfigObj по ключу
		// дальше мы закидываем его в папку где мы есть либо если --sp то где он должен быть
		// добавляем в локальным ConfigsList
	} else if len(args) == 2 {
		if AllFlag || SpFlag {
			return ErrPullAllAndSpFlagsNotSupported
		}
		// здесь надо получить от SyncService ConfigObj по ключу
		// и четко положить его в папку куда сказал пользователь
	} else {
		return ErrPullMoreThanTwoArgumentsProvided
	}

	fmt.Printf("Flag all: %t\n", AllFlag)
	fmt.Printf("Flag sp: %t\n", SpFlag)
	fmt.Println(args)
	return nil
}

func NewSyncPullCmd(d *sync_services.Deps) *SyncPullCmd {
	syncPullCmd := SyncPullCmd{
		SyncDeps: d,
	}

	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pulls cofnigs from your cloud",
		Long:  ``,
		RunE:  syncPullCmd.run,
	}

	cmd.Flags().BoolVar(&syncPullCmd.All, "all", false, "Pull all the configs")
	cmd.Flags().BoolVar(&syncPullCmd.SamePlace, "sp", false, "Pull selected config/s")

	syncPullCmd.Command = cmd

	return &syncPullCmd
}
