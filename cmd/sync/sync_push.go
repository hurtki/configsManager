package sync_cmd

import (
	"fmt"
	"path/filepath"

	"github.com/hurtki/configsManager/services"
	"github.com/hurtki/configsManager/services/sync"
	"github.com/spf13/cobra"
)

type SyncPushCmd struct {
	syncService       sync_services.SyncService
	configListService services.ConfigsListService
	osService         services.OsService
	Command           *cobra.Command
	Force             bool
}

func (c *SyncPushCmd) run(cmd *cobra.Command, args []string) error {
	ForceFlag := c.Force

	fmt.Printf("Flag force: %t\n", ForceFlag)

	configList, err := c.configListService.Load()
	if err != nil {
		return err
	}
	configObjs := []*sync_services.ConfigObj{}
	for _, key := range configList.GetAllKeys() {
		cfgObj := sync_services.ConfigObj{}
		cfgObj.KeyName = key
		cfgObj.DeterminedPath, _ = configList.GetPath(key)
		cfgObj.FileName = filepath.Base(cfgObj.DeterminedPath)
		data, err := c.osService.GetFileData(cfgObj.DeterminedPath)

		if err != nil {
			return err
		}
		cfgObj.Content = data
		configObjs = append(configObjs, &cfgObj)
	}

	results := c.syncService.Push(configObjs, true)
	for _, res := range results {
		if res.Error != nil {
			fmt.Printf("error pushing cfg: , error: %s\n", res.Error.Error())
		} else {
			fmt.Printf("pushed cfg: %s successfully\n", res.ConfigObj.KeyName)
		}
	}

	// здесь нужно думаю сгенерировать из ConfigsList ConfigObj-ты и передать их в SyncService
	// вот так вот передать Push(configs []*ConfigObj) map[*ConfigObj]error
	// еще можно туда как раз флаг --force передавать чтобы если были противоречия они либо возвращали соответствующую ошибку либо делали overwrite
	// дальше нам вернули map[*ConfigObj]error и мы просто показываем всем ошибки, красиво выводим
	return nil
}

func NewSyncPushCmd(syncService sync_services.SyncService,
	configListService services.ConfigsListService,
	osService services.OsService,
) *SyncPushCmd {
	syncPushCmd := SyncPushCmd{
		syncService:       syncService,
		configListService: configListService,
		osService:         osService,
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
