# linux-egpu-setup

> **Disclaimer**\
> Works only with NVIDIA cards and after the nvidia-drivers were installed.\
> Uses the `nvidia-xconfig` command to detect GPUs for setup and startup.\
> Not properly tested on different setups, use at your own risk.

This script will try to detect your GPUs and prompt you to answer which is the internal and which the external one.

It will then create a `xorg.conf.egpu` and a `xorg.conf.internal` file in your `/etc/X11` directory, with each of the GPUs defined via their BusID.

Additionally, an `init.d` script will be copied over to `/etc/init.d`, which will detect at startup if the egpu is connected. If the EGPU is connected the `/etc/X11/xorg.conf` file will be created as symlink to `/etc/X11/xorg.conf.egpu`. If no egpu is detected the symlink will point at the `/etc/X11/xorg.conf.internal` file.



# Prerequisites
1. Install the proprietary NVIDIA drivers

# Setup
1. Clone the git repo to your local machine or download as zip and unzip it in a directory on your computer.
1. Open the terminal and change directory into the script folder
1. Start the setup with `sudo ./setup.sh` and follow the instructions


