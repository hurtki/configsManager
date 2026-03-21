package sync_cmd

import (
	"fmt"
	"path/filepath"

	"github.com/hurtki/configsManager/internal/domain"
	sync_services "github.com/hurtki/configsManager/internal/services/sync"

	"github.com/spf13/cobra"
)

type SyncPushCmd struct {
	syncService        SyncService
	configsListService ConfigsListService
	osService          OsService
	Command            *cobra.Command
	// not realised feature
	// Force             bool
}

func (c *SyncPushCmd) run(cmd *cobra.Command, args []string) error {

	configsList, err := c.configsListService.Load()
	if err != nil {
		return err
	}
	// building sync_services.ConfigObj from every config from local configList
	configObjs := []*domain.ConfigObj{}
	homeDir, _ := c.osService.GetHomeDir()
	for _, key := range configsList.GetAllKeys() {
		cfgObj := domain.ConfigObj{}
		// key
		cfgObj.KeyName = key
		// path
		absCfgPath, _ := configsList.GetPath(key)
		// buuilding determind path
		cfgObj.DeterminedPath = domain.NewDeterminedPath(absCfgPath, homeDir)
		// getting filename
		cfgObj.FileName = filepath.Base(absCfgPath)
		data, err := c.osService.GetFileData(absCfgPath)

		if err != nil {
			fmt.Printf("Config's data with key: %s is not found on path: %s\n", key, absCfgPath)
			fmt.Println("Skipping...")
			continue
		}

		cfgObj.Content = data
		// adding to slice
		configObjs = append(configObjs, &cfgObj)
	}
	if len(configObjs) == 0 {
		return ErrNoLocalConfigsForPush
	}
	results, err := c.syncService.Push(configObjs, true)
	if err != nil {
		return err
	}
	return c.printResults(results)

}

func (c *SyncPushCmd) printResults(results []*sync_services.SyncResult) error {
	if len(results) == 0 {
		fmt.Println("✅ pushed changes successfully")
		return nil
	}

	var successes []string
	var failures []string

	for _, res := range results {
		if res.Error != nil {
			failures = append(failures, fmt.Sprintf("  - %s: %s", res.ConfigObj.KeyName, res.Error.Error()))
		} else {
			successes = append(successes, fmt.Sprintf("  - %s", res.ConfigObj.KeyName))
		}
	}

	if len(successes) > 0 {
		fmt.Println("Successfully pushed:")
		for _, s := range successes {
			fmt.Println(s)
		}
	}

	if len(failures) > 0 {
		fmt.Println("Failed to push:")
		for _, f := range failures {
			fmt.Println(f)
		}
	}

	return nil
}

func NewSyncPushCmd(syncService SyncService,
	configListService ConfigsListService,
	osService OsService,
) *SyncPushCmd {
	syncPushCmd := SyncPushCmd{
		syncService:        syncService,
		configsListService: configListService,
		osService:          osService,
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
