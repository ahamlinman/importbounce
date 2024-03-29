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

This repository includes a CloudFormation template (`CloudFormation.yaml`) and
a helper script (`hfc`) to deploy importbounce as an AWS Lambda function in the
us-east-1 region serving a single custom domain name through CloudFront.

To prepare the deployment:

1. Install and configure the AWS CLI, and ensure that it has access to S3 and
   CloudFormation.
2. Create an S3 bucket in the us-east-1 region to hold the Lambda function code.
3. Copy `hfc.local.example.toml` to `hfc.local.toml`, and update all the values
   in the file to match your desired configuration.
4. Run `./hfc build-deploy [stack]` to start deploying your first stack, **but
   see the directions below to complete the initial deployment**. To deploy
   additional stacks using the same build, run `./hfc deploy [stack]`.

While the initial deployment of the stack is in progress, you will need to
manually add DNS records to your domain to validate your ownership of it. To do
this:

1. Go to [AWS Certificate Manager][acm] in the AWS console, and look at the
   list of certificates for the us-east-1 region.
2. Find the pending certificate for the new stack, and click on its ID.
3. In the "Domains" box, look for the "CNAME name" and "CNAME value" columns.
4. Create a new CNAME record with your DNS host so that the "CNAME name"
   resolves to the "CNAME value."

Once you've created the CNAME record for validation, the certificate should be
issued and the stack deployment should proceed within a few minutes. Don't
delete the CNAME, as AWS will continue to look for it every time it renews the
certificate. Future stack deployments won't require any further interaction
from you.

[acm]: https://us-east-1.console.aws.amazon.com/acm/home?region=us-east-1#/certificates/list

After the initial deployment finishes, you will need to upload your TOML
configuration to the S3 bucket created by CloudFormation, then set up a CNAME
to the CloudFront domain with your DNS host. `hfc` will print the S3 path and
CloudFront domain after every successful stack deployment.

Note that the template can only be deployed to us-east-1, as CloudFront
requires the generated TLS certificate to be there. The `hfc` helper script
will automatically override the region for all AWS CLI commands.

If you don't wish to use the serverless AWS Lambda deployment, you can run
importbounce as a standard HTTP server by passing the `-http` flag with a
listening address (e.g. `-http 0.0.0.0:8080`).

## Future Work

* Mechanism so that accidentally deploying a bad config doesn't take down your
  whole domain.
* Enabling non-zero cache TTLs in CloudFront.
