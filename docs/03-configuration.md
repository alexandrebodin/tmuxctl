---
title: Configuration
---

Tmuxctl works with configuration files. Here is the list of possible options

## name

You can set a session name

```toml
name="session-name"
```

## dir

You can set a base directory for your windows

```toml
dir="~/zad"
```

## select-window

You can choose which window to attach to on start

```toml
select-window="window-name"
```

## select-pane

You can choose which pane to attach to on start (the index of the pane in the window starting from 1)

```toml
select-pane=2
```

## clear-panes

You can clear the panes (C-l) after initialisation

```toml
clear-panes=true
```

## window-scripts

You can run scripts in every window (runs before anything else in a window)

```toml
window-scripts=[
  "cd folder",
  "ls -larth"
]
```

## [[windows]]

You can add windows like follow

```toml
name="session-name"

[[windows]]

[[windows]]

[[windows]]
```

## name

You can give a name to a window

```toml
name="window-name"
```

## dir

You can start a window in a directory (every pane in this window will start there by default)

```toml
dir="/some-dir"
```

## scripts

You can run scripts in the window before it is splitted in panes (they will run once)

```toml
scripts=[
  "do sth"
]
```

## pane-scripts

You can run scripts in every pane of a window if needed

```toml
pane-scripts=[
  "do sth"
]
```

## sync

You can synchronize all the panes of a window

```toml
sync=true
```

## [[windows.panes]]

You can declare all the paes you want to have in a window

```toml
name="session-name"

[[windows]]
  name="window-name"

  [[windows.panes]]

  [[windows.panes]]

  [[windows.panes]]
```

## dir

You can start a pane in a specific dir

```toml
dir="/some-dir"
```

## scripts

You can run a list of scripts in a specific pane

```toml
scripts=["do stuff"]
```

## split

You can specify the split options to setup a pane size. See [the docs](https://www.systutorials.com/docs/linux/man/1-tmux/) for split-window options

```toml
split="-f -h"
```

## zoom

You can set a pane as zoomed (only visible pane in the window) on startup

```toml
zoom=true
```
