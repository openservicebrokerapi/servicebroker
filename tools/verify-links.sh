#!/bin/bash

# Copyright 2017 The Authors.
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

REPO_ROOT=$(dirname "${BASH_SOURCE}")/..

verbose=""
debugFlag=""
stop=""
tmp=/tmp/out${RANDOM}

trap clean EXIT
unset seenFiles
declare -A seenFiles
let numSeenFiles=0 || true

function debug {
  if [[ "$debugFlag" != "" ]]; then
    (>&2 echo $*)
  fi
}

function clean {
  rm -f ${tmp}*
}

while [[ "$#" != "0" && "$1" == "-"* ]]; do
  opts="${1:1}"
  while [[ "$opts" != "" ]]; do
    case "${opts:0:1}" in
      v) verbose="1" ;;
      d) debugFlag="1" ; verbose="1" ;;
      -) stop="1" ;;
      ?) echo "Usage: $0 [OPTION]... [DIR|FILE]..."
         echo "Verify all links in markdown files."
         echo
         echo "  -v   show each file as it is checked"
         echo "  -d   show each href as it is found"
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
# echo debugFlag:$debugFlag
# echo args:$*

arg=""

if [ "$*" == "" ]; then
  arg="${REPO_ROOT}"
fi

mdFiles=$(find $* $arg -name "*.md" | sort | grep -v vendor | grep -v glide)

clean
for file in ${mdFiles}; do
  # echo scanning $file
  dir=$(dirname $file)

  [[ -n "$verbose" ]] && echo "> $file"

  # Replace ) with )\n so that each possible href is on its own line.
  # Then only grab lines that have [..](..) in them - put results in tmp file.
  # If the file doesn't have any lines with [..](..) then skip this file
  # Steps:
  #  tr   - convert all \n to a space since newlines shouldn't change anything
  #  sed  - add a \n after each ) since ) ends what we're looking for.
  #         This makes it so that each href is on a line by itself
  #  sed  - prefix each line with a space so the grep can do [^\\]
  #  grep - find all lines that match [...](...)
  cat $file | \
    tr '\n' ' ' | \
    sed "s/)/)\n/g" | \
    sed "s/^/ /g" | \
    grep "[^\\]\[.*\](.*)" > ${tmp}1 || continue

  # This sed will extract the href portion of the [..](..) - meaning
  # the stuff in the parens.
  sed "s/.*\[*\]\([^()]*\)/\1/" < ${tmp}1 > ${tmp}2  || continue

  cat ${tmp}2 | while read line ; do
    # Strip off the leading and trailing parens, and then spaces
    ref=${line#*(}
    ref=${ref%)*}
    ref=$(echo $ref | sed "s/ *//" | sed "s/ *$//")

    # Show all hrefs - mainly for verifying in our tests
    debug "Found: '$ref'"

    # An external href (ie. starts with http)
    if [ "${ref:0:4}" == "http" ]; then
      if ! curl -f -s -k --connect-timeout 10 ${ref} \
          > /dev/null 2>&1 ; then
        echo $file: Can\'t load url: ${ref} | tee -a ${tmp}3
      fi
      continue
    fi

    # Skip "mailto:" refs
    if [ "${ref:0:7}" == "mailto:" ]; then
      continue
    fi

    # Local file link (i.e. ref contains a #)
    if [[ "${ref/\#}" != "${ref}" ]]; then

      # If ref doesn't start with "#" then update filepath
      if [ "${ref:0:1}" != "#" ]; then
        # Split ref into filepath and the section link
        reffile=$(echo ${ref} | awk -F"#" '{print $1}')
        fullpath=${dir}/${reffile}
        ref=$(echo ${ref} | awk -F"#" '{$1=""; print $0}')
      else
        fullpath=${file}
        ref=${ref:1}
      fi

      if [[ ! -e "${fullpath}" ]]; then
        echo "$file: Can't find referenced file '${fullpath}'" | \
          tee -a ${tmp}3
        continue
      fi

      # Remove leading and trailing spaces
      ref=$(echo ${ref} | sed 's/^[[:space:]]*//' | sed 's/[[:space:]]*$//')

      # If we've seen this file before then grab its processed tmp file
      if [[ "${seenFiles[${fullpath}]+x}" != "" ]]; then
        anchorFile="${seenFiles[${fullpath}]}"
      else
        let numSeenFiles+=1
        anchorFile="${tmp}-anchors-${numSeenFiles}"
        seenFiles[${fullpath}]="${anchorFile}"

        # Search file for sections
        unset used
        declare -A used  # anchors used, seen+twiddled ones

        # Find all section headers in the file.
        # Remove leading & trailing spaces.
        # Lower case it.
        # Convert spaces to "-".
        # Drop all non alphanumeric chars.
        # Twiddle section anchor if we've seen it before.
        grep "^[[:space:]]*#" < ${fullpath} | \
          sed 's/[[:space:]]*##*[[:space:]]*//' | \
          sed 's/[[:space:]]*$//' | \
          tr '[:upper:]' '[:lower:]' | \
          sed "s/  */-/g" | \
          sed "s/[^-a-zA-Z0-9]//g" | while read section ; do
            # If we haven't used this exact anchor before just use it now
            if [[ "${used[${section}]+x}" == "" ]]; then
              anchor=${section}
            else
              # We've used this anchor before so add "-#" to the end.
              # Keep adding 1 to "#" until we find a free spot.
              let num=1
              while true; do
                anchor="${section}-${num}"
                if [[ "${used[${anchor}]+x}" == "" ]]; then
                  break
                fi
                let num+=1
              done
            fi

            echo "${anchor}"
            used[${anchor}]="1"

            debug "Mapped section '${section}' to '${anchor}'"

          done > ${anchorFile} || true

        # Add sections of the form <a name="xxx">
        grep "<a name=" <${fullpath} | \
          sed 's/<a name="/\n<a name="/g' | \
          sed 's/^.*<a name="\(.*\)">.*$/\1/' | \
          sort | uniq >> ${anchorFile} || true

        # echo sections ; cat ${tmp}sections1
      fi

      # Skip refs of the form #L<num> and assume its pointing to a line
      # number of a file and those don't have anchors
      if [[ "${ref}" =~ ^L([0-9])+$ ]]; then
        continue
      fi

      # Finally, look for the ref in the list of sections/anchors
      debug "Anchor file(${fullpath}): ${anchorFile}"
      if ! grep "^${ref}$" ${anchorFile} > /dev/null 2>&1 ; then
        echo $file: Can\'t find section \'\#${ref}\' in ${fullpath} | \
          tee -a ${tmp}3
      fi

      continue

    fi

    newPath=${dir}/${ref}

    # And finally make sure the file is there
    # debug line: echo ref: $ref "->" $newPath
    if [[ ! -e "${newPath}" ]]; then
      echo $file: Can\'t find: ${newPath} | tee -a ${tmp}3
    fi

  done
done
rc=0
if [ -a ${tmp}3 ]; then
  rc=1
fi
exit $rc
