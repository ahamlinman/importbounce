# This file is meant to be sourced.
# shellcheck disable=SC2034

# Because this stack generates a TLS certificate for use by CloudFront, it can
# only be deployed to the us-east-1 region. This setting overrides the region
# in the current AWS profile.
export AWS_DEFAULT_REGION=us-east-1

project=importbounce
go_src=../cmd/importbounce

params_usage="$(cat <<EOF
  DomainName=<domain>  The domain name served by the redirector.
EOF
)"

print-stack-output () (
  stack_name="$1"
  echo "Upload your importbounce configuration to:"
  aws cloudformation describe-stacks \
    --stack-name "$stack_name" \
    --output text \
    --query "Stacks[0].Outputs[?OutputKey=='ConfigS3URI'].OutputValue"
  echo
  echo "Point your CNAME to:"
  aws cloudformation describe-stacks \
    --stack-name "$stack_name" \
    --output text \
    --query "Stacks[0].Outputs[?OutputKey=='ApiDomain'].OutputValue"
)
