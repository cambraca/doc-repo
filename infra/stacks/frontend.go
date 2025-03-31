package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfront"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfrontorigins"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3deployment"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"path/filepath"
)

type FrontendStackProps struct {
	awscdk.NestedStackProps
	ApiUrlOutput *string
}

type FrontendStack struct {
	awscdk.NestedStack
	CloudFrontUrlOutput *string
}

func NewFrontendStack(scope constructs.Construct, id string, props *FrontendStackProps) *FrontendStack {
	stack := awscdk.NewNestedStack(scope, &id, &props.NestedStackProps)

	// 1. Create an S3 bucket to host the frontend files
	bucket := awss3.NewBucket(stack, jsii.String("FrontendBucket"), &awss3.BucketProps{
		PublicReadAccess:  jsii.Bool(false),             // CloudFront will handle access
		RemovalPolicy:     awscdk.RemovalPolicy_DESTROY, // Adjust as needed
		AutoDeleteObjects: jsii.Bool(true),              // Adjust as needed
	})

	// 2. Deploy the frontend files to the S3 bucket
	deployment := awss3deployment.NewBucketDeployment(stack, jsii.String("DeployFrontend"), &awss3deployment.BucketDeploymentProps{
		Sources: &[]awss3deployment.Source{
			awss3deployment.Source_Asset(jsii.String(filepath.Join("..", "frontend", "dist"))), // Assuming your Ember build output is in 'frontend/dist'
		},
		DestinationBucket: bucket,
	})

	// 3. Create a CloudFront distribution to serve the S3 bucket content
	distribution := awscloudfront.NewDistribution(stack, jsii.String("FrontendDistribution"), &awscloudfront.DistributionProps{
		DefaultRootObject: jsii.String("index.html"), // Your main HTML file
		ErrorResponses: &[]*awscloudfront.ErrorResponse{
			{
				HttpStatus:         jsii.Number(403),
				ResponseHttpStatus: jsii.Number(200),
				ResponsePagePath:   jsii.String("/index.html"),
			},
			{
				HttpStatus:         jsii.Number(404),
				ResponseHttpStatus: jsii.Number(200),
				ResponsePagePath:   jsii.String("/index.html"),
			},
		},
		DefaultBehavior: &awscloudfront.BehaviorOptions{
			Origin:               awscloudfrontorigins.NewS3Origin(bucket, &awscloudfrontorigins.S3OriginProps{}),
			ViewerProtocolPolicy: awscloudfront.ViewerProtocolPolicy_REDIRECT_TO_HTTPS,
		},
	})

	// 4. Output the CloudFront distribution URL
	cloudFrontUrlOutput := awscdk.NewCfnOutput(stack, jsii.String("FrontendUrl"), &awscdk.CfnOutputProps{
		Value: distribution.DomainName(),
	})

	// 5. Output the API URL for the frontend to use (configuration)
	apiUrlConfigOutput := awscdk.NewCfnOutput(stack, jsii.String("ApiUrlConfig"), &awscdk.CfnOutputProps{
		Value: props.ApiUrlOutput,
	})

	return &FrontendStack{
		NestedStack:         stack,
		CloudFrontUrlOutput: cloudFrontUrlOutput.Value(),
	}
}
