package sync_cmd

import (
	"fmt"
	"path/filepath"

	"github.com/hurtki/configsManager/services"
	"github.com/hurtki/configsManager/services/sync"
	"github.com/spf13/cobra"
)

type SyncPullCmd struct {
	syncService sync_services.SyncService
	osService   services.OsService
	Command     *cobra.Command
	All         bool
	SamePlace   bool
}

func (c *SyncPullCmd) run(cmd *cobra.Command, args []string) error {

	if len(args) == 0 {
		if !c.All || !c.SamePlace {
			return ErrPullBothFlagsRequired
		}
		res, err := c.syncService.PullAll()
		if err != nil {
			return err
		}

		for i, res := range res {

			if res.Error != nil {
				if res.ConfigObj.KeyName == "" {
					fmt.Printf("error for index: %s, error: %d\n", fmt.Sprint(i), res.Error)
				} else {
					fmt.Printf("for '%s', error: %d\n", res.ConfigObj.KeyName, res.Error)
				}

			} else {
				homeDir, _ := c.osService.GetHomeDir()
				cfgPath := res.ConfigObj.DeterminedPath.BuildPath(homeDir)
				if err := c.osService.MakePathAndFile(cfgPath); err != nil {
					return err
				}
				if err := c.osService.WriteFile(cfgPath, res.ConfigObj.Content); err != nil {
					return err
				}
				fmt.Printf("pulled config to: %s\n", cfgPath)
			}
		}

	} else if len(args) == 1 {
		if c.All {
			return ErrPullAllFlagNotSupported
		}
		res := c.syncService.PullOne(args[0])
		if res.Error != nil {
			return res.Error
		}
		if c.SamePlace {
			homeDir, _ := c.osService.GetHomeDir()
			cfgPath := res.ConfigObj.DeterminedPath.BuildPath(homeDir)
			if err := c.osService.MakePathAndFile(cfgPath); err != nil {
				return err
			}
			if err := c.osService.WriteFile(cfgPath, res.ConfigObj.Content); err != nil {
				return err
			}
			fmt.Printf("pulled config to: %s\n", cfgPath)
		} else {
			if err := c.osService.WriteFile(res.ConfigObj.FileName, res.ConfigObj.Content); err != nil {
				return err
			}
			fmt.Printf("pulled config: %s to executing folder\n", res.ConfigObj.FileName)
		}

	} else if len(args) == 2 {
		if c.SamePlace || c.All {
			return ErrPullAllAndSpFlagsNotSupported
		}
		res := c.syncService.PullOne(args[0])
		if res.Error != nil {
			return res.Error
		}
		path := filepath.Join(args[1], res.ConfigObj.FileName)
		if err := c.osService.MakePathAndFile(path); err != nil {
			return err
		}
		if err := c.osService.WriteFile(path, res.ConfigObj.Content); err != nil {
			return err
		}
		fmt.Printf("pulled config to: %s\n", path)
	} else {
		return ErrPullMoreThanTwoArgumentsProvided
	}
	return nil
}

func NewSyncPullCmd(syncService sync_services.SyncService, osService services.OsService) *SyncPullCmd {
	syncPullCmd := &SyncPullCmd{syncService: syncService, osService: osService}

	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pulls configs from your cloud",
		Long:  ``,
		RunE:  syncPullCmd.run,
	}

	cmd.Flags().BoolVar(&syncPullCmd.All, "all", false, "Pull all the configs")
	cmd.Flags().BoolVar(&syncPullCmd.SamePlace, "sp", false, "Pull selected config/s")

	syncPullCmd.Command = cmd

	return syncPullCmd
}
