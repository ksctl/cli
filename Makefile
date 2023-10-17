GOOS_LINUX = linux
GOOS_WINDOWS = windows
GOOS_MACOS = darwin

GOARCH_LINUX = amd64
GOARCH_WINDOWS = amd64
GOARCH_MACOS = arm64
GOARCH_MACOS_INTEL = amd64

CURR_TIME = $(shell date +%s)

install_linux:
	@echo "Started to Install ksctl"
	cd scripts && \
		env GOOS=${GOOS_LINUX} GOARCH=${GOARCH_LINUX} ./builder.sh

install_macos:
	@echo "Started to Install ksctl"
	cd scripts && \
		env GOOS=${GOOS_MACOS} GOARCH=${GOARCH_MACOS} ./builder.sh

install_macos_intel:
	@echo "Started to Install ksctl"
	cd scripts && \
		env GOOS=${GOOS_MACOS} GOARCH=${GOARCH_MACOS_INTEL} ./builder.sh

uninstall:
	@echo "Started to Uninstall ksctl"
	cd scripts && \
		./uninstall.sh
