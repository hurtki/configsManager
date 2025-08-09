# ConfigsManager
ConfigsManager is a simple and efficient CLI application that helps you 
manage your configuration files by associating keys with file paths.
**Supports only UNIX-like systems.**

Built with [Cobra CLI](https://github.com/spf13/cobra) library to provide a powerful and user-friendly command line interface.

### With configsManager, you can:
- Manage configs 
- Retrieve config paths by keys
- View content of config files
- List all stored config keys
- Open config files in your preferred editor

# Fast start 
### Linux (x86_64 / AMD64)
```
curl -L https://github.com/hurtki/configsManager/releases/latest/download/cm-linux-amd64 -o cm
chmod +x cm
sudo mv cm /usr/local/bin/
cm init
```

### macOS Intel 
```
curl -L https://github.com/hurtki/configsManager/releases/latest/download/cm-darwin-amd64 -o cm
chmod +x cm
sudo mv cm /usr/local/bin/
cm init
```

### Linux ARM64
```
curl -L https://github.com/hurtki/configsManager/releases/latest/download/cm-linux-arm64 -o cm
chmod +x cm
sudo mv cm /usr/local/bin/
cm init
```

### macOS Apple Silicon(ARM)
```
curl -L https://github.com/hurtki/configsManager/releases/latest/download/cm-darwin-arm64 -o cm
chmod +x cm
sudo mv cm /usr/local/bin/
cm init
```

# Initializing 
`cm init`

### Available Commands:
  ### add
  Add a new configuration key with its associated file path
  ```sh
  cm add [config name] [config absolute/realive path]
  ```
  Auto creating a name for config
  ```sh
  cm add [config absolute/realive path]
  ```
  Add a new config using stdIN
  ```sh
  realpath [path] | cm add [config name] #or nothing if you want to auto create a name for config
  # tool waits from you to give a path using pipe, not config name
  ```
  > **If tool sees that given path doesn't exist, it will automatically create it, and will notify you where it was created!**
  ---
  ### rm
  Remove a config from list by key
  ```
  cm rm [config key]
  ```
  > **Returns nothing if operation was successful**
  ---
  ### cat
  Print the content of the config file for a given key
  ```
  cm cat [config name]
  ```
  ---
  ### keys
  List all configuration keys
  ```
  cm keys
  ```
  --- 
  ### open
  Open a config in editor
  ```
  cm open [config name]
  ```
  --- 
  ### path
  Retrieve the file path associated with a configuration key
  ```
  cm path  [config name]
  ```
  --- 
  ### help
  Help about any command
  ```
  cm help
  ```

## 🔧 Application Config

Yes, this tool uses a config file!  
You can open it anytime with:

```bash
cm open cm_config
````

The tool folder is automatically created:
`~/.config/configsManager/`

Example structure:

```json
{
  "editor": "vim",
  "if_key_exists": "default",
}
```

---

### 🖊️ `editor`

**Default:** `"vim"`
This is the command used to open your configuration files.

Example:

```bash
[your_editor_command] /path/to/your/config.cfg
```

---

### ⚠️ `if_key_exists`

**Default:** `default`

* `default` - tool will ask you what to do + will notice to change this setting
* `o` - tool will overwrite the existing key
* `n` - tool will automatically create a new name
* `ask` - tool will always ask you what to do

---



