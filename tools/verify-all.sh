#!/bin/bash

# Copyright 2017 The authors
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

# This script will scan all md (markdown) files for bad references.
# It will look for strings of the form [...](...) and make sure that
# the (...) points to either a valid file in the source tree or, in the
# case of it being an http url, it'll make sure we don't get a 404.
#
# Usage: verify-links.sh [ dir | file ... ]
# default arg is root of our source tree

set -o errexit
set -o nounset
set -o pipefail

REPODIR=$(dirname "${BASH_SOURCE}")/..

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

exit $rc
