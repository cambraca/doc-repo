package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type VpcStackProps struct {
	awscdk.NestedStackProps
}

type VpcStack struct {
	awscdk.NestedStack
	Vpc ec2.IVpc // Exported field
}

func NewVpcStack(scope constructs.Construct, id string, props *VpcStackProps) *VpcStack {
	stack := awscdk.NewNestedStack(scope, &id, &props.NestedStackProps)

	vpc := ec2.NewVpc(stack, jsii.String("VPC"), &ec2.VpcProps{
		MaxAzs: jsii.Number(2),
	})

	return &VpcStack{
		NestedStack: stack,
		Vpc:         vpc, // Assign to the exported field
	}
}
