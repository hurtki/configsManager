# ConfigsManager
ConfigsManager is a simple and efficient CLI application that helps you 
manage your configuration files by associating keys with file paths, then easly sync them with cloud.
**Supports only UNIX-like systems.**

Built with [Cobra CLI](https://github.com/spf13/cobra) library to provide a powerful and user-friendly command line interface.

### With configsManager, you can:
- Manage configs 
- **Sync configs with cloud**
- Retrieve config paths by keys
- View content of config files
- List all stored config keys
- Open config files in your preferred editor


### [Full documentaion](https://hurtki.github.io/configsManager/)

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
