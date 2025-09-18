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

### Fast Start

Log into Dropbox:

```sh
cm sync auth --dropbox
```

> You will receive a link from the tool. Follow the OAuth2 workflow by opening the link and entering the code provided.

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
