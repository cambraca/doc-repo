package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
	"infra/stacks"
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

	databaseStack := stacks.NewDatabaseStack(parentStack, "DatabaseStack", &stacks.DatabaseStackProps{
		NestedStackProps: awscdk.NestedStackProps{
			Description: jsii.String("Main application database."),
		},
		Vpc: vpcStack.Vpc,

		// TODO: configure this by env
		AllowAccessFromEverywhere: true,
	})

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

	webAppStack := stacks.NewWebAppStack(parentStack, "WebApp", &stacks.WebAppStackProps{
		NestedStackProps: awscdk.NestedStackProps{
			Description: jsii.String("Web app with API service running on ECS and frontend as an S3 deployment."),
		},
		Vpc:             vpcStack.Vpc,
		DocumentsBucket: documentsBucketStack.Bucket,
		DbInstance:      databaseStack.Instance,
		DbSecurityGroup: databaseStack.SecurityGroup,
	})

	awscdk.NewCfnOutput(parentStack, jsii.String("Url"), &awscdk.CfnOutputProps{
		ExportName: jsii.String("Url"),
		Key:        jsii.String("Url"),
		Value:      webAppStack.FrontendUrlOutput.ImportValue(),
	})

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
