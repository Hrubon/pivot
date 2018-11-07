#/bin/sh -eu

prefix="pivot"
ip netns ls | grep -o "^$prefix[^ ]*" | xargs --no-run-if-empty -n1 ip netns del
