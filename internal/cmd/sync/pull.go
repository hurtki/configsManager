package sync_cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

type SyncPullCmd struct {
	syncService SyncService
	osService   OsService
	Command     *cobra.Command
	All         bool
	SamePlace   bool
}

func (c *SyncPullCmd) run(cmd *cobra.Command, args []string) error {
	switch len(args) {
	case 0:
		if !c.All || !c.SamePlace {
			return ErrPullBothFlagsRequired
		}
		res, err := c.syncService.PullAll()
		if err != nil {
			return err
		}
		fmt.Printf("Loading %d configs\n", res.ExpectedConfigsCount)
		i := 1
		for res := range res.Configs {
			if res.Error != nil {
				if res.ConfigObj.KeyName == "" {
					fmt.Printf("%d:error: %s\n", i, res.Error.Error())
				} else {
					fmt.Printf("%d:for '%s', error: %s\n", i, res.ConfigObj.KeyName, res.Error.Error())
				}
			} else {
				homeDir, _ := c.osService.GetHomeDir()
				cfgPath := res.ConfigObj.DeterminedPath.BuildPath(homeDir)
				if err := c.osService.MakePathAndFile(cfgPath); err != nil {
					return fmt.Errorf("%d:can't make path on config's determined path: %s, %w", i, cfgPath, err)
				}
				if err := c.osService.WriteFile(cfgPath, res.ConfigObj.Content); err != nil {
					return fmt.Errorf("%d:can't write confi's data on its determined path: %w", i, err)
				}
				fmt.Printf("%d:pulled %s to: %s\n", i, res.ConfigObj.KeyName, cfgPath)
			}

			i++
		}
	case 1:
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
	case 2:
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
	default:
		return ErrPullMoreThanTwoArgumentsProvided
	}
	return nil
}

func NewSyncPullCmd(syncService SyncService, osService OsService) *SyncPullCmd {
	syncPullCmd := &SyncPullCmd{syncService: syncService, osService: osService}

	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pulls configs from your cloud",
		Long: `Pull your configuration files from the cloud.
You can pull a single config, all configs, to a specific folder,
or restore them to their original paths.`,
		RunE: syncPullCmd.run,
	}

	cmd.Flags().BoolVar(&syncPullCmd.All, "all", false, "Pull all the configs")
	cmd.Flags().BoolVar(&syncPullCmd.SamePlace, "sp", false, "Pull selected config/s")

	syncPullCmd.Command = cmd

	return syncPullCmd
}
