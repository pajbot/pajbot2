#!/usr/bin/env bash
# this script is mostly yoinked from the built in checkout script CircleCI uses

# fail on process exit code != 0
set -e

# Workaround old docker images with incorrect $HOME
# check https://github.com/docker/docker/issues/2968 for details
if [ "${HOME}" = "/" ]; then
  export HOME=$(getent passwd $(id -un) | cut -d: -f6)
fi

# github.com and bitbucket.org identities
mkdir -p ~/.ssh
echo 'github.com ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAq2A7hRGmdnm9tUDbO9IDSwBK6TbQa+PXYPCPy6rbTrTtw7PHkccKrpp0yVhp5HdEIcKr6pLlVDBfOLX9QUsyCOV0wzfjIJNlGEYsdlLJizHhbn2mUjvSAHQqZETYP81eFzLQNnPHt4EVVUh7VfDESU84KezmD5QlWpXLmvU31/yMf+Se8xhHTvKSCZIFImWwoG6mbUoWf9nzpIoaSjB+weqqUUmpaaasXVal72J+UX2B+2RPW3RcT0eOzQgqlJL3RKrTJvdsjE3JEAvGq3lGHSZXy28G3skua2SmVi/w4yCE6gbODqnTWlg7+wC604ydGXA8VJiS5ap43JXiUFFAaQ==
bitbucket.org ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAubiN81eDcafrgMeLzaFPsw2kNvEcqTKl/VqLat/MaB33pZy0y3rJZtnqwR2qOOvbwKZYKiEO1O6VqNEBxKvJJelCq0dTXWT5pbO2gDXC6h6QDXCaHo6pOHGPUy+YBaGQRGuSusMEASYiWunYN0vCAI8QaXnWMXNMdFP3jHAJH0eDsoiGnLPBlBp4TNm6rYI74nMzgz3B9IikW4WVK+dc8KZJZWYjAuORU3jc1c/NPskD2ASinf8v3xnfXeukU0sJ5N6m5E8VLjObPEO+mN2t/FZTMZLiFqPWc/ALSqnMnnhwrNi2rbfg/rd/IpL8Le3pSBne8+seeFVBoGqzHM9yXw==
' >> ~/.ssh/known_hosts

# for private repositories
(umask 077; touch ~/.ssh/id_rsa)
chmod 0600 ~/.ssh/id_rsa
(cat <<EOF > ~/.ssh/id_rsa
$CHECKOUT_KEY
EOF
)

# do the actual clone
# $CIRCLE_REPOSITORY_URL is something like "git@github.com:pajlada/pajbot2.git"
# this INTENTIONALLY clones into the username folder "pajlada" because all imports
# are hardcoded to this github username in the go source (this is a go "feature")
git clone --recursive "$CIRCLE_REPOSITORY_URL" "$GOPATH/src/github.com/pajlada/$CIRCLE_PROJECT_REPONAME"
cd "$GOPATH/src/github.com/pajlada/$CIRCLE_PROJECT_REPONAME"

# fetch tag, if exists
if [ -n "$CIRCLE_TAG" ]; then
  git fetch --force origin "refs/tags/${CIRCLE_TAG}"
else
  git fetch --force origin "develop:remotes/origin/develop"
fi

# checkout specific tag or branch
if [ -n "$CIRCLE_TAG" ]; then
  git reset --hard "$CIRCLE_SHA1"
  git checkout -q "$CIRCLE_TAG"
elif [ -n "$CIRCLE_BRANCH" ]; then
  git reset --hard "$CIRCLE_SHA1"
  git checkout -q -B "$CIRCLE_BRANCH"
fi

# checkout specific SHA1 hash
git reset --hard "$CIRCLE_SHA1"
