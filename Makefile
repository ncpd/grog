BUILDKIT=DOCKER_BUILDKIT=1

build-armv7:
	$(BUILDKIT) docker build . --build-arg arch=armv7 -t olacin/grog

build-amd64:
	$(BUILDKIT) docker build . -t olacin/grog