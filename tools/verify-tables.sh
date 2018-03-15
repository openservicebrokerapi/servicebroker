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

# Error processing
trap clean EXIT
err=tmp-${RANDOM}
function clean {
  rm -f ${err}*
}

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
# $3 == previous line
function checkForPunc() {
  if [[ "$2" != "|"* ]]; then
    return
  fi

  # Split the first row of each table
  if [[ "$3" == "" ]]; then
    return
  fi

  if [[ "$2" != *" |" ]]; then
    echo "$file - $1: row doesn't end with ' |' - watch for trailing spaces"
    return
  fi

  line=${2#|}      # remove leading |
  line=${line%|}   # remove trailing |
  colNum=0
  # Split the row into columns - loop while "line" isn't empty
  while [[ "${line}" != "" ]]; do
    colNum=$((colNum+1))

    # Grab just the first column of the row
    col=${line%%|*}      # remove everything after the |
    len=${#col}
    line=${line:len+1}   # remove leading "len+1" chars

    # Skip the first column of each table row
    if (( "$colNum" < 2 )); then
      continue
    fi

    # Anything less than 4 words we ignore
    count=$(echo "$col" | sed "s/(.*)//g" | sed "s/\[.*\]//g" | wc -w)
    if (( "$count" < 4 )); then
      continue
    fi

    if [[ "$col" == *". " || "$col" == *"? " ]]; then
      continue
    fi

    echo "$file - $1: column $colNum has more than 3 words and doesn't end with a '. |' or '? |' - watch for extra/missing spaces."

  done
}

# $1 == line number
# $2 == line
function checkForSpaces() {
  if [[ "$2" != "|"* ]]; then
    return
  fi

  if [[ "$line" != "| "* ]]; then
    echo "$file - $1: Missing a space after leading |"
    return
  fi

  if [[ "$line" != *" |" ]]; then
    echo "$file - $1: Line doesn't end with ' |'"
    return
  fi

  if [[ "$line" =~ "|.*[^ ]|" ]]; then
    echo "$file - $1: Missing a space before |"
    return
  fi

  if [[ "$line" =~ "|.*|[^ ]" ]]; then
    echo "$file - $1: Missing a space after |"
    return
  fi

  if [[ "$line" == *"  "* ]]; then
    echo "$file - $1: Has a double-space in it"
    return
  fi
}

arg=""

if [ "$*" == "" ]; then
  arg="${REPO_ROOT}"
fi

Files=$(find $* $arg \( -name "*.md" -o -name "*.htm*" \) | sort)

for file in ${Files}; do
  # echo scanning $file
  dir=$(dirname $file)

  [[ -n "$verbose" ]] && echo "> $file"

  lineNum=0
  previous=""
  cat ${file} | while read line; do
    ((lineNum++)) || true

    checkForPunc ${lineNum} "${line}" "${previous}" | tee -a ${err}
    checkForSpaces ${lineNum} "${line}" | tee -a ${err}
    previous="${line}"
  done
done

if [ -s ${err} ]; then exit 1 ; fi
