#!/bin/bash
echo "## Extracting app name and version from code in pkg/version"
SOURCE_CODE=pkg/version/version.go
APP_NAME=$(grep -E 'APP\s+=' $SOURCE_CODE| awk '{ print $3 }'  | tr -d '"')
APP_VERSION=$(grep -E 'VERSION\s+=' $SOURCE_CODE| awk '{ print $3 }'  | tr -d '"')
APP_REVISION=$(grep -E 'REVISION\s+=' $SOURCE_CODE| awk '{ print $3 }'  | tr -d '"')
echo "## Found APP: ${APP_NAME}, VERSION: ${APP_VERSION}, REVISION: ${APP_REVISION}  in source file ${SOURCE_CODE}"
export APP_VERSION APP_NAME APP_REVISION
