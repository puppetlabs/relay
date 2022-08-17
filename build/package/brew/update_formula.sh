#!/bin/bash
GITHUB_TOKEN=$(ni get -p {.github_token})
TAG=$(ni get -p {.tag})
SHA=$(ni get -p {.sha})
if [[ $GITHUB_TOKEN =~ .{25,} ]] && [[ $TAG =~ ^v ]] && [[ $SHA =~ .{63,} ]] ; then
  git clone https://${GITHUB_TOKEN}@github.com/puppetlabs/homebrew-puppet
  cd homebrew-puppet
  git config user.name "Relay Autobot" && git config user.email "relay@users.noreply.github.com"
  PUBLISH_BRANCH=relay_${TAG}
  git checkout -b ${PUBLISH_BRANCH}
  sed -e "s/version \".*\"/version \"${TAG}\"/g" -i ./Formula/relay.rb
  sed -e "s/sha256 \".*\"/sha256 \"${SHA}\"/g" -i ./Formula/relay.rb
  COMMIT_MESSAGE="Update Relay to tagged version ${TAG}"
  git commit -am "${COMMIT_MESSAGE}"
  git push origin ${PUBLISH_BRANCH}
  PULLS_URI="https://api.github.com/repos/puppetlabs/homebrew-puppet/pulls"
  AUTH_HEADER="Authorization: token $GITHUB_TOKEN"
  NEW_PR_RESP=$(curl --data "{\"title\": \"${COMMIT_MESSAGE}\", \"head\": \"${PUBLISH_BRANCH}\", \"base\": \"main\"}" -X POST -s -H "${AUTH_HEADER}" ${PULLS_URI})
  if [[ $? == 0 ]]; then
    PR_URL=$(echo $NEW_PR_RESP | jq ._links.html.href)
  ni output set --key result --value "Success! Pull request for $TAG submitted at: $PR_URL"
  else
    ni output set --key result --value "error submitting pull request: $NEW_PR_RESP"
    exit 1
  fi
else
  ni output set --key result --value "bad input for one or more of: tag: [$TAG], sha: [$SHA], token: [sha256:$(echo $GITHUB_TOKEN | sha256sum)]"
  exit 1
fi
