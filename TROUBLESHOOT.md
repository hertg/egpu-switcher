# Troubleshooting

This is a collection of tips and common problems that might occur when trying to use an eGPU on the Linux Desktop. This document is not supposed to be a step-by-step guide on how to setup an eGPU, but rather provide some TL;DR assistance if you are running into problems.

Please feel free to provide your own solutions to problems you might have had by issuing a Pull Request to this repository.

## Common Problems
For general problems with your graphics card or graphics drivers, please do a web search first. Consider reading the corresponding articles in the ArchWiki (eg. [NVIDIA/Troubleshooting - ArchWiki](https://wiki.archlinux.org/index.php/NVIDIA/Troubleshooting)) in order to find a solution to your problem.

- **My eGPU is not recognized on `egpu-switcher setup` / `egpu-switcher config`**
  
  Please check for the usual errors first, maybe by trying the eGPU unit with a different device or operating system:
  - Does the GPU work?
  - Does the power supply of your eGPU case work?
  - Does your TB3 cable work?
  - Does your USB-C Port support TB3?

  If you can exclude the issues mentioned above, check for the following:
  - Does your system meet the [requirements](https://github.com/hertg/egpu-switcher#requirements)?
  - Is Thunderbolt enabled in your BIOS?\
  Also check whether you've set `Thunderbolt Security` to `USB-DP only` by accident.
  - Did you authorize your eGPU? (see [boltctl](https://www.mankier.com/1/boltctl))
  - Check if your eGPU shows up when running `lspci`. The egpu-switcher script is running this command internally to detect GPUs. If `lspci` doesn't find your GPU, neither will egpu-switcher.

- **eGPU is not selected automatically on bootup**
  - Make sure the GPU is connected **before** you start your computer.
  - Check if egpu-switcher even detects your eGPU at runtime by executing `egpu-switcher switch auto` with the eGPU connected. If it doesn't recognize the eGPU please see the troubleshooting tips above.
  - Try to re-run the `egpu-switcher config` and `egpu-switcher setup` command to reconfigure egpu-switcher.
  - Check whether the `egpu.service` is still enabled in systemd and check its output in `journalctl`.
  - In case there is a race-condition happening between `bolt` and `egpu-switcher`, try enabling `Pre-Boot ACL` in the BIOS and re-authorize your eGPU. With this setting enabled, your eGPU gets connected much faster. Be aware that this setting only makes sense with the Thunderbolt Security set to `user`  (see [#50](https://github.com/hertg/egpu-switcher/issues/50)).
  - Try changing your Thunderbolt Security Level and see if it changes anything (I personally use Thunderbolt Security Level `user` and `Pre-Boot ACL` enabled).


## Tips
Below you'll find a none exhaustive list of some general tips that may prevent you from running into certain problems. The list is based on personal experience and not necessarily on best-practices, please take into consideration that it might be outdated, depending on the time you read it.

- **Drivers**\
If you happen to have an NVIDIA GPU, it's preferrable to use the proprietary NVIDIA graphics drivers rather than Nouveau (see [NVIDIA - ArchWiki](https://wiki.archlinux.org/index.php/NVIDIA)).

- **Thunderbolt**\
There have been less issues reported when `Pre-Boot ACL` was enabled in the BIOS. Enabling `Pre-Boot ACL` allows authorized Thunderbolt devices to connect during pre-boot and leads to the eGPU connecting much faster on bootup, therefore limiting the impact of race-conditions between `egpu-switcher` and `bolt`. Please be aware that this only makes sense with Thunderbolt Security set to `user`.

- **Don't choose an explicit internal GPU**\
Although `egpu-switcher config` will ask you whether you want to specify a specific internal GPU, it's not recommended to do so. Many users have reported that specifying the internal GPU explicitly will cause problems in certain situations, especially with the intel drivers (see [#33](https://github.com/hertg/egpu-switcher/issues/33)).

- **Display Settings (BIOS)**\
If you are using a laptop with hybrid graphics (dedicated GPU + internal graphics) try changing the `Display Settings` in your BIOS if a problem occurs and see if it changes anything. I've personally experienced system freezes in distro installers (I think it was Ubuntu 19.04) that could only be solved by changing `Display Settings` from `Hybrid Graphics` to `Discrete Graphics` before starting the installation. After the installation went successful, the display settings could be changed back to `Hybrid Graphics` without any issues.