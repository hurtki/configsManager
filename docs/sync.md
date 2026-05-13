# ConfigsManager Docs

* [Main Page](index.html)
* [Install](installation.html)
* [Commands](commands.html)
* **[Sync](sync.html)**
* [Configuration](cm_configuration.html)

---

## Sync Feature

The **ConfigsManager Sync** feature allows you to quickly save and restore all your configuration files in the cloud.

### Key Features

* **Push** – Upload your configs to the cloud.
* **Pull** – Download your configs from the cloud.

---

### How tokens are stored

ConfigsManager keeps your cloud OAuth tokens in a single encrypted file at
`~/.config/configsManager/tokens.dat` (permissions `0600`).

* Encryption: **AES-256-GCM** (authenticated).
* Key derivation: **Argon2id** with parameters `time=4`, `memory=128 MiB`,
  `threads=4`. The KDF parameters and salt are written inside the file, so
  they can be bumped in future releases without breaking existing installs.
* No system keyring is used — no KWallet, no Secret Service, no Keychain,
  no `pass`, no external daemons.

You set a passphrase the first time you run `cm sync auth`. Every later
command that needs the tokens will ask for that passphrase (or use a
session key — see [Unlock a terminal session](#unlock-a-terminal-session)
below).

---

### Fast Start

Log into Dropbox:

```sh
cm sync auth --dropbox
```

The first time you run this you will be asked to **create a passphrase**
(entered twice). This passphrase protects the encrypted token file —
remember it; there is no recovery.

> You will then receive a link from the tool. Follow the OAuth2 workflow
> by opening the link and entering the code provided.

---

### Push

Push all your configuration files and metadata to the cloud:

```sh
cm sync push
```

> Stores not only the files but also their paths and metadata at the moment of the push.

---

### Pull

Pull configs from the cloud with different options:

**1. Pull to the current folder:**

```sh
cm sync pull [config_key_in_cloud]
```

**2. Pull to a specific folder:**

```sh
cm sync pull [config_key_in_cloud] [path/]
```

> The tool creates folders automatically if they do not exist.

**3. Pull to the original path (`--sp` flag):**

```sh
cm sync pull [config_key_in_cloud] --sp
```

> `--sp` stands for **Determined Path** – the original path where the config was located when pushed.
> Example: Fish shell config is stored in `~/.config/fish/config.fish`. When pulling with `--sp`, the tool restores it to the correct location automatically.

**4. Pull all configs to their determined paths (killer feature!):**

```sh
cm sync pull --all --sp
```

> Automatically restores **all synced configs** to the proper folders.

---

### Unlock a terminal session

Argon2id is slow on purpose — it adds about a second to every command that
touches the token file, and you have to retype the passphrase each time.
For a series of commands in one terminal there is a shortcut:

```sh
eval "$(cm sync unlock)"
```

* `cm sync unlock` asks for the passphrase **once**, derives the AES key,
  and prints `export CM_SYNC_KEY=<base64>` to stdout. The `eval` puts that
  variable into your current shell.
* Every later `cm sync …` command in the same shell reads `CM_SYNC_KEY`,
  skips the passphrase prompt **and** the slow Argon2id step, and decrypts
  the file directly.
* `CM_SYNC_KEY` lives only in this shell. Closing the terminal forgets it.

To clear it without closing the terminal:

```sh
eval "$(cm sync lock)"
```

That just runs `unset CM_SYNC_KEY` — the next `cm sync …` will ask for
the passphrase again.

#### Security notes

* `CM_SYNC_KEY` grants the same access to your tokens as the passphrase
  would. Treat it like a password: do not log it, do not pass it to other
  processes, do not copy it between machines.
* Anything that can read your process environment (same UID, `ps eww`,
  `/proc/<pid>/environ` on Linux) can read the key while it is set.
* If the variable becomes stale (for example you changed the passphrase
  on another machine and pulled a new `tokens.dat`), the tool will print a
  warning and fall back to a passphrase prompt — just run
  `eval "$(cm sync unlock)"` again to refresh.

---

### Logout

The `logout` command removes saved authentication tokens from your system.

**Usage:**
Logout from all cloud services
```sh
cm sync logout
```
Logout only from Dropbox
```sh
cm sync logout --dropbox
```

> The command safely deletes access and refresh tokens stored by ConfigsManager.
---

