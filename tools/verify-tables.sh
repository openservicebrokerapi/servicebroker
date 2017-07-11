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

# This script will scan all md (markdown) files to make sure our tables
# are formatted correctly.
#
# Usage: verify-table.sh [ dir | file ... ]
# default arg is root of our source tree

set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT=$(dirname "${BASH_SOURCE}")/..

verbose=""
debug=""
stop=""

while [[ "$#" != "0" && "$1" == "-"* ]]; do
  opts="${1:1}"
  while [[ "$opts" != "" ]]; do
    case "${opts:0:1}" in
      v) verbose="1" ;;
      d) debug="1" ; verbose="1" ;;
      -) stop="1" ;;
      ?) echo "Usage: $0 [OPTION]... [DIR|FILE]..."
         echo "Verify all RFC2119 keywords in files."
         echo
         echo "  -v   show each file as it is checked"
         echo "  -?   show this help text"
         echo "  --   treat remainder of args as dir/files"
         exit 0 ;;
      *) echo "Unknown option '${opts:0:1}'"
         exit 1 ;;
    esac
    opts="${opts:1}"
  done
  shift
  if [[ "$stop" == "1" ]]; then
    break
  fi
done

# echo verbose:$verbose
# echo debug:$debug
# echo args:$*

# $1 == line number
# $2 == line
function checkForPunc() {
  if [[ "$2" != "|"* ]]; then
    return 0
  fi

  # Grab just the last column of the table and count how many words there are
  last=$(echo "$2" | sed "s/.*|\(.*\) |$/\1/")
  count=$(echo "$last" | wc -w)

  # Anything less than 2 words we ignore
  if (( "$count" < 2 )); then
    return 0
  fi

  if [[ "$last" == *"." || "$last" == *"?" ]]; then
    return 0
  fi

  echo "$file - $1: Last column doesn't end with a '. |' or '? |' - watch for extra spaces"
  return 1
}

# $1 == line number
# $2 == line
function checkForSpaces() {
  if [[ "$2" != "|"* ]]; then
    return 0
  fi

  if [[ "$line" != "| "* ]]; then
    echo "$file - $1: Missing a space after |"
    return 1
  fi

  if [[ "$line" != *" |" ]]; then
    echo "$file - $1: Missing a space before |"
    return 1
  fi

  if [[ "$line" =~ "|.*[^ ]|" ]]; then
    echo "$file - $1: Missing a space before |"
    return 1
  fi

  if [[ "$line" =~ "|.*|[^ ]" ]]; then
    echo "$file - $1: Missing a space after |"
    return 1
  fi

  if [[ "$line" == *"  "* ]]; then
    echo "$file - $1: Has a double-space in it"
    return 1
  fi

  return 0
}

arg=""

if [ "$*" == "" ]; then
  arg="${REPO_ROOT}"
fi

Files=$(find $* $arg \( -name "*.md" -o -name "*.htm*" \) | sort)

for file in ${Files}; do
  # echo scanning $file
  dir=$(dirname $file)

  [[ -n "$verbose" ]] && echo "Verifying: $file"

  lineNum=0
  cat ${file} | while read line; do 
    ((lineNum++)) || true

    checkForPunc ${lineNum} "${line}"
    checkForSpaces ${lineNum} "${line}"
  done
done
exit 0
