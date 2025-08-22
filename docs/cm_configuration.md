# ConfigsManager Docs

* [Main Page](index.html)
* [Install](installation.html)
* [Commands](commands.html)
* **[Configuration](cm_configuration.html)**

## Configuration

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

