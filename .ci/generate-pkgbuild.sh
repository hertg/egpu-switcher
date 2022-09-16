#!/bin/sh

# $1: Version
# $2: SHA256 Checksum

#### egpu-switcher ####

mkdir -p ./.pkgbuild/egpu-switcher
cp ./.pkg/aur/egpu-switcher/* ./.pkgbuild/egpu-switcher

cat << EOF > ./.pkgbuild/egpu-switcher/PKGBUILD
# Maintainer: hertg <aur@her.tg>
# This file is generated automatically
_pkgname=egpu-switcher
_pkgver=$1
EOF

cat ./.pkgbuild/egpu-switcher/PKGBUILD.template >> ./.pkgbuild/egpu-switcher/PKGBUILD
rm ./.pkgbuild/egpu-switcher/PKGBUILD.template

#### egpu-switcher-bin ####

mkdir -p ./.pkgbuild/egpu-switcher-bin
cp ./.pkg/aur/egpu-switcher-bin/* ./.pkgbuild/egpu-switcher-bin

cat << EOF > ./.pkgbuild/egpu-switcher-bin/PKGBUILD
# Maintainer: hertg <aur@her.tg>
# This file is generated automatically
_pkgname=egpu-switcher-bin
_pkgver=$1
_sha256sum=$2
_source=\${_pkgname}-\${_pkgver}::https://github.com/hertg/egpu-switcher/releases/download/\${_pkgver}/egpu-switcher-amd64
EOF

cat ./.pkgbuild/egpu-switcher-bin/PKGBUILD.template >> ./.pkgbuild/egpu-switcher-bin/PKGBUILD
rm ./.pkgbuild/egpu-switcher-bin/PKGBUILD.template


