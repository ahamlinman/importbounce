# importbounce

importbounce is a Go import redirector designed to:

* Run on AWS Lambda (although you can also run it as a normal HTTP server)
* Serve redirects for multiple packages on one or more domain names
* Support dynamic re-configuration without re-deployment

## Configuration

On every request, importbounce reads a TOML configuration file from a local or
remote source and uses it to decide where to redirect. For every Go package
prefix, a repository root and user-facing web redirect can be configured. See
`importbounce.sample.toml` for details.

The location of the config file can be set with the `-config` flag or
`IMPORTBOUNCE_CONFIG_URL` environment variable. The value is a URL-style string
using one of the following schemes:

* `http://{path...}` or `https://{path...}` for HTTP or HTTPS
* `file://{path...}` to read from the local filesystem
* `s3://{bucket}/{key...}` to read from an Amazon S3 bucket (you must have
  appropriate AWS credentials configured in the environment)

## Deployment

This repo includes an AWS CloudFormation template (`CloudFormation.yaml`) that
deploys importbounce as an AWS Lambda function serving your custom domain name.
To use the template, run `make` to build a Lambda-compatible binary, use `aws
cloudformation package` to upload the binary to S3, then use `aws
cloudformation deploy` to create the stack. Before the stack can be fully
deployed, you will need to go into the AWS Certificate Manager and follow the
directions to validate the TLS certificate generated for your domain. You will
also need to upload your TOML configuration as `importbounce.toml` to the S3
bucket created by CloudFormation (the filename can be changed with a stack
parameter if desired).

Alternatively, you can run importbounce as a standard HTTP server by passing
the `-http` flag with a listening address (e.g. `-http 0.0.0.0:8080`).

## Future Work

* DNS validation for AWS Certificate Manager certificates instead of email
  validation.
* Mechanism so that accidentally deploying a bad config doesn't take down your
  whole domain.
* Making the CloudFormation template easier for new users to deploy.
