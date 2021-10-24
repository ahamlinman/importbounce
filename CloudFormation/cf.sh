#!/usr/bin/env bash
set -euo pipefail

usage () {
  cat <<EOF
cf.sh - Deploy importbounce to your AWS account using CloudFormation

You must install the AWS CLI and Skopeo to use this script.

$0 build
  (Re)build the container image that will be deployed to AWS Lambda.

$0 upload <ECR repository>
  Upload the built container image to the provided ECR repository with Skopeo
  using a unique tag.

$0 deploy <stack name> [overrides...]
  Deploy the latest version of the CloudFormation stack using the latest
  container image, then print the URL for the deployed API.

  Additional arguments are passed to the "--parameter-overrides" option of "aws
  cloudformation deploy". When deploying the stack for the first time, pass
  "DomainName=<domain>" to set the domain name of the redirector.

$0 build-deploy <ECR repository> <stack name> [overrides...]
  Build, upload, and deploy all in one step.

$0 help
  Print this message.
EOF
}

build () (
  os=linux
  arch=arm64

  set -x

  CGO_ENABLED=0 GOOS=$os GOARCH=$arch \
    go build -v \
    -ldflags='-s -w' \
    -o importbounce \
    ../cmd/importbounce

  go run go.alexhamlin.co/zeroimage@main \
    -os $os -arch $arch \
    importbounce
)

upload () (
  if ! type skopeo &>/dev/null; then
    echo "must install skopeo to upload container images" 1>&2
    return 1
  fi
  if [ ! -s importbounce.tar ]; then
    echo "must build a container image before uploading" 1>&2
    return 1
  fi

  repository_name="$1"
  repository="$(aws ecr describe-repositories \
    --region us-east-1 \
    --repository-names "$repository_name" \
    --query 'repositories[0].repositoryUri' \
    --output text)"
  registry="${repository%%/*}"
  tag="$(date +%s)"
  image="$repository:$tag"

  set -x
  if ! skopeo list-tags docker://"$repository" &>/dev/null; then
    aws ecr get-login-password --region us-east-1 \
    | skopeo login --username AWS --password-stdin "$registry"
  fi

  skopeo copy oci-archive:importbounce.tar docker://"$image"
  echo "$image" > latest-image.txt
)

deploy () (
  if [ ! -f latest-image.txt ]; then
    echo "must upload a container image before deploying" 1>&2
    return 1
  fi

  stack_name="$1"
  shift

  (
    set -x
    aws cloudformation deploy \
      --region us-east-1 \
      --template-file Template.yaml \
      --capabilities CAPABILITY_IAM \
      --stack-name "$stack_name" \
      --no-fail-on-empty-changeset \
      --parameter-overrides \
          ImageUri="$(cat latest-image.txt)" \
          "$@"
  )

  echo -e "\\nUpload your importbounce configuration to:"
  aws cloudformation describe-stacks \
    --region us-east-1 \
    --stack-name "$stack_name" \
    --output text \
    --query "Stacks[0].Outputs[?OutputKey=='ConfigS3URI'] | [0].OutputValue"

  echo -e "\\nPoint your CNAME to:"
  aws cloudformation describe-stacks \
    --region us-east-1 \
    --stack-name "$stack_name" \
    --output text \
    --query "Stacks[0].Outputs[?OutputKey=='ApiDomain'] | [0].OutputValue"
)

build-deploy () {
  local ecr_repository="$1"
  local stack_name="$2"
  shift 2

  build
  upload "$ecr_repository"
  deploy "$stack_name" "$@"
}

cmd="${1:-help}"
[ "$#" -gt 0 ] && shift

case "$cmd" in
  build)
    build
    ;;
  upload)
    upload "$@"
    ;;
  deploy)
    deploy "$@"
    ;;
  build-deploy)
    build-deploy "$@"
    ;;
  help)
    usage
    ;;
  *)
    usage
    exit 1
    ;;
esac
