# Installation
> I am using this script with my notebook, which already has a dedicated gpu. I've experienced problems when working in "Hybrid Graphics" mode, meaning that both will be activated, the integrated graphics *and* the dedicated graphics. When booting from the Ubuntu Live-CD the installer freezed multiple times. After changing the BIOS display settings from "Hybrid Graphics" to "Dedicated Graphics", the problem was gone. After successful installation of Ubuntu, the BIOS setting can be set back to "Hybrid Graphics".\
\
> I am using a Lenovo ThinkPad, your device will probably have different BIOS settings.

Install the latest Ubuntu 19.04. You can choose the minimal installation, but make sure to check the box "Install third-party software for graphics and Wi-Fi hardware".

With Ubuntu 19.04 the latest proprietary nvidia-drivers will already be installed and you won't have to tinker around with the nouveau drivers.

# Setup
Connect your egpu to the computer and run the setup with `sudo egpu-switcher setup`. Choose your internal and external graphics card from the list.

`sudo nano /etc/systemd/system/egpu.service`
`sudo systemctl daemon-reload`
`sudo systemctl enable egpu`
