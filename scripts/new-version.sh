#!/usr/bin/env bash
# Description: Create a new version

#set -x

function incr_version {
   local v="$1"
   echo "[$v]" | tr . , | jq -r '.[2]+=1 | @csv' | tr , .
}
cur_version="$(cat VERSION.md )"
new_version="$1"
if [ -z "$new_version" ] ; then
   echo "current version is $cur_version"
   echo "Usage: $0 [version]"
   echo "pass i to [version] to increment"
   exit 1
fi



# Increment version
if [ "$new_version" == "i" ] ; then
  new_version="$(incr_version "$cur_version")"
  if [ -z "$new_version" ] ; then
    echo "Error: Unable to increment new version"
    exit 1
  fi

  echo "next version is $new_version"
  #exit 0
fi

git_tag=v${new_version}
echo $new_version > VERSION.md
git commit -m" new version $git_tag " VERSION.md
git push

git tag $git_tag
git push origin $git_tag
git status
