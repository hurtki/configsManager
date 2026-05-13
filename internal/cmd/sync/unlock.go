package sync_cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type SyncUnlockCmd struct {
	syncService SyncService
	Command     *cobra.Command
}

func (c *SyncUnlockCmd) run(cmd *cobra.Command, args []string) error {
	encoded, err := c.syncService.Unlock()
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "Session key derived. Apply with: eval \"$(cm sync unlock)\"")
	fmt.Printf("export CM_SYNC_KEY=%s\n", encoded)
	return nil
}

func NewSyncUnlockCmd(syncService SyncService) *SyncUnlockCmd {
	syncUnlockCmd := SyncUnlockCmd{
		syncService: syncService,
	}

	cmd := &cobra.Command{
		Use:   "unlock",
		Short: "Derive a per-session key and print an `export CM_SYNC_KEY=...` line",
		Long: `Prompts for your passphrase, derives the AES key once via Argon2id and
prints a shell-evaluable line to stdout that sets CM_SYNC_KEY in your
environment. Subsequent cm commands in the same shell will use that key
directly, skipping the passphrase prompt and the slow KDF step.

Usage:

    eval "$(cm sync unlock)"

Then run cm sync push/pull/etc. without re-entering the passphrase.

SECURITY: CM_SYNC_KEY grants the same access to your tokens as the
passphrase would. Do not export it to other processes, log it, or copy
it between machines. It dies with the shell — that is intentional.`,
		RunE: syncUnlockCmd.run,
	}

	syncUnlockCmd.Command = cmd
	return &syncUnlockCmd
}
