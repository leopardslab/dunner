#!/bin/sh

set -e

REPO="dunner-deb"
PACKAGE="dunner"
DISTRIBUTIONS="stable"
COMPONENTS="main"

if [ -z "$USER" ]; then
  echo "USER is not set"
  exit 1
fi

if [ -z "$API_KEY" ]; then
  echo "API_KEY is not set"
  exit 1
fi

setVersion () {
  VERSION=$(curl --silent "https://api.github.com/repos/leopardslab/dunner/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/');
}

setUploadDirPath () {
  UPLOADDIRPATH="pool/d/$PACKAGE"
}

downloadDebianArtifacts() {
  echo "Dowloading debian artifacts"
  FILES=$(curl -s https://api.github.com/repos/leopardslab/dunner/releases/latest \
| grep "browser_download_url.*deb" \
| cut -d : -f 3 \
| sed -e 's/^/https:/' \
| tr -d '"' );
  echo $FILES
  for i in $FILES; do
    RESPONSE_CODE=$(curl -O  -w "%{response_code}" "$i")
    code=$(echo "$RESPONSE_CODE" | head -c2)
    if [ $code != "20" ] && [ $code != "30" ]; then
      echo "Unable to download $i HTTP response code: $RESPONSE_CODE"
    fi
  done;
  echo "Finished downloading"
}

bintrayUpload () {
  for i in $FILES; do
    FILENAME=${i##*/}
    ARCH=$(echo ${FILENAME##*_} | cut -d '.' -f 1)
    if [ $ARCH == "386" ]; then
      ARCH="i386"
    fi

    URL="https://api.bintray.com/content/leopardslab/$REPO/$PACKAGE/$VERSION/$UPLOADDIRPATH/$FILENAME;deb_distribution=$DISTRIBUTIONS;deb_component=$COMPONENTS;deb_architecture=$ARCH?publish=1&override=1"
    echo "Uploading $URL"

    RESPONSE=$(curl -T $FILENAME -u$USER:$API_KEY "$URL" -I -s -w "HTTPSTATUS:%{http_code}");
    HTTP_STATUS=$(echo $RESPONSE | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    echo "$RESPONSE"

    if [ $HTTP_STATUS != "20" ] && [ $HTTP_STATUS != "30" ]; then
      echo "Unable to upload, HTTP response code: $HTTP_STATUS"
      exit 1
    fi
    echo "HTTP response code: $HTTP_STATUS"
  done;
}

bintraySetDownloads () {
  for i in $FILES; do
    FILENAME=${i##*/}
    ARCH=$(echo ${FILENAME##*_} | cut -d '.' -f 1)
    if [ $ARCH == "386" ]; then
      ARCH="i386"
    fi
    URL="https://api.bintray.com/file_metadata/leopardslab/$REPO/$UPLOADDIRPATH/$FILENAME"

    echo "Putting $FILENAME in $PACKAGE's download list"
    RESPONSE=$(curl -X PUT -d "{ \"list_in_downloads\": true }" -H "Content-Type: application/json" -u$USER:$API_KEY "$URL" -s -w "HTTPSTATUS:%{http_code}");
    HTTP_STATUS=$(echo $RESPONSE | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    echo "$RESPONSE"

    if [ $HTTP_STATUS != "20" ]; then
        echo "Unable to put in download list, HTTP response: $RESPONSE"
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

cleanArtifacts () {
  rm -f "$(pwd)/*.deb"
}

cleanArtifacts
downloadDebianArtifacts
setVersion
printMeta
setUploadDirPath
bintrayUpload
snooze
bintraySetDownloads
