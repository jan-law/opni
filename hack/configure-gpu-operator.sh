#!/bin/bash

config_dir="$(realpath $(dirname $0)/../config)"
cd "${config_dir}"

tmp="$(mktemp -d)"

git clone --depth=1 "https://gitlab.com/nvidia/kubernetes/gpu-operator.git" "${tmp}"
pushd "${tmp}" &>/dev/null

pushd "config" &>/dev/null
bases=$(find ./crd/bases -name "*.yaml")
sed -i 's/nvidia.com/nvidia.opni.io/g' "${bases}"
rename "nvidia.com" "nvidia.opni.io" "${bases}"
bases=$(find ./crd/bases -name "*.yaml")
rm -rf "${config_dir}/crd/nvidia"
mkdir -p "${config_dir}"/crd/nvidia
mv ${bases} "${config_dir}/crd/nvidia"
popd &>/dev/null

rm -rf "${config_dir}/assets/gpu-operator"
mv assets "${config_dir}/assets/gpu-operator"

popd &>/dev/null
