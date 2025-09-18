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
	// not realised feature
	// Force             bool
}

func (c *SyncPushCmd) run(cmd *cobra.Command, args []string) error {

	configList, err := c.configListService.Load()
	if err != nil {
		return err
	}
	configObjs := []*sync_services.ConfigObj{}
	homeDir, _ := c.osService.GetHomeDir()
	for _, key := range configList.GetAllKeys() {
		cfgObj := sync_services.ConfigObj{}
		cfgObj.KeyName = key
		absCfgPath, _ := configList.GetPath(key)
		cfgObj.DeterminedPath = sync_services.NewDeterminedPath(absCfgPath, homeDir)
		cfgObj.FileName = filepath.Base(absCfgPath)
		data, err := c.osService.GetFileData(absCfgPath)

		if err != nil {
			return err
		}
		cfgObj.Content = data
		configObjs = append(configObjs, &cfgObj)
	}

	results, err := c.syncService.Push(configObjs, true)
	if err != nil {
		return err
	}
	for _, res := range results {
		if res.Error != nil {
			fmt.Printf("error pushing %s, error: %s\n", res.ConfigObj.KeyName, res.Error.Error())
		} else {
			fmt.Printf("pushed %s successfully\n", res.ConfigObj.KeyName)
		}
	}
	if len(results) == 0 {
		fmt.Println("pushed changes successfully")
	}

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
		Short: "Push saves configs to cloud",
		Long: `Push your local configuration files to the cloud.
This command uploads the files, their metadata, and the paths
where they were stored at the time of pushing.

You can use it to back up your configs quickly and safely,
so that later you can pull them to any machine or folder.`,
		RunE: syncPushCmd.run,
	}

	// not realised feature
	//cmd.Flags().BoolVar(&syncPushCmd.Force, "force", false, "Ignore inappropriate configs while pushing")

	syncPushCmd.Command = cmd

	return &syncPushCmd
}
