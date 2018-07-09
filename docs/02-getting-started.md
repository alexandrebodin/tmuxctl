---
title: Getting started
---

Tmuxctl uses config files to start a session. 

Start by creating a toml config file (by default `tmuxctl` will look for a `.tmuxctlrc` in the current dir, up the parent directories)

```toml
name="azdaz"

[[windows]]
  name="win-1"
```

You can see more examples in the [examples](https://github.com/alexandrebodin/tmuxctl/tree/master/__examples__) folder