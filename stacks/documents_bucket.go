package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type DocumentsBucketStackProps struct {
	awscdk.NestedStackProps
}

type DocumentsBucketStack struct {
	awscdk.NestedStack
	BucketArnOutput  *string // Exported field
	BucketNameOutput *string // Exported field
	bucket           awss3.IBucket
}

func NewDocumentsBucketStack(scope constructs.Construct, id string, props *DocumentsBucketStackProps) *DocumentsBucketStack {
	stack := awscdk.NewNestedStack(scope, &id, &props.NestedStackProps)

	documentsBucket := awss3.NewBucket(stack, jsii.String("DocumentsBucket"), &awss3.BucketProps{
		Versioned: jsii.Bool(true),
		//RemovalPolicy: awscdk.RemovalPolicy_RETAIN, // TODO: only do this for prod
	})

	bucketArnOutput := awscdk.NewCfnOutput(stack, jsii.String("DocumentsBucketArn"), &awscdk.CfnOutputProps{
		ExportName: jsii.String("DocumentsBucketArn"),
		Value:      documentsBucket.BucketArn(),
	})
	bucketNameOutput := awscdk.NewCfnOutput(stack, jsii.String("DocumentsBucketName"), &awscdk.CfnOutputProps{
		ExportName: jsii.String("DocumentsBucketName"),
		Value:      documentsBucket.BucketName(),
	})

	return &DocumentsBucketStack{
		NestedStack:      stack,
		BucketArnOutput:  bucketArnOutput.ImportValue(),  // Assign to exported field
		BucketNameOutput: bucketNameOutput.ImportValue(), // Assign to exported field
		bucket:           documentsBucket,
	}
}
