> **Note**: The main branch is now tracking a completely new version of egpu-switcher
> that hasn't been fully released yet. To get the `README` of the most recent release
> head over to [the `legacy` branch](https://github.com/hertg/egpu-switcher/tree/legacy).

<div align="center">
  <h1><strong>egpu-switcher</strong></h1>
  <p>
		<strong>Distribution agnostic eGPU script that works with NVIDIA and AMD cards.</strong>
  </p>
  <p>
    <!--<a href="https://goreportcard.com/report/github.com/hertg/egpu-switcher">
      <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/hertg/egpu-switcher" />
    </a>-->
    <a href="#">
			<img alt="License Information" src="https://img.shields.io/github/license/hertg/egpu-switcher">
    </a>
  </p>
</div>

## Description
The goal of this CLI is to lower the barrier for Linux users to use their eGPU on the Linux Desktop.
With the `egpu-switcher config` command the user can choose their external GPU. On every bootup the service will check if the eGPU is connected and if so, make X.Org prefer it.

---

## Limitations
- No hotplugging is possible. Users still need to reboot their computer to connect / disconnect the eGPU.

---

## Requirements
- Running X.Org
- Thunderbolt connection to eGPU is authorized
- Necessary graphics drivers for eGPU are installed

---

## Installation

### Ubuntu (apt)
*TODO*

### Arch (aur)
*TODO*

### Manual
Installation and setup:
```bash
$ git clone git@github.com:hertg/egpu-switcher.git
$ cd egpu-switcher
$ make install
$ sudo egpu-switcher setup
```

Uninstall:
```bash
$ sudo egpu-switcher cleanup
$ make uninstall
```

---

## Commands
*TODO*

---

## Troubleshooting
If you run into problems, please have a look at [TROUBLESHOOT.md](https://github.com/hertg/egpu-switcher/blob/master/TROUBLESHOOT.md) before reporting any issues.
