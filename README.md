# egpu-switcher

> **Disclaimer**\
> Works with NVIDIA as well as AMD cards.\
> Tested on Ubuntu, but may work on any Linux (with X-Server) that supports `#!/bin/bash` scripts.
>
> For more information and user feedback, take a look at my [Thread](https://egpu.io/forums/thunderbolt-linux-setup/ubuntu-19-04-easy-to-use-setup-script-for-your-egpu/) over on egpu.io.

![Screenshot of setup](https://raw.githubusercontent.com/hertg/egpu-switcher/master/images/screenshot_setup.png)

# TL;DR

## Ubuntu (apt)
Installation and setup:
```bash
$ sudo add-apt-repository ppa:hertg/egpu-switcher
$ sudo apt update
$ sudo apt install egpu-switcher
$ sudo egpu-switcher setup
```

Uninstall:
```bash
$ apt remove --purge egpu-switcher
```

## Other
Installation and setup:
```bash
$ git clone git@github.com:hertg/egpu-switcher.git
$ cd egpu-switcher
$ make install
$ sudo egpu-switcher setup
```

Uninstall: 
> **Warning**: Do not use this command on any version prior to 0.10.2!
> There was a critical typo in the Makefile which would delete your `/usr/bin` folder. Please do a manual uninstall by removing the `egpu-switcher` folder in the `/usr/bin/` and the `/usr/share/` directory.
```bash
$ sudo egpu-switcher cleanup
$ make uninstall
```

# Goal
The goal of this script is to make the initial egpu setup for (new) Linux users less of a pain. With this script your X-Server configs for the different GPUs will be automatically created, you just have to choose which one is the external and which the internal graphics card.

After the setup, your linux installation will at each startup check if your EGPU is connected or not, and then automatically choose the right X-Server configuration in the background.

**This does not provide you with a plug-and-play functionality like you may know from Windows. If you want do connect / disconnect your EGPU, you will have to restart your computer**.

# Prerequisites
1. You have already authorized your Thunderbolt EGPU and are able to connect
1. You have already installed the latest (proprietary) drivers for your GPUs

> When installing Ubuntu 19.04, please check the box "Install third-party software for graphics and Wi-Fi hardware". After that, all required drivers will be installed automatically.

> **Hint for people with hybrid graphics**\
> I am using a Lenovo notebook with hybrid graphics (internal graphics **and** a dedicated GPU). I've experienced freezes in the Ubuntu 19.04 installer which could only be resolved by changing the display settings in the BIOS from ~~Hybrid Graphics~~ to **Discrete Graphics**. After the installation was complete, i was able to change this setting back to **Hybrid Graphics**, without any issues.

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

# Background information
> A backup of your current `xorg.conf` will be created, nothing gets deleted. If the script doesn't work for you, you can revert the changes by executing `egpu-switcher cleanup` or just completely uninstall the script with `apt remove --purge egpu-switcher`. This will purge all files it has created and also restore your previous `xorg.conf` file.

This script will create two configuration files in your X-Server folder `/etc/X11`.
The file `xorg.conf.egpu` holds the settings for your EGPU and the file `xorg.conf.internal` holds the settings for your internal graphics.

Then a symlink `xorg.conf` will be generated which points to the corresponding config file, depending on wheter your egpu is connected or not.

Additionally a custom `systemd` service with the following content will be created.

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

# Available commands
You usually don't need these commands (apart from the `setup`, of course).

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

> This command is executed automatically when uninstalling via `apt remove --purge egpu-switcher`.

# Build (notes to myself)
1. `sudo apt install devscripts`
1. `sudo apt install debhelper`
1. Update changelog: `dch`
1. Build: `debuild -S | tee /tmp/debuild.log 2>&1`
1. Upload to ppa: `dput ppa:hertg/egpu-switcher egpu-switcher_X.X.X_source.changes`
