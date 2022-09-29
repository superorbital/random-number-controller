#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

kubectl config use-context colima

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

make -C ${DIR}/.. docker-build
make -C ${DIR}/.. install
make -C ${DIR}/.. deploy

make -C ${DIR}/.. e2e-test

make -C ${DIR}/.. undeploy

IGNORE_NOT_FOUND=true make -C ${DIR}/.. uninstall