#!/bin/sh

set -e

REPO="dunner-rpm"
PACKAGE="dunner"
GORELEASER_DIR="dist"

if [ -z "$USER" ]; then
  echo "USER is not set"
  exit 1
fi

if [ -z "$API_KEY" ]; then
  echo "API_KEY is not set"
  exit 1
fi

setVersion () {
  echo "Fetching latest dunner version from Github releases.."
  VERSION=$(curl --silent "https://api.github.com/repos/leopardslab/dunner/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/');
  echo "Latest dunner version: $VERSION"
  if [ "$VERSION" == "" ]; then
    exit 1
  fi
  VERSIONNUMBER=$(echo $VERSION | cut -d 'v' -f 2)
  FILES=( "dunner_$(echo $VERSIONNUMBER)_linux_arm.rpm" "dunner_$(echo $VERSIONNUMBER)_linux_arm64.rpm" "dunner_$(echo $VERSIONNUMBER)_linux_i386.rpm" "dunner_$(echo $VERSIONNUMBER)_linux_x86_64.rpm" )
}

setUploadDirPath () {
  UPLOADDIRPATH="$PACKAGE/$VERSION"
}

bintrayUpload () {
  for i in "${FILES[@]}"; do
    FILENAME=$i
    ARCH=$(echo ${FILENAME##*_} | cut -d '.' -f 1)
    if [ "$ARCH" == "386" ]; then
      ARCH="i386"
    fi
    if [ "$ARCH" == "64" ]; then
      ARCH="x86_64"
    fi

    URL="https://api.bintray.com/content/leopardslab/$REPO/$PACKAGE/$VERSION/$UPLOADDIRPATH/$FILENAME?publish=1&override=1"
    echo "Uploading $URL"

    RESPONSE=$(curl -T ./$GORELEASER_DIR/$FILENAME -u$USER:$API_KEY "$URL" -I -s -w "HTTPSTATUS:%{http_code}");
    HTTP_STATUS=$(echo $RESPONSE | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    echo "$RESPONSE"

    if [[ "$(echo $HTTP_STATUS | head -c2)" != "20" ]]; then
      echo "Unable to upload, HTTP response code: $HTTP_STATUS"
      exit 1
    fi
    echo "HTTP response code: $HTTP_STATUS"
  done;
}

bintraySetDownloads () {
  for i in "${FILES[@]}"; do
    FILENAME=$i
    ARCH=$(echo ${FILENAME##*_} | cut -d '.' -f 1)
    if [ "$ARCH" == "386" ]; then
      ARCH="i386"
    fi
    if [ "$ARCH" == "64" ]; then
      ARCH="x86_64"
    fi
    URL="https://api.bintray.com/file_metadata/leopardslab/$REPO/$UPLOADDIRPATH/$FILENAME"

    echo "Putting $FILENAME in $PACKAGE's download list"
    RESPONSE=$(curl -X PUT -d "{ \"list_in_downloads\": true }" -H "Content-Type: application/json" -u$USER:$API_KEY "$URL" -s -w "HTTPSTATUS:%{http_code}");
    HTTP_STATUS=$(echo $RESPONSE | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    echo "$RESPONSE"

    if [ "$(echo $HTTP_STATUS | head -c2)" != "20" ]; then
        echo "Unable to put in download list, HTTP response code: $HTTP_STATUS"
        exit 1
    fi
    echo "HTTP response code: $HTTP_STATUS"
  done
}

snooze () {
    echo "\nSleeping for 30 seconds. Have a coffee..."
    sleep 30s;
    echo "Done sleeping\n"
}

printMeta () {
    echo "Publishing: $PACKAGE"
    echo "Version to be uploaded: $VERSION"
}

setVersion
printMeta
setUploadDirPath
bintrayUpload
snooze
bintraySetDownloads

