# egpu-switcher

> **Disclaimer**\
> Works only with NVIDIA cards and after the nvidia-drivers were installed.\
> Uses the `nvidia-xconfig --query-gpu-info` command to detect GPUs for setup and startup.\
> Not properly tested on different setups, **use at your own risk.**

This script will try to detect your GPUs and prompt you to answer which is the internal and which the external one. It will then create a `xorg.conf.egpu` and a `xorg.conf.internal` file in your `/etc/X11` directory, with each of the GPUs defined via their BusID. 

Additionally a custom `init.d` script with the following content will be created.

*/etc/init.d/egpu-switcher*
```bash
egpu-switcher switch auto
```

A symlink to this script is created in `/etc/rc5.d/egpu-switcher`.
This will enable the automatic detection wheter your egpu is connected or not on startup.

# Prerequisites
1. Install the proprietary NVIDIA drivers
1. Make sure that the `nvidia-xconfig` is installed

# Install
```bash
$ sudo add-apt-repository ppa:hertg/egpu-switcher
$ sudo apt update
$ sudo apt install egpu-switcher
```

# Usage
## Setup
`egpu-switcher setup <method>`\
Will start the setup process. 
> If no method is passed, the `lspci` will be used by default.\
> The following methods are available: 
> - `lspci`
> - `nvidia-xconfig`

## Switch
`egpu-switcher switch auto <method>`\
This command will automatically detect if the egpu is connected and update the `xorg.conf` symlink accordingly.
> If no method is passed, the `nvidia-xconfig` will be used by default.\
> The following methods are available: 
> - `nvidia-xconfig`

`egpu-switcher switch egpu`\
This command will point the `xorg.conf` symlink to `xorg.conf.egpu`

`egpu-switcher switch internal`\
This command will point the `xorg.conf` symlink to `xorg.conf.internal`

## Cleanup
`egpu-switcher cleanup`\
This command will revert the whole `setup` process and remove all files it has created.
Additionally the command will restore your previous `xorg.conf` that you had before running the `setup`.

> This command is executed automatically while doing a `apt remove --purge egpu-switcher`.

# Build
1. Update changelog: `dch`
1. Build: `debuild -S | tee /tmp/debuild.log 2>&1`
1. Upload to ppa: `dput ppa:hertg/egpu-switcher egpu-switcher_X.X.X_source.changes`
