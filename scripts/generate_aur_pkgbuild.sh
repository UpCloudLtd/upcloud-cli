#/usr/bin/env bash

# Maintainer: UpCloud Team <contact at upcloud dot com>
# vim:set ts=2 sw=2

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <pkg_version>"
    exit 1
fi

PKG_VERSION=$1
PKG_HASH="SKIP"

cat > PKGBUILD <<EOF
pkgname=upctl
pkgver=${PKG_VERSION}
pkgrel=1
pkgdesc="UpCloud CLI."
arch=('x86_64')
url="https://upcloud.com"
license=('Apache')
makedepends=('go' 'git')
source=("https://github.com/UpCloudLtd/\$pkgname/archive/\${pkgver}.tar.gz")
sha256sums=('${PKG_HASH}')

build() {
  cd "\$pkgname-\$pkgver"
  make build
}

check() {
  cd "\$pkgname-\$pkgver"
  make test
}

package() {
    cd "\$pkgname-\$pkgver"
    install -Dm755 bin/upctl "\$pkgdir/usr/local/bin/upctl"
}
EOF

exit 0
