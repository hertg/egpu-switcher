#!/bin/sh

version=$1
checksum=$2

#### egpu-switcher ####

mkdir -p ./.pkgbuild/egpu-switcher
cp ./.pkg/aur/egpu-switcher/* ./.pkgbuild/egpu-switcher

cat << EOF > ./.pkgbuild/egpu-switcher/PKGBUILD
# Maintainer: hertg <aur@her.tg>
# This file is generated automatically
_version=$version
_versionWithoutPrefix=${version#v}
_pkgname=egpu-switcher
_pkgver=$(echo $version | sed 's/\([^-]*-g\)/r\1/;s/-/./g')
_source=\${_pkgname}-\${_version}::https://github.com/hertg/egpu-switcher/archive/refs/tags/$version.tar.gz
EOF

cat ./.pkgbuild/egpu-switcher/PKGBUILD.template >> ./.pkgbuild/egpu-switcher/PKGBUILD
rm ./.pkgbuild/egpu-switcher/PKGBUILD.template


#### egpu-switcher-bin ####

mkdir -p ./.pkgbuild/egpu-switcher-bin
cp ./.pkg/aur/egpu-switcher-bin/* ./.pkgbuild/egpu-switcher-bin

cat << EOF > ./.pkgbuild/egpu-switcher-bin/PKGBUILD
# Maintainer: hertg <aur@her.tg>
# This file is generated automatically
_version=$version
_pkgname=egpu-switcher-bin
_pkgver=$(echo $version | sed 's/\([^-]*-g\)/r\1/;s/-/./g')
_sha256sum=$checksum
_source=\${_pkgname}-\${_pkgver}::https://github.com/hertg/egpu-switcher/releases/download/$version/egpu-switcher-amd64
EOF

cat ./.pkgbuild/egpu-switcher-bin/PKGBUILD.template >> ./.pkgbuild/egpu-switcher-bin/PKGBUILD
rm ./.pkgbuild/egpu-switcher-bin/PKGBUILD.template


