echo "Repository: ghcr.io/$GITHUB_REPOSITORY"
echo "SHA: $GITHUB_SHA"

. "$(werf ci-env github --as-file)"

if [[ $? -ne 0 ]]; then
  echo "Error: failed to source werf environment."
  exit 1
fi

werf export service --tag "ghcr.io/$GITHUB_REPOSITORY:$GITHUB_SHA"

if [[ $? -ne 0 ]]; then
  echo "Error: werf export service command failed."
  exit 1
fi
