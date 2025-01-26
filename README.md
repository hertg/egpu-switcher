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

The goal of this CLI is to lower the barrier for Linux users to
use their eGPU on the Linux Desktop. With the `egpu-switcher config`
command the user can choose their external GPU.
On every bootup the service will check if the eGPU is connected
and if so, make X.Org prefer it.

> [!NOTE]
> **Limitations**: No hotplugging is possible. Users still need to reboot their computer to connect / disconnect the eGPU.

## Requirements

- Running X.Org
- Thunderbolt connection to eGPU is authorized
- Necessary graphics drivers for eGPU are installed

## Installation

### Ubuntu (apt)

*The PPA is no longer maintained for now (see [#90](https://github.com/hertg/egpu-switcher/issues/90))*

### Arch (aur)

```bash
paru -S egpu-switcher
```

> [!TIP]
> :deciduous_tree::zap: Save time and energy by using the pre-compiled `egpu-switcher-bin` package


### Manual

#### Installation and setup

Download binary from [latest release](https://github.com/hertg/egpu-switcher/releases)

Copy binary to `/opt`, apply proper permissions, and link it in `/usr/bin`

```bash
sudo cp <downloaded-binary> /opt/egpu-switcher
sudo chmod 755 /opt/egpu-switcher
sudo ln -s /opt/egpu-switcher /usr/bin/egpu-switcher
sudo egpu-switcher enable
```

#### Uninstall

```bash
sudo egpu-switcher disable --hard
sudo rm /usr/bin/egpu-switcher
sudo rm /opt/egpu-switcher
```

### Build

#### Prerequisites

Install the [go toolchain](https://go.dev/doc/install)

#### Installation and setup

```bash
git clone git@github.com:hertg/egpu-switcher.git
cd egpu-switcher
make build -s
sudo make install -s
sudo egpu-switcher enable
```

#### Uninstall

```bash
sudo egpu-switcher disable --hard
sudo make uninstall -s
```


## Commands

```txt
Usage:
  egpu-switcher [command]

Available Commands:
  config      Choose your external GPU
  disable     Disable egpu-switcher from running at startup
  enable      Enable egpu-switcher to run at startup
  help        Help about any command
  switch      Check if eGPU is present and configure X.org accordingly
  version     Print version information

Flags:
  -h, --help      help for egpu-switcher
  -v, --verbose   verbose output

Use "egpu-switcher [command] --help" for more information about a command.


```


## Configuration

The config file is created automatically and can be found at `/etc/egpu-switcher/config.yaml`.
Below you can see an example of a configuration file, annotated with additional information.

```yaml
egpu:
    # the 'driver' and 'id' configs are generated by 'egpu-switcher config'.
    # you probably shouldn't change this manually unless you understand why.
    driver: amdgpu
    id: 1153611719250962689
    
    # OPTIONAL: do not load 'modesetting' in the egpu config
    nomodesetting: false

# OPTIONAL: how many times 'egpu-switcher switch auto' should retry finding the egpu.
# this can be helpful if the egpu takes some time to connect on your machine,
# the following values are the default.
detection:
  retries: 6
  interval: 500 # milliseconds

# OPTIONAL: if you want to execute a script after switching to egpu/internal.
# the values must be absolute paths to a shell script, this script will
# then be run with '/bin/sh $script'.
# 
# it is required that the script is owned by root (uid 0) 
# and has a permission of -rwx------ (0700).
hooks:
  internal: /home/michael/tmp/internal.sh
  egpu: /home/michael/tmp/egpu.sh
```

---

## Troubleshooting

If you run into problems, please have a look at
[TROUBLESHOOT.md](https://github.com/hertg/egpu-switcher/blob/master/TROUBLESHOOT.md)
before reporting any issues.
