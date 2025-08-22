# ConfigsManager Docs

* [Main Page](index.html)
* [Install](installation.html)
* **[Commands](commands.html)**
* [Configuration](cm_configuration.html)

## Commands:

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
**If tool sees that given path doesn't exist, it will automatically create it and will notify you where it was created!**

---

### rm
Remove a config/s from list by key/s
```
cm rm key1 key2 key3...
```
**Returns nothing if operation was successful and if key didn't exist**

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