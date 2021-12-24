#!/bin/bash
[ -f "$1" ] || { echo "Usage: $0 TEST_IN" >&2; exit 1; }

TEST_IN="$1"
TEST_GV="${TEST_IN}.gv"

grep -q '\.in$' <<<"$TEST_IN" || { echo "[-] '${TEST_IN}' is not a in file."; exit 1; }

cat >"$TEST_GV" <<-EOF
	digraph {
	$(sed 's/^\([^ ]*\) \(.*\)/  "\1" -> "\2"/g' "$TEST_IN")
	}
EOF
