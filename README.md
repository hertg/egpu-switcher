# egpu-switcher

> **Disclaimer**\
> Works only with NVIDIA cards and after the nvidia-drivers were installed.\
> Uses the `nvidia-xconfig --query-gpu-info` command to detect GPUs for setup and startup.\
> Not properly tested on different setups, **use at your own risk.**

This script will try to detect your GPUs and prompt you to answer which is the internal and which the external one. It will then create a `xorg.conf.egpu` and a `xorg.conf.internal` file in your `/etc/X11` directory, with each of the GPUs defined via their BusID. 

# Prerequisites
1. Install the proprietary NVIDIA drivers

# Usage
## Setup
`egpu-switcher setup`

## Switch
`egpu-switcher switch`\
This command will automatically detect if the egpu is connected and update the `xorg.conf` symlink accordingly.

`egpu-switcher switch egpu`\
This command will point the `xorg.conf` symlink to `xorg.conf.egpu`

`egpu-switcher switch internal`\
This command will point the `xorg.conf` symlink to `xorg.conf.internal`

