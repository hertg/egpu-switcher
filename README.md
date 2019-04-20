# egpu-switcher

> **Disclaimer**\
> Works only with NVIDIA cards and after the nvidia-drivers were installed.\
> Not properly tested on different setups, **use at your own risk.**\
> Latest version is only tested with Ubuntu 19.04.

This package allows you to use your egpu and the displays directly connected to it. **There is no plug-and-play functionality**. You have to restart your computer in order to connect / disconnect your egpu.

The main goal of this package is to make the initial setup easier and allowing your computer to decide **on startup** if your egpu is connected and can be used.

## More information
This script will try to detect your GPUs and prompt you to answer which is the internal and which the external one. It will then create a `xorg.conf.egpu` and a `xorg.conf.internal` file in your `/etc/X11` directory, with each of the GPUs defined via their BusID. 

Additionally a custom `systemd` service with the following content will be added.

*/etc/systemd/system/egpu.service*
```bash
[Unit]
Description=EGPU Service

[Service]
ExecStart=egpu-switcher switch auto

[Install]
WantedBy=multi-user.target
```

This will enable the automatic detection wheter your egpu is connected or not on startup.

# Prerequisites
**Only tested on Ubuntu 19.04**
> I am using this script with my notebook, which already has a dedicated gpu. I've experienced problems when working in "Hybrid Graphics" mode, meaning that both will be activated, the integrated graphics *and* the dedicated graphics. When booting from the Ubuntu Live-CD the installer freezed multiple times. After changing the BIOS display settings from "Hybrid Graphics" to "Dedicated Graphics", the problem was gone. After successful installation of Ubuntu, the BIOS setting can be set back to "Hybrid Graphics".\
\
> I am using a Lenovo ThinkPad, your device will probably have different BIOS settings.

Install the latest Ubuntu 19.04. You can choose the minimal installation, but make sure to check the box "Install third-party software for graphics and Wi-Fi hardware".

With Ubuntu 19.04, the latest proprietary nvidia-drivers will already be installed and you won't have to tinker around with manually disabling the nouveau drivers, etc.

**Just make sure that you have the proprietary nvidia drivers installed**

# Install
```bash
$ sudo add-apt-repository ppa:hertg/egpu-switcher
$ sudo apt update
$ sudo apt install egpu-switcher
```

Just start the setup with the following command, after the setup is finished, you are good to go.

```bash
$ sudo egpu-switcher setup
```

# Uninstall
Run the following command to uninstall the package. All created files will be removed and your previous `xorg.conf` will be restored if you had one.
```bash
$ apt remove --purge egpu-switcher
```

# Commands
## Setup
`egpu-switcher setup <method>`\
Will start the setup process. 
> If no method is passed, the `lspci` will be used by default.\
> The following methods are available: 
> - `lspci` (recommended)
> - `nvidia-xconfig`

## Switch
`egpu-switcher switch auto <method>`\
This command will automatically detect if the egpu is connected and update the `xorg.conf` symlink accordingly.
> If no method is passed, the `lspci` will be used by default.\
> The following methods are available: 
> - `lspci` (recommended)
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
1. `sudo apt install devscripts`
1. `sudo apt install debhelper`
1. Update changelog: `dch`
1. Build: `debuild -S | tee /tmp/debuild.log 2>&1`
1. Upload to ppa: `dput ppa:hertg/egpu-switcher egpu-switcher_X.X.X_source.changes`
