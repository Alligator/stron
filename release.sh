#!/bin/sh
set +e

upload=0

OPTIND=1
while getopts "hu" opt; do
  case $opt in
    h)
      echo "usage: ./release.sh [-u]"
      exit 0
      ;;
    u)
      upload=1
      ;;
  esac
done

shift $((OPTIND0))

if [[ ! -d "builds" ]]; then
  mkdir builds
fi

if [[ $upload -eq 1 ]]; then
  git fetch --tags
  git tag | grep "^$version$"
  if [[ $? -ne 0 ]]; then
    git tag "$version"
    gh release create --draft -n "$version" "$version"
  fi
fi

build () {
  plat=$1
  arch=$2
  echo "$plat-$arch"

  binary="stron"
  if [[ $plat == "windows" ]]; then
    binary="stron.exe"
  fi

  GOOS=${plat} GOARCH=${arch} go build .

  archive="stron-${plat}-${arch}.tgz"
  if [[ $plat == "windows" ]]; then
    archive="stron-${plat}-${arch}.zip"
    zip "builds/$archive" "$binary"
  else
    tar --create --gzip --file="builds/$archive" "$binary"
  fi
  rm "$binary"

  if [[ $upload -eq 1 ]]; then
    gh release upload --clobber "$version" "builds/$archive"
  fi
}

build "darwin" "arm64"
build "darwin" "amd64"

build "linux" "arm64"
build "linux" "amd64"
build "linux" "386"

build "windows" "amd64"
build "windows" "386"
