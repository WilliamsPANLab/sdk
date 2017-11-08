#!/usr/bin/env bash
set -euo pipefail
unset CDPATH; cd "$( dirname "${BASH_SOURCE[0]}" )"; cd "$(pwd -P)"

USAGE=$(cat <<USAGE_EOF
    Create github release, and upload all files in folder as build artifacts.

    Usage:   $( basename "$0" ) <github-auth-token> <github-tag> <artifact-dir>
USAGE_EOF
)

if [ "$#" -ne 3 ]; then
    >&2 echo "Invalid number of positional arguments"
    >&2 echo "$USAGE"
    exit 1
fi

GITHUB_AUTH_TOKEN="${1}"
TAG_LABEL="${2}"
DIST_FILES="${3}/*"


function CreateRelease {(
    GITHUB_AUTH_TOKEN="${1}"
    TAG_LABEL="${2}"

    release_json=$(cat <<SETVAR
{
"tag_name": "${TAG_LABEL}",
"name": "${TAG_LABEL}",
"body": "Description of the release",
"draft": true,
"prerelease": false
}
SETVAR
    )

    curl -X POST \
        --data "${release_json}" \
        -H "Authorization: token ${GITHUB_AUTH_TOKEN}" \
        -H "Accept: application/vnd.github.v3+json" \
        "https://api.github.com/repos/flywheel-io/sdk/releases" \
    | jq -r '.url'
)}


function UploadReleaseArtifact {(
    GITHUB_AUTH_TOKEN="${1}"
    RELEASE_URL="${2}"
    ARTIFACT_FILE_PATH="${3}"

    upload=$(
    curl  -H "Authorization: token ${GITHUB_AUTH_TOKEN}" \
            -H "Accept: application/vnd.github.v3+json" \
            "$RELEASE_URL" \
    | jq -r '.upload_url' | cut -d "{" -f1
    )

    curl -X POST \
        -H "Authorization: token ${GITHUB_AUTH_TOKEN}" \
        -H "Accept: application/vnd.github.v3+json" \
        -H "Content-Type: $(file -b --mime-type ${ARTIFACT_FILE_PATH})" \
        --data-binary "@${ARTIFACT_FILE_PATH}" \
        "${upload}?name=$(basename ${ARTIFACT_FILE_PATH})"

)}


REL_URL=$(CreateRelease "${GITHUB_AUTH_TOKEN}" "${TAG_LABEL}")

for f in ${DIST_FILES}
do
  >&2 echo "Processing ${f} file..."
  UploadReleaseArtifact "${GITHUB_AUTH_TOKEN}" "${REL_URL}" "${f}"
done
