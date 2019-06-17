#!/bin/bash

# Copyright 2017 - 2019 The authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This script runs all the verification checks, both for CI and development
#
# Usage: verify-all.sh

set -o errexit
set -o nounset
set -o pipefail

REPODIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"/..; pwd -P)"

rc=0

echo Verifying hrefs
"${REPODIR}/tools/verify-links.sh" -v "${REPODIR}"/spec.md "${REPODIR}"/profile.md "${REPODIR}"/compatibility.md || rc=1

echo Verify tables
"${REPODIR}/tools/verify-tables.sh" -v "${REPODIR}"/spec.md "${REPODIR}"/profile.md "${REPODIR}"/compatibility.md || rc=1

echo Verify terminology and RFC keywords
"${REPODIR}/tools/verify-phrases.sh" -v "${REPODIR}"/spec.md "${REPODIR}"/profile.md "${REPODIR}"/compatibility.md || rc=1

echo Verify description json fields
"${REPODIR}/tools/verify-descriptions.sh" -v "${REPODIR}"/spec.md "${REPODIR}"/profile.md "${REPODIR}"/compatibility.md || rc=1

echo Verify spaces after periods
"${REPODIR}/tools/verify-spaces.sh" -v "${REPODIR}"/spec.md "${REPODIR}"/profile.md "${REPODIR}"/compatibility.md || rc=1

echo Verify OpenAPI
docker run --rm -v "${REPODIR}":/local openapitools/openapi-generator-cli validate -i /local/openapi.yaml || rc=1

echo Verify Swagger
docker run --rm -v "${REPODIR}":/local openapitools/openapi-generator-cli validate -i /local/swagger.yaml || rc=1

exit $rc
