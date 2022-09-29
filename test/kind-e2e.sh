#!/bin/bash

set -o errexit
set -o pipefail

KIND_LOG_LEVEL="1"

if ! [ -z $DEBUG ]; then
  set -x
  KIND_LOG_LEVEL="6"
fi

set -o nounset
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cleanup() {
  if [[ "${KUBETEST_IN_DOCKER:-}" == "true" ]]; then
    kind "export" logs --name ${KIND_CLUSTER_NAME} "${ARTIFACTS}/logs" || true
  fi

  kind delete cluster \
    --verbosity=${KIND_LOG_LEVEL} \
    --name ${KIND_CLUSTER_NAME}

  git restore ${DIR}/../config/manager/kustomization.yaml
}

trap cleanup EXIT

export KIND_CLUSTER_NAME=${KIND_CLUSTER_NAME:-random-number-controller}

if ! command -v kind --version &> /dev/null; then
  echo "kind is not installed. Use the package manager or visit the official site https://kind.sigs.k8s.io/"
  exit 1
fi

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
export KUBECONFIG="${KUBECONFIG:-$HOME/.kube/kind-config-$KIND_CLUSTER_NAME}"

if [ "${SKIP_CLUSTER_CREATION:-false}" = "false" ]; then
  echo "[dev-env] creating Kubernetes cluster with kind"

  # find image tags at
  # https://github.com/kubernetes-sigs/kind/releases
  export K8S_VERSION=${K8S_VERSION:-v1.24.6@sha256:97e8d00bc37a7598a0b32d1fabd155a96355c49fa0d4d4790aab0f161bf31be1}

  kind create cluster \
    --verbosity=${KIND_LOG_LEVEL} \
    --name ${KIND_CLUSTER_NAME} \
    --config ${DIR}/kind.yaml \
    --retain \
    --image "kindest/node:${K8S_VERSION}"

  echo "Kubernetes cluster:"
  kubectl get nodes -o wide
fi


export IMG=controller:kind-e2e

make -C ${DIR}/.. docker-build
make -C ${DIR}/.. install

KIND_WORKERS=$(kind get nodes --name="${KIND_CLUSTER_NAME}" | grep worker | awk '{printf (NR>1?",":"") $1}')

kind load docker-image --name="${KIND_CLUSTER_NAME}" --nodes=${KIND_WORKERS} ${IMG}

make -C ${DIR}/.. deploy

kubectl rollout status -n random-number-controller-system deploy/random-number-controller-controller-manager --timeout=60s

make -C ${DIR}/.. e2e-test

make -C ${DIR}/.. undeploy

IGNORE_NOT_FOUND=true make -C ${DIR}/.. uninstall
