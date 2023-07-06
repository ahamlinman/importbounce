// Code generated by smithy-go-codegen DO NOT EDIT.

package s3

import (
	"context"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	s3cust "github.com/aws/aws-sdk-go-v2/service/s3/internal/customizations"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"time"
)

// Retrieves all the metadata from an object without returning the object itself.
// This action is useful if you're interested only in an object's metadata. To use
// GetObjectAttributes , you must have READ access to the object.
// GetObjectAttributes combines the functionality of HeadObject and ListParts . All
// of the data returned with each of those individual calls can be returned with a
// single call to GetObjectAttributes . If you encrypt an object by using
// server-side encryption with customer-provided encryption keys (SSE-C) when you
// store the object in Amazon S3, then when you retrieve the metadata from the
// object, you must use the following headers:
//   - x-amz-server-side-encryption-customer-algorithm
//   - x-amz-server-side-encryption-customer-key
//   - x-amz-server-side-encryption-customer-key-MD5
//
// For more information about SSE-C, see Server-Side Encryption (Using
// Customer-Provided Encryption Keys) (https://docs.aws.amazon.com/AmazonS3/latest/dev/ServerSideEncryptionCustomerKeys.html)
// in the Amazon S3 User Guide.
//   - Encryption request headers, such as x-amz-server-side-encryption , should
//     not be sent for GET requests if your object uses server-side encryption with
//     Amazon Web Services KMS keys stored in Amazon Web Services Key Management
//     Service (SSE-KMS) or server-side encryption with Amazon S3 managed keys
//     (SSE-S3). If your object does use these types of keys, you'll get an HTTP 400
//     Bad Request error.
//   - The last modified property in this case is the creation date of the object.
//
// Consider the following when using request headers:
//   - If both of the If-Match and If-Unmodified-Since headers are present in the
//     request as follows, then Amazon S3 returns the HTTP status code 200 OK and the
//     data requested:
//   - If-Match condition evaluates to true .
//   - If-Unmodified-Since condition evaluates to false .
//   - If both of the If-None-Match and If-Modified-Since headers are present in
//     the request as follows, then Amazon S3 returns the HTTP status code 304 Not
//     Modified :
//   - If-None-Match condition evaluates to false .
//   - If-Modified-Since condition evaluates to true .
//
// For more information about conditional requests, see RFC 7232 (https://tools.ietf.org/html/rfc7232)
// . Permissions The permissions that you need to use this operation depend on
// whether the bucket is versioned. If the bucket is versioned, you need both the
// s3:GetObjectVersion and s3:GetObjectVersionAttributes permissions for this
// operation. If the bucket is not versioned, you need the s3:GetObject and
// s3:GetObjectAttributes permissions. For more information, see Specifying
// Permissions in a Policy (https://docs.aws.amazon.com/AmazonS3/latest/dev/using-with-s3-actions.html)
// in the Amazon S3 User Guide. If the object that you request does not exist, the
// error Amazon S3 returns depends on whether you also have the s3:ListBucket
// permission.
//   - If you have the s3:ListBucket permission on the bucket, Amazon S3 returns an
//     HTTP status code 404 Not Found ("no such key") error.
//   - If you don't have the s3:ListBucket permission, Amazon S3 returns an HTTP
//     status code 403 Forbidden ("access denied") error.
//
// The following actions are related to GetObjectAttributes :
//   - GetObject (https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObject.html)
//   - GetObjectAcl (https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObjectAcl.html)
//   - GetObjectLegalHold (https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObjectLegalHold.html)
//   - GetObjectLockConfiguration (https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObjectLockConfiguration.html)
//   - GetObjectRetention (https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObjectRetention.html)
//   - GetObjectTagging (https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObjectTagging.html)
//   - HeadObject (https://docs.aws.amazon.com/AmazonS3/latest/API/API_HeadObject.html)
//   - ListParts (https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListParts.html)
func (c *Client) GetObjectAttributes(ctx context.Context, params *GetObjectAttributesInput, optFns ...func(*Options)) (*GetObjectAttributesOutput, error) {
	if params == nil {
		params = &GetObjectAttributesInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "GetObjectAttributes", params, optFns, c.addOperationGetObjectAttributesMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*GetObjectAttributesOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type GetObjectAttributesInput struct {

	// The name of the bucket that contains the object. When using this action with an
	// access point, you must direct requests to the access point hostname. The access
	// point hostname takes the form
	// AccessPointName-AccountId.s3-accesspoint.Region.amazonaws.com. When using this
	// action with an access point through the Amazon Web Services SDKs, you provide
	// the access point ARN in place of the bucket name. For more information about
	// access point ARNs, see Using access points (https://docs.aws.amazon.com/AmazonS3/latest/userguide/using-access-points.html)
	// in the Amazon S3 User Guide. When you use this action with Amazon S3 on
	// Outposts, you must direct requests to the S3 on Outposts hostname. The S3 on
	// Outposts hostname takes the form
	// AccessPointName-AccountId.outpostID.s3-outposts.Region.amazonaws.com . When you
	// use this action with S3 on Outposts through the Amazon Web Services SDKs, you
	// provide the Outposts access point ARN in place of the bucket name. For more
	// information about S3 on Outposts ARNs, see What is S3 on Outposts? (https://docs.aws.amazon.com/AmazonS3/latest/userguide/S3onOutposts.html)
	// in the Amazon S3 User Guide.
	//
	// This member is required.
	Bucket *string

	// The object key.
	//
	// This member is required.
	Key *string

	// Specifies the fields at the root level that you want returned in the response.
	// Fields that you do not specify are not returned.
	//
	// This member is required.
	ObjectAttributes []types.ObjectAttributes

	// The account ID of the expected bucket owner. If the bucket is owned by a
	// different account, the request fails with the HTTP status code 403 Forbidden
	// (access denied).
	ExpectedBucketOwner *string

	// Sets the maximum number of parts to return.
	MaxParts int32

	// Specifies the part after which listing should begin. Only parts with higher
	// part numbers will be listed.
	PartNumberMarker *string

	// Confirms that the requester knows that they will be charged for the request.
	// Bucket owners need not specify this parameter in their requests. For information
	// about downloading objects from Requester Pays buckets, see Downloading Objects
	// in Requester Pays Buckets (https://docs.aws.amazon.com/AmazonS3/latest/dev/ObjectsinRequesterPaysBuckets.html)
	// in the Amazon S3 User Guide.
	RequestPayer types.RequestPayer

	// Specifies the algorithm to use when encrypting the object (for example, AES256).
	SSECustomerAlgorithm *string

	// Specifies the customer-provided encryption key for Amazon S3 to use in
	// encrypting data. This value is used to store the object and then it is
	// discarded; Amazon S3 does not store the encryption key. The key must be
	// appropriate for use with the algorithm specified in the
	// x-amz-server-side-encryption-customer-algorithm header.
	SSECustomerKey *string

	// Specifies the 128-bit MD5 digest of the encryption key according to RFC 1321.
	// Amazon S3 uses this header for a message integrity check to ensure that the
	// encryption key was transmitted without error.
	SSECustomerKeyMD5 *string

	// The version ID used to reference a specific version of the object.
	VersionId *string

	noSmithyDocumentSerde
}

type GetObjectAttributesOutput struct {

	// The checksum or digest of the object.
	Checksum *types.Checksum

	// Specifies whether the object retrieved was ( true ) or was not ( false ) a
	// delete marker. If false , this response header does not appear in the response.
	DeleteMarker bool

	// An ETag is an opaque identifier assigned by a web server to a specific version
	// of a resource found at a URL.
	ETag *string

	// The creation date of the object.
	LastModified *time.Time

	// A collection of parts associated with a multipart upload.
	ObjectParts *types.GetObjectAttributesParts

	// The size of the object in bytes.
	ObjectSize int64

	// If present, indicates that the requester was successfully charged for the
	// request.
	RequestCharged types.RequestCharged

	// Provides the storage class information of the object. Amazon S3 returns this
	// header for all objects except for S3 Standard storage class objects. For more
	// information, see Storage Classes (https://docs.aws.amazon.com/AmazonS3/latest/dev/storage-class-intro.html)
	// .
	StorageClass types.StorageClass

	// The version ID of the object.
	VersionId *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationGetObjectAttributesMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsRestxml_serializeOpGetObjectAttributes{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsRestxml_deserializeOpGetObjectAttributes{}, middleware.After)
	if err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddClientRequestIDMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddComputeContentLengthMiddleware(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = v4.AddComputePayloadSHA256Middleware(stack); err != nil {
		return err
	}
	if err = addRetryMiddlewares(stack, options); err != nil {
		return err
	}
	if err = addHTTPSignerV4Middleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = swapWithCustomHTTPSignerMiddleware(stack, options); err != nil {
		return err
	}
	if err = addOpGetObjectAttributesValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opGetObjectAttributes(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addMetadataRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecursionDetection(stack); err != nil {
		return err
	}
	if err = addGetObjectAttributesUpdateEndpoint(stack, options); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = v4.AddContentSHA256HeaderMiddleware(stack); err != nil {
		return err
	}
	if err = disableAcceptEncodingGzip(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opGetObjectAttributes(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "s3",
		OperationName: "GetObjectAttributes",
	}
}

// getGetObjectAttributesBucketMember returns a pointer to string denoting a
// provided bucket member valueand a boolean indicating if the input has a modeled
// bucket name,
func getGetObjectAttributesBucketMember(input interface{}) (*string, bool) {
	in := input.(*GetObjectAttributesInput)
	if in.Bucket == nil {
		return nil, false
	}
	return in.Bucket, true
}
func addGetObjectAttributesUpdateEndpoint(stack *middleware.Stack, options Options) error {
	return s3cust.UpdateEndpoint(stack, s3cust.UpdateEndpointOptions{
		Accessor: s3cust.UpdateEndpointParameterAccessor{
			GetBucketFromInput: getGetObjectAttributesBucketMember,
		},
		UsePathStyle:                   options.UsePathStyle,
		UseAccelerate:                  options.UseAccelerate,
		SupportsAccelerate:             true,
		TargetS3ObjectLambda:           false,
		EndpointResolver:               options.EndpointResolver,
		EndpointResolverOptions:        options.EndpointOptions,
		UseARNRegion:                   options.UseARNRegion,
		DisableMultiRegionAccessPoints: options.DisableMultiRegionAccessPoints,
	})
}
