#!/bin/bash
set -euo pipefail

selected=$(opencode session list --format json | jq -r '.[] | "\(.id) | \(.title)"' | fzf --delimiter="|" --with-nth=2.. --preview "echo {}" | cut -d' ' -f1)
echo "Selected session: $selected"

# if [[ -n "$selected" ]]; then
#   exec opencode -s "$selected"
# fi
