NAME=pkg
VERSION=0.1.0
MKFILE := $(abspath $(lastword $(MAKEFILE_LIST)))
DIR := $(dir $(MKFILE))


build:
	echo "Cross-compiling System Package Manager"
	GOOS=freebsd GOARCH=amd64 go build -o bin/freebsd/amd64/${NAME} main.go
	GOOS=freebsd GOARCH=arm64 go build -o bin/freebsd/arm64/${NAME} main.go
	GOOS=linux GOARCH=amd64 go build -o bin/linux/amd64/${NAME} main.go
	GOOS=linux GOARCH=arm64 go build -o bin/linux/arm64/${NAME} main.go
	GOOS=linux GOARCH=ppc64 go build -o bin/linux/ppc64/${NAME} main.go
	GOOS=linux GOARCH=mips64 go build -o bin/linux/mips64/${NAME} main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin/amd64/${NAME} main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/darwin/aarm64/${NAME} main.go
	GOOS=windows GOARCH=amd64 go build -o bin/windows/amd64/${NAME}.exe main.go
	GOOS=windows GOARCH=arm64 go build -o bin/windows/arm64/${NAME}.exe main.go

dep:
	go mod tidy

run:
	go run main.go

clean:
	go clean
	rm -rf ./bin
	rm -rf ./packages

package:
	mkdir -p packages/debian
	mkdir -p packages/fedora
	mkdir -p packages/tarball
	mkdir -p packages/darwin
	echo "Building Packages"
	fpm \
	  -s dir -t deb \
	  -p packages/debian/system-package-manager-${VERSION}-amd64.deb \
	  --architecture=amd64 \
	  --name=system-package-manager \
	  --license=gpl3 \
	  --version=${VERSION} \
	  --depends bash \
	  ${DIR}bin/linux/amd64/pkg=/usr/bin/pkg
	fpm \
          -s dir -t deb \
          -p packages/debian/system-package-manager-${VERSION}-arm64.deb \
          --architecture=arm64 \
          --name=system-package-manager \
          --license=gpl3 \
          --version=${VERSION} \
          --depends bash \
          ${DIR}bin/linux/amd64/pkg=/usr/bin/pkg
	fpm \
          -s dir -t rpm \
          -p packages/fedora/system-package-manager-${VERSION}-amd64.rpm \
          --architecture=amd64 \
          --name=system-package-manager \
          --license=gpl3 \
          --version=${VERSION} \
          --depends bash \
          ${DIR}bin/linux/amd64/pkg=/usr/bin/pkg
	fpm \
          -s dir -t rpm \
          -p packages/fedora/system-package-manager-${VERSION}-arm64.rpm \
          --architecture=arm64 \
          --name=system-package-manager \
          --license=gpl3 \
          --version=${VERSION} \
          --depends bash \
          ${DIR}bin/linux/amd64/pkg=/usr/bin/pkg
	fpm \
          -s dir -t osxpkg \
          -p packages/darwin/system-package-manager-${VERSION}-amd64.pkg \
          --architecture=amd64 \
          --osxpkg-identifier-prefix=org.systemos \
          --name=system-package-manager \
          --license=gpl3 \
          --version=${VERSION} \
          --depends bash \
          ${DIR}bin/darwin/amd64/pkg=/usr/bin/pkg
	fpm \
          -s dir -t osxpkg \
          -p packages/darwin/system-package-manager-${VERSION}-aarm64.pkg \
          --architecture=arm64 \
          --osxpkg-identifier-prefix=org.systemos \
          --name=system-package-manager \
          --license=gpl3 \
          --version=${VERSION} \
          --depends bash \
          ${DIR}bin/darwin/aarm64/pkg=/usr/bin/pkg
	fpm \
          -s dir -t tar \
          -p packages/tarball/system-package-manager-${VERSION}-amd64.tar.gz \
          --architecture=amd64 \
          --name=system-package-manager \
          --license=gpl3 \
          --version=${VERSION} \
          --depends bash \
          ${DIR}bin/linux/amd64/pkg=/usr/bin/pkg
	fpm \
          -s dir -t tar \
          -p packages/tarball/system-package-manager-${VERSION}-arm64.tar.gz \
          --architecture=arm64 \
          --name=system-package-manager \
          --license=gpl3 \
          --version=${VERSION} \
          --depends bash \
          ${DIR}bin/linux/amd64/pkg=/usr/bin/pkg

all: clean dep build package
