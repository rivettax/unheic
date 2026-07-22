#!/usr/bin/env bash
set -euo pipefail

# Hosted agents run jobs in a PTY without `less`, so any aws command whose
# stdout isn't captured would die with "Unable to redirect output to pager"
# (exit 253) — after the API call succeeds but before the wait/verify below.
export AWS_PAGER=""

# Rolls the ECS service $ECS_CLUSTER/$ECS_SERVICE to the image published by
# the build step (meta-data key `image`): registers a new revision of
# $TASK_FAMILY with $CONTAINER_NAME's image swapped, points the service at it,
# and waits for the rollout to stabilize. The ECS deployment circuit breaker
# rolls back automatically if the new tasks never go healthy.

image="$(buildkite-agent meta-data get image)"

task_def="$(aws ecs describe-task-definition \
  --task-definition "${TASK_FAMILY}" \
  --query taskDefinition --output json)"

# A container-name mismatch would make the image swap below a silent no-op —
# the old image would redeploy and report success — so fail fast instead.
if ! echo "${task_def}" | jq -e --arg name "${CONTAINER_NAME}" \
  '.containerDefinitions | any(.name == $name)' >/dev/null; then
  echo "container ${CONTAINER_NAME} not found in task definition ${TASK_FAMILY}" >&2
  exit 1
fi

# Strip the read-only fields describe returns; register-task-definition
# rejects them.
new_task_def="$(echo "${task_def}" | jq \
  --arg img "${image}" \
  --arg name "${CONTAINER_NAME}" \
  '.containerDefinitions |= map(if .name == $name then .image = $img else . end)
   | del(.taskDefinitionArn, .revision, .status, .requiresAttributes,
         .compatibilities, .registeredAt, .registeredBy)')"

new_arn="$(aws ecs register-task-definition \
  --cli-input-json "${new_task_def}" \
  --query taskDefinition.taskDefinitionArn --output text)"

aws ecs update-service \
  --cluster "${ECS_CLUSTER}" --service "${ECS_SERVICE}" \
  --task-definition "${new_arn}"

aws ecs wait services-stable \
  --cluster "${ECS_CLUSTER}" --services "${ECS_SERVICE}"

# services-stable also reports success when the deployment circuit breaker
# rolled back to the previous revision — the service is stable, just not on
# our image. Verify the primary deployment runs the new revision so a
# rollback fails this step (and blocks the staging stage in deploy-core).
primary="$(aws ecs describe-services \
  --cluster "${ECS_CLUSTER}" --services "${ECS_SERVICE}" \
  --query 'services[0].deployments[?status==`PRIMARY`].taskDefinition | [0]' \
  --output text)"
if [[ "${primary}" != "${new_arn}" ]]; then
  echo "deployment rolled back: ${ECS_SERVICE} is running ${primary}, expected ${new_arn}" >&2
  exit 1
fi

echo "Deployed ${image} to ${ECS_CLUSTER}/${ECS_SERVICE}"
