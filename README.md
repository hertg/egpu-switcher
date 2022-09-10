# egpu-switcher

> **Note**: This is the `legacy` branch for the `bash` version of egpu-switcher.

Distribution agnostic eGPU script that works with **NVIDIA** and **AMD** cards.

---

- [Description](#description)
- [Screenshot](#screenshot)
- [Requirements](#requirements)
- [Installation](#installation)
  - [Ubuntu (apt)](#ubuntu-apt)
  - [Other](#other)
- [Commands](#commands)
- [Troubleshooting](#troubleshooting)
- [Background Information](#background-information)

---

## Description
The goal of this script is to lower the barrier for Linux users to use their eGPU on the Linux Desktop.
An interactive setup allows the user to choose their external GPU, which will then be automatically chosen as the primary GPU if it's connected on bootup.

**This does not provide you with a plug-and-play functionality like you may know from Windows.<br> You still need to reboot your computer in order to connect / disconnect your eGPU.**.

## Screenshot
![Screenshot of setup](https://raw.githubusercontent.com/hertg/egpu-switcher/master/images/screenshot_setup.png)

## Requirements
1. Your OS is running X-Server.
1. You have at least **pciutils 3.3.x or higher** installed (check with `lspci --version`).
1. You have at least **Bash 4.x or higher** installed.
1. You have **already authorized your Thunderbolt EGPU** and are able to connect.
1. You have already **installed the latest (proprietary) drivers for your GPUs**.

---

## Installation

### Ubuntu (apt)
Installation and setup:
```bash
$ sudo add-apt-repository ppa:hertg/egpu-switcher
$ sudo apt update
$ sudo apt install egpu-switcher
$ sudo egpu-switcher setup
```

Uninstall:
```bash
$ apt remove egpu-switcher
```

### Other
Installation and setup:
```bash
$ git clone git@github.com:hertg/egpu-switcher.git
$ cd egpu-switcher
$ make install
$ sudo egpu-switcher setup
```

Uninstall: 
> **Critical Warning**: **Do not use this command on any version prior to `0.10.2`!**\
> There was a critical typo in the Makefile which would delete your `/usr/bin` folder. Please do a manual uninstall by removing the `egpu-switcher` folder in the `/usr/bin/` and the `/usr/share/` directory.
```bash
$ sudo egpu-switcher cleanup
$ make uninstall
```

---

## Commands
<pre>
<b>egpu-switcher setup</b> [--override] [--noprompt]
    This will generate the "xorg.conf.egpu" and "xorg.conf.internal" files and
    symlink the "xorg.conf" file to one of them.
    
    It will also create the systemd service, that runs the "switch" command on each bootup.
    
    This will NOT delete any already existing files. If an "xorg.conf" file already exists, 
    it will be backed up to "xorg.conf.backup.{datetime}". This can later be reverted by
    executing the "cleanup" command.

    <b>--override</b>
        If an AMD GPU or open-source NVIDIA drivers are used, the "switch" command 
        will prevent from switching to the eGPU if there are no displays directly attached to it. 
        This flag will make sure to switch to the EGPU even if there are no displays attached.

    <b>--noprompt</b>
        Prevent the setup from prompting for user interaction if there is
        no existing configuration file found. (Is currently only used by the "postinst" script)
</pre>

<pre>
<b>egpu-switcher switch auto|egpu|internal</b> [--override]
    Switches to the specified GPU. if the <u>auto</u> parameter is used, the script will check if 
    the eGPU is attached and switch accordingly. 
    
    The computer (or display-manager) needs to be restarted for this to take effect.

    <b>--override</b>
        If an AMD GPU or open-source NVIDIA drivers are used, the "switch" command 
        will prevent from switching to the eGPU if there are no displays directly attached to it. 
        This flag will make sure to switch to the EGPU even if there are no displays attached.
</pre>

<pre>
<b>egpu-switcher cleanup</b> [--hard]
    Remove all files egpu-switcher has created previously and restore the backup
    of previous "xorg.conf" files.

    <b>--hard</b>
        Remove configuration files too.
</pre>

<pre>
<b>egpu-switcher config</b>
    Prompts the user to specify their external/internal GPU and saves their answer
    to the configuration file.
</pre>

<pre>
<b>egpu-switcher remove</b>
    Allows the user to remove their eGPU without a complete reboot.
    This method will still restart the display-manager, and therefore terminate all its child-processes.
</pre>

---

## Troubleshooting
If you run into problems, please have a look at [TROUBLESHOOT.md](https://github.com/hertg/egpu-switcher/blob/master/TROUBLESHOOT.md) before reporting any issues.

---

## Background information
> A backup of your current `xorg.conf` will be created, nothing gets deleted. If the script doesn't work for you, you can revert the changes by executing `egpu-switcher cleanup` or just completely uninstall the script with `apt remove egpu-switcher`. This will remove all files it has created and also restore your previous `xorg.conf` file.

This script will create two configuration files in your X-Server folder `/etc/X11`.
The file `xorg.conf.egpu` holds the settings for your EGPU and the file `xorg.conf.internal` holds the settings for your internal graphics.

Then a symlink `xorg.conf` will be generated which points to the corresponding config file, depending on wheter your egpu is connected or not.

Additionally a custom `systemd` service with the following content will be created.

*/etc/systemd/system/egpu.service*
```bash
[Unit]
Description=EGPU Service
Before=display-manager.service
After=bolt.service

[Service]
Type=oneshot
ExecStart=/usr/bin/egpu-switcher switch auto

[Install]
WantedBy=graphical.target
```

This will enable the automatic detection wheter your egpu is connected or not on startup.
