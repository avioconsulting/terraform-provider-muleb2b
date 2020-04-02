#!/usr/bin/env sh
SCRIPT=`basename "$0"`
OS=unknown
ARCH=$(uname -m)
DOCKER_IMAGE=${DOCKER_IMAGE:-"golang:latest"}

function usage() {
    echo "Usage: $SCRIPT <version>"
    echo "  version - the version number to download or build"
}

# Check input
if [[ $# != 1 ]]; then
  usage
  exit 1
fi

VERSION=$1

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
  OS=linux
elif [[ "$OSTYPE" == "darwin"* ]]; then
  OS=darwin
elif [[ "$OSTYPE" == "cygwin" || "$OSTYPE" == "msys" ]]; then
  OS=windows
fi

if [[ "$ARCH" == "x86_64" ]]; then
  ARCH="amd64"
elif [[ "$ARCH" == "i386" ]]; then
  ARCH="386"
elif [[ "$ARCH" == "armv7"* ]]; then
  ARCH="arm"
elif [[ "$ARCH" == "arm64" ]]; then
  ARCH="arm64"
fi

echo "Downloading https://github.com/avioconsulting/terraform-provider-muleb2b/releases/download/${VERSION}/terraform-provider-muleb2b_${OS}_${ARCH}.tar.gz"
curl -LO "https://github.com/avioconsulting/terraform-provider-muleb2b/releases/download/${VERSION}/terraform-provider-muleb2b_${OS}_${ARCH}.tar.gz"
tar -xzf "terraform-provider-muleb2b_${OS}_${ARCH}.tar.gz"

if [[ OS != "windows" ]]; then
  if [[ ! -d ~/.terraform.d/plugins ]]; then
    mkdir -p ~/.terraform.d/plugins
  fi
  mv terraform-provider-muleb2b_${VERSION} ~/.terraform.d/plugins/
else
  if [[ ! -d $APPDATA/terraform.d/plugins ]]; then
    mkdir -p $APPDATA/terraform.d/plugins
  fi
  mv terraform-provider-muleb2b_${VERSION} $APPDATA/terraform.d/plugins/
fi

rm "terraform-provider-muleb2b_${OS}_${ARCH}.tar.gz"