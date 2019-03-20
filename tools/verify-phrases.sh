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

# This script will scan all md (markdown) files for bad keyword usages.
#
# Usage: verify-phrases.sh [ dir | file ... ]
# default arg is root of our source tree

set -o errexit
set -o nounset
set -o pipefail

casePhrases=(   # case matters
"Service Binding" "Service Bindings"
"Service Broker" "Service Brokers"
"Service Offering" "Service Offerings"
"Service Plan" "Service Plans"
"Service Instance" "Service Instances"
MUST
"MUST NOT"
REQUIRED
SHALL
"SHALL NOT"
SHOULD
"SHOULD NOT"
RECOMMENDED
MAY
OPTIONAL
)

bannedPhrases=(  # case does not matter
"OSB API"
"Marketplace"
)

REPO_ROOT=$( cd $(dirname "${BASH_SOURCE}")/.. && pwd)

verbose=""
debug=""
stop=""

# Error file processing
err=tmpCC-$RANDOM
trap clean EXIT
function clean {
  rm -f ${err}* /tmp/tmpout
}

while [[ "$#" != "0" && "$1" == "-"* ]]; do
  opts="${1:1}"
  while [[ "$opts" != "" ]]; do
    case "${opts:0:1}" in
      v) verbose="1" ;;
      d) debug="1" ; verbose="1" ;;
      -) stop="1" ;;
      ?) echo "Usage: $0 [OPTION]... [DIR|FILE]..."
         echo "Verify all terms defined in spec are cased correctly."
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

arg=""

if [ "$*" == "" ]; then
  arg="${REPO_ROOT}"
fi

Files=$(find -L $* $arg \( -name "*.md" -o -name "*.htm*" \) | sort)

function checkFile {
  # Determine the max # of words we need to look for
  upperCasePhrases=( "" )
  upperBannedPhrases=( "" )
  maxWords=1
  i=0
  for phrase in "${casePhrases[@]}"; do
    words=( ${phrase[@]} )
    if (( ${#words[@]} > $maxWords )); then
      maxWords=${#words[@]}
    fi
	upper=$(echo $phrase | tr '[:lower:]' '[:upper:]')
    upperCasePhrases[i]=$upper
	((++i))
  done

  i=0
  for phrase in "${bannedPhrases[@]}"; do
	upper=$(echo $phrase | tr '[:lower:]' '[:upper:]')
    upperBannedPhrases[i]=$upper
	((++i))
  done

  lines=( "" )
  words=( "" )

  # Prepend each line of the file with its line number
  # echo $(date) start parsing
  cat -n $1 | while read num line ; do
    # Put each word on its own line with its line number before it.
    echo "$line" | \
    sed "s/(http[^[[:space:]]]*)|([a-zA-Z_\-]+)/ & /g" | \
	tr -s ' ' '\n' | \
	sed "s/^/$num /"
  done > /tmp/tmpout

  pairs=( "" )
  upperPairs=( "" )

  # echo $(date) start arraying
  i=0
  while read line word ; do
    [[ "${word}" == "" ]] && continue
    lines[i]=$line
	upperWord=$(echo -n "${word}" | tr '[:lower:]' '[:upper:]')

	pairs[i]=${word}
	upperPairs[i]=${upperWord}

	for (( j=0 ; j < maxWords-1 ; j++ )); do
	  if (( i > j )); then
	    pairs[i-1-j]="${pairs[i-1-j]} ${word}"
	    upperPairs[i-1-j]="${upperPairs[i-1-j]} ${upperWord}"
	  fi
	done

    ((++i))
  done < /tmp/tmpout

  # echo $(date) start scanning
  for (( i=0 ; i < ${#lines[@]} ; i++ )); do
    # echo $i ${pairs[i]}

	for (( j=0 ; j < ${#casePhrases[@]} ; j++ )); do
	  phrase=${casePhrases[j]}

	  if [[ "${upperPairs[i]} " == "${upperCasePhrases[j]} "* && \
	        "${pairs[i]} " != "${phrase} "* ]]; then
        ll=${pairs[i]}
        echo line ${lines[i]}: \'${ll:0:${#phrase}}\' should be \'${phrase}\'
      fi
	done

	for (( j=0 ; j < ${#upperBannedPhrases[@]} ; j++ )); do
	  phrase=${upperBannedPhrases[j]}
	  
	  # echo "${upperPairs[i]} "
	  if [[ "${upperPairs[i]} " == "${phrase} "* ]]; then
        ll=${pairs[i]}
        echo line ${lines[i]}: \'${ll:0:${#phrase}}\' is banned
      fi
	done
  done
}

for file in ${Files}; do
  # echo scanning $file
  dir=$(dirname $file)

  [[ -n "$verbose" ]] && echo "> $file"

  checkFile $file | tee -a $err
done

if [ -s ${err} ]; then exit 1 ; fi
