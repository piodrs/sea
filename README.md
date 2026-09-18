# Sea

A experimental small terminal text editor written in Go

## Keybindings

### Movement

| Key | Action |
| --- | --- |
| `Ctrl+F` / `Ctrl+B` | Move forward / backward one character |
| `Ctrl+N` / `Ctrl+P` | Move down / up one line |
| `Ctrl+A` / `Ctrl+E` | Move to the beginning / end of the line |
| `Ctrl+V` / `Alt+V` | Move down / up one screen |
| `Alt+<` / `Alt+>` | Move to the beginning / end of the file |

### Editing

| Key | Action |
| --- | --- |
| `Ctrl+D` / `Backspace` | Delete forward / backward |
| `Ctrl+K` | Cut to the end of the line (cuts the newline when already at the end) |
| `Ctrl+Y` | Paste the last cut |

### File & control

| Key | Action |
| --- | --- |
| `Ctrl+X Ctrl+S` | Save |
| `Ctrl+X Ctrl+C` | Quit (repeat the sequence to discard unsaved changes) |
| `Ctrl+G` / `Escape` | Cancel a pending command |