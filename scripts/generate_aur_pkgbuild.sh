#/usr/bin/env bash

# Maintainer: UpCloud Team <hello at upcloud dot com>
# vim:set ts=2 sw=2

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <pkg_version>"
    exit 1
fi

PKG_VERSION="${1:1}"
PKG_HASH="SKIP"

cat > PKGBUILD <<EOF
pkgname=upcloud-cli
pkgver=${PKG_VERSION}
pkgrel=1
pkgdesc="upctl - a CLI tool for managing UpCloud services."
arch=('x86_64')
url="https://upcloud.com"
license=('Apache')
source=("https://github.com/UpCloudLtd/\${pkgname}/releases/download/v\${pkgver}/\${pkgname}_\${pkgver}_linux_x86_64.tar.gz")
sha256sums=('${PKG_HASH}')

package() {
    install -Dm755 upctl "\$pkgdir/usr/local/bin/upctl"
}
EOF

exit 0
