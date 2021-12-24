#!/bin/sh
[ -f "$1" ] || { echo "Usage: $0 TEST_INPUT" >&2; exit 1; }
which dot || { echo "[!] graphviz not installed." >&2; exit 1; }

TEST_INPUT="$1"
DOT="$(mktemp)"
PNG="$(mktemp)"

cat >"$DOT" <<EOF
digraph {
$(sed 's/^\([^ ]*\) \(.*\)/"\1" -> "\2"/g' "$TEST_INPUT")
}
EOF

dot -Tpng -o"$PNG" "$DOT"
xdg-open "$PNG"
sleep 60
rm -f "$PNG" "$DOT"
