# Set these to the desired values
ARTIFACT_ID=redsage
VERSION=0.1.0

GOTAG=1.14.13
# overwrite ADDITIONAL_LDFLAGS to disable static compilation
# this should fix https://github.com/golang/go/issues/13470
ADDITIONAL_LDFLAGS=""
MAKEFILES_VERSION=4.4.0
.DEFAULT_GOAL:=default

include build/make/variables.mk

include build/make/info.mk
include build/make/dependencies-gomod.mk
include build/make/build.mk
include build/make/test-common.mk
include build/make/test-integration.mk
include build/make/test-unit.mk
include build/make/static-analysis.mk
include build/make/clean.mk
include build/make/package-debian.mk
include build/make/deploy-debian.mk
include build/make/digital-signature.mk
include build/make/self-update.mk

default: compile