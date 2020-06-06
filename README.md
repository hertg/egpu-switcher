# egpu-switcher

> **Disclaimer**\
> Works with **NVIDIA** as well as **AMD** cards.\
> Tested on Ubuntu, Arch and many other distros.
>
> For more information and user feedback, take a look at my [Thread](https://egpu.io/forums/thunderbolt-linux-setup/ubuntu-19-04-easy-to-use-setup-script-for-your-egpu/) over on egpu.io or open an issue on Github.

## Description
The goal of this script is to make the initial egpu setup for (new) Linux users less of a pain. With this script your X-Server configs for the different GPUs will be automatically created, you just have to choose which one is the external and which the internal graphics card.

After the setup, your linux installation will at each startup check if your EGPU is connected or not, and then automatically choose the right X-Server configuration in the background.

**This does not provide you with a plug-and-play functionality like you may know from Windows. If you want do connect / disconnect your EGPU, you will have to restart your computer**.

## Screenshot
![Screenshot of setup](https://raw.githubusercontent.com/hertg/egpu-switcher/master/images/screenshot_setup.png)

## Requirements
1. Your OS is running X-Server.
1. You have at least **pciutils 3.3.x or higher** installed (check with `lspci --version`).
1. You have at least **Bash 4.x or higher** installed.
1. You have **already authorized your Thunderbolt EGPU** and are able to connect.
1. You have already **installed the latest (proprietary) drivers for your GPUs**.

## TL;DR
> **Please note**: There has been feedback that defining a specific internal GPU causes trouble in some cases (Especially when choosing *Intel Integrated Graphics*). It is therefore **not recommended** to specify the internal GPU (see #33, #28, #36, #37).

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

> **Hint**: Up until version `0.12.0` there was an [issue](https://github.com/hertg/egpu-switcher/issues/25) that the cleanup command needed the egpu-switcher to be set up for the cleanup to work. If you are having trouble to cleanup, remove or upgrade the egpu-switcher take a look at [this possible fix](https://github.com/hertg/egpu-switcher/issues/25#issuecomment-590728815).

## Commands
<pre>
<b>egpu-switcher setup</b> [--override] [--noprompt]
    This will generate your "xorg.conf.egpu" and "xorg.conf.internal" files and 
    symlink the "xorg.conf" file to one of them.

    It will also create the systemd service, that runs the "switch" command on each startup.

    The setup will NOT delete any of your files, if you have an existing "xorg.conf"
    file it will be backed up to "xorg.conf.backup.{datetime}". You can later revert
    it by executing the "cleanup" command (see below).

    <b>--override</b>
        If you have an AMD GPU or use the open-source nvidia drivers, 
        the "switch" command will prevent you from switching to the EGPU
        if there are no displays directly attached to it. This flag will
        make sure to switch to the EGPU even if there are no displays attached.

    <b>--noprompt</b>
        Prevent the setup from prompting for user interaction if there is
        no existing configuration file found. 
        (Is currently only used by the "postinst" script)
</pre>

<pre>
<b>egpu-switcher switch auto|egpu|internal</b> [--override]
    Switches to the specified GPU (egpu|internal). If the <b>auto</b> parameter 
    is used, the script will check if the egpu is attached and switch accordingly.
    
    You need to restart your computer (or display-manager) for this to take effect.

    <b>--override</b>
        If you have an AMD GPU or use the open-source nvidia drivers, 
        the "switch" command will prevent you from switching to the EGPU
        if there are no displays directly attached to it. This flag will
        make sure to switch to the EGPU even if there are no displays attached.
</pre>

<pre>
<b>egpu-switcher config</b>
    Prompts the user to specify their external/internal GPU and saves their answer
    to the configuration file.
</pre>

<pre>
<b>egpu-switcher cleanup</b> [--hard]
    Remove all files egpu-switcher has created previously and restore the backup
    of previous "xorg.conf" files.

    <b>--hard</b>
        Remove configuration files too.
</pre>

---

## Hints

1. **Ubuntu 19.04 or later**\
When installing Ubuntu 19.04 or later, please check the box "Install third-party software for graphics and Wi-Fi hardware". After that, all required drivers will be installed automatically.

1. **Notebooks with hybrid graphics**\
I am using a Lenovo notebook with hybrid graphics (internal graphics **and** a dedicated GPU). I've experienced freezes in the Ubuntu 19.04 installer which could only be resolved by changing the display settings in the BIOS from **Hybrid Graphics** to **Discrete Graphics**. After the installation was complete, i was able to change this setting back to **Hybrid Graphics**, without any issues.

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

[Service]
Type=oneshot
ExecStart=egpu-switcher switch auto

[Install]
WantedBy=multi-user.target
```

This will enable the automatic detection wheter your egpu is connected or not on startup.

## Build (notes to myself)
1. `sudo apt install devscripts`
1. `sudo apt install debhelper`
1. Update changelog: `dch`
1. Build: `debuild -S | tee /tmp/debuild.log 2>&1`
1. Upload to ppa: `dput ppa:hertg/egpu-switcher egpu-switcher_X.X.X_source.changes`
