package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
	"infra/stacks"
	"log"
)

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	parentStack := awscdk.NewStack(app, jsii.String("DocRepo-dev"), &awscdk.StackProps{
		Env: env(),
	})

	vpcStack := stacks.NewVpcStack(parentStack, "VpcStack", &stacks.VpcStackProps{
		NestedStackProps: awscdk.NestedStackProps{
			Description: jsii.String("Network infrastructure for the application."),
		},
	})

	documentsBucketStack := stacks.NewDocumentsBucketStack(parentStack, "DocumentsBucketStack", &stacks.DocumentsBucketStackProps{
		NestedStackProps: awscdk.NestedStackProps{
			Description: jsii.String("S3 bucket to store documents."),
		},
	})

	//databaseStack := stacks.NewDatabaseStack(parentStack, "DatabaseStack", &stacks.DatabaseStackProps{
	//	NestedStackProps: awscdk.NestedStackProps{
	//		Description: jsii.String("Database for the application backend."),
	//	},
	//	Vpc: vpcStack.Vpc,
	//})

	apiStack := stacks.NewApiStack(parentStack, "ApiStack", &stacks.ApiStackProps{
		NestedStackProps: awscdk.NestedStackProps{
			Description: jsii.String("Backend API service running on ECS."),
		},
		Vpc:             vpcStack.Vpc,
		DocumentsBucket: documentsBucketStack.Bucket,
		//DatabaseEndpointAddress: databaseStack.EndpointAddressOutput,
		//DatabaseEndpointPort:    databaseStack.EndpointPortOutput,
		//DatabaseSecurityGroupId: databaseStack.SecurityGroupIdOutput,
	})

	log.Print("Stacks", vpcStack, documentsBucketStack, apiStack)
	//
	//stepFunctionStack := stacks.NewStepFunctionStack(parentStack, "StepFunctionStack", &stacks.StepFunctionStackProps{
	//	NestedStackProps: awscdk.NestedStackProps{
	//		Description: jsii.String("Step Function workflow for asynchronous tasks."),
	//	},
	//	Vpc:                     vpcStack.Vpc,
	//	DatabaseEndpointAddress: databaseStack.EndpointAddressOutput,
	//	DatabaseEndpointPort:    databaseStack.EndpointPortOutput,
	//	DatabaseSecurityGroupId: databaseStack.SecurityGroupIdOutput,
	//	ApiUrl:                  apiStack.ApiUrlOutput,
	//})
	//
	//frontendStack := stacks.NewFrontendStack(parentStack, "FrontendStack", &stacks.FrontendStackProps{
	//	NestedStackProps: awscdk.NestedStackProps{
	//		Description: jsii.String("Frontend web application running on ECS."),
	//	},
	//	Vpc:    vpcStack.Vpc,
	//	ApiUrl: apiStack.ApiUrlOutput,
	//})

	app.Synth(nil)
}

// Environment configuration
func env() *awscdk.Environment {
	// TODO: make this configurable
	return &awscdk.Environment{
		Account: jsii.String("434274429363"),
		Region:  jsii.String("us-east-2"),
	}
}
