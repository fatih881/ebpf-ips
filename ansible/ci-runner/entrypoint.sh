#!/bin/bash
set -e
if [ -z "$GH_PAT" ] || [ -z "$GH_REPO_URL" ]; then
    echo "Error: GH_PAT and GH_REPO_URL environment variables are required."
    exit 1
fi
REPO_NAME=$(echo "${GH_REPO_URL}" | sed -e 's|https://github.com/||' -e 's|/$||')

echo "--> Getting Registration Token for ${REPO_NAME}..."
REG_TOKEN=$(curl -sX POST \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer ${GH_PAT}" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  https://api.github.com/repos/${REPO_NAME}/actions/runners/registration-token | jq .token --raw-output)

if [ "$REG_TOKEN" == "null" ] || [ -z "$REG_TOKEN" ]; then
    echo "Error: Failed to get registration token. Check your PAT scopes (repo, workflow)."
    exit 1
fi

echo "--> Configuring Runner..."

./config.sh \
    --url "${GH_REPO_URL}" \
    --token "${REG_TOKEN}" \
    --name "${RUNNER_NAME:-$(hostname)}" \
    --work "_work" \
    --labels "self-hosted,linux,x64,docker,bare-metal" \
    --unattended \
    --replace


cleanup() {
    echo "--> Removing runner..."
    ./config.sh remove --token "${REG_TOKEN}"
}
trap 'cleanup; exit 130' INT
trap 'cleanup; exit 143' TERM

echo "--> Starting Runner..."
./run.sh & wait $!