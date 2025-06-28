# ConfigsManager
ConfigsManager is a simple and efficient CLI application that helps you 
manage your configuration files by associating keys with file paths.

Built with [Cobra CLI](https://github.com/spf13/cobra) library to provide a powerful and user-friendly command line interface.

### With configsManager, you can:
- Add new config entries with keys and paths
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
```

### macOS Intel 
```
curl -L https://github.com/hurtki/configsManager/releases/latest/download/cm-darwin-amd64 -o cm
chmod +x cm
sudo mv cm /usr/local/bin/
```

### Linux ARM64
```
curl -L https://github.com/hurtki/configsManager/releases/latest/download/cm-linux-arm64 -o cm
chmod +x cm
sudo mv cm /usr/local/bin/
```

### macOS Apple Silicon(ARM)
```
curl -L https://github.com/hurtki/configsManager/releases/latest/download/cm-darwin-arm64 -o cm
chmod +x cm
sudo mv cm /usr/local/bin/
```
### Windows AMD64 (64-bit)
Download latest release 
https://github.com/hurtki/configsManager/releases/latest/download/cm-windows-amd64.exe

Rename to cm.exe
Place it in a folder that is in the PATH environment variable, for example:
`C:\Users\YourUser\AppData\Local\Microsoft\WindowsApp`

# Usage
`cm [command]`

### Available Commands:
  Add a new configuration key with its associated file path
  ```
  cm add [config name] [config path]
  ```  
  Print the content of the config file for a given key
  ```
  cm cat [config name]
  ```  
  List all configuration keys
  ```
  cm keys
  ```  
  Open a config in editor
  ```
  cm open [config name]
  ```
  Retrieve the file path associated with a configuration key
  ```
  cm path  [config name]
  ```
  Help about any command
  ```
  cm help
  ```
# Application config 
Yeah we have a config, you need to specify here your editor
You can open it with
```
cm open cm_config
```
It creates automatically, with this structure in `user/.config/configManager.json`:
```
{
  "editor": "vim"
}
```
Default editor "vim"
It's just a command that runs before the config path 
```
[your_editor_command] /path/to/your/config.cfg
```


