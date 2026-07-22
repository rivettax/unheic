#!/usr/bin/env bash
set -euo pipefail

# Builds the Docker image from $IMAGE_CONTEXT, pushes it to the
# shared-services ECR repo $ECR_REPO tagged with the commit SHA, and publishes
# the pushed reference as the `image` meta-data key for the deploy steps.
#
# The registry account comes from the credentials the assume-role plugin
# already established, so the account ID lives in exactly one place per
# pipeline file: the role ARN.

account_id="$(aws sts get-caller-identity --query Account --output text)"
image="${account_id}.dkr.ecr.${AWS_REGION}.amazonaws.com/${ECR_REPO}"
tag="sha-${BUILDKITE_COMMIT}"

# --push makes BuildKit push straight from the builder: hosted agents build
# via a remote BuildKit driver, so the image never lands in the local daemon
# and a separate `docker push` would fail with "tag does not exist".
docker build --push --platform linux/amd64 -t "${image}:${tag}" "${IMAGE_CONTEXT}"
buildkite-agent meta-data set image "${image}:${tag}"
echo "Pushed ${image}:${tag}"
