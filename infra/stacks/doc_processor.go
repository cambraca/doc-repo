package stacks

//import (
//	"fmt"
//
//	"github.com/aws/aws-cdk-go/awscdk/v2"
//	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
//	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
//	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
//	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
//	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
//	"github.com/aws/aws-cdk-go/awscdk/v2/awsstepfunctions"
//	"github.com/aws/aws-cdk-go/awscdk/v2/awsstepfunctionstasks"
//	"github.com/aws/constructs-go/constructs/v10"
//	"github.com/aws/jsii-runtime-go"
//)
//
//type StepFunctionStackProps struct {
//	awscdk.NestedStackProps
//	Vpc                     awsec2.IVpc
//	DatabaseEndpointAddress *string
//	DatabaseEndpointPort    *string
//	DatabaseSecurityGroupId *string
//	ApiUrl                  *string
//}
//
//type StepFunctionStack struct {
//	awscdk.NestedStack
//	StepFunctionApiUrlOutput *string
//}
//
//func NewStepFunctionStack(scope constructs.Construct, id string, props *StepFunctionStackProps) *StepFunctionStack {
//	stack := awscdk.NewNestedStack(scope, &id, &props.NestedStackProps)
//
//	stepFunctionLogGroup := awslogs.NewLogGroup(stack, jsii.String("StepFunctionLogs"), &awslogs.LogGroupProps{
//		Retention:     awslogs.RetentionDays_ONE_WEEK,
//		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
//	})
//	rdsInteractionLambda := awslambda.NewFunction(stack, jsii.String("RDSInteractionLambda"), &awslambda.FunctionProps{
//		Runtime:        awslambda.Runtime_GO_1_X(),                                            // Or your preferred runtime
//		Handler:        jsii.String("main"),                                                   // Adjust handler
//		Code:           awslambda.Code_FromAsset(jsii.String("path/to/your/rds_lambda_code")), // Replace
//		Vpc:            props.Vpc,
//		SecurityGroups: &[]awsec2.ISecurityGroup{awsec2.SecurityGroup_FromSecurityGroupId(stack, jsii.String("DBAccessSG"), *props.DatabaseSecurityGroupId, &awsec2.SecurityGroupImportProps{Vpc: props.Vpc})},
//		Environment: &map[string]*string{
//			"DATABASE_HOST":     props.DatabaseEndpointAddress,
//			"DATABASE_PORT":     props.DatabaseEndpointPort,
//			"DATABASE_USER":     jsii.String("admin"),                                                       // Consider Secrets Manager
//			"DATABASE_PASSWORD": awscdk.SecretValue_PlainText(jsii.String("yourStrongPassword")).ToString(), // Consider Secrets Manager
//		},
//		LogGroup: stepFunctionLogGroup,
//	})
//	rdsInteractionLambda.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
//		Actions:    &[]*string{jsii.String("rds:Connect")},
//		Resources:  &[]*string{jsii.String(fmt.Sprintf("arn:aws:rds:*:%s:db:*", *awscdk.Stack_Of(stack).Account()))},
//		Conditions: &map[string]interface{}{"ArnEquals": map[string]*string{"rds:db-id": jsii.String("postgresdb")}}, // Adjust
//	}))
//	invokeRDSLambda := awsstepfunctionstasks.NewLambdaInvoke(stack, jsii.String("InvokeRDSLambda"), &awsstepfunctionstasks.LambdaInvokeProps{
//		LambdaFunction: rdsInteractionLambda,
//	})
//	finalState := awsstepfunctions.NewPass(stack, jsii.String("FinalState"))
//	definition := invokeRDSLambda.Next(finalState)
//	stateMachine := awsstepfunctions.NewStateMachine(stack, jsii.String("MyStateMachine"), &awsstepfunctions.StateMachineProps{
//		Definition: definition,
//		Logs: &awsstepfunctions.LogOptions{
//			Destination:          stepFunctionLogGroup,
//			IncludeExecutionData: jsii.Bool(true),
//			Level:                awsstepfunctions.LogLevel_ALL,
//		},
//	})
//
//	stepFunctionApi := awsapigateway.NewRestApi(
//		stack, jsii.String("StepFunctionApi"), &awsapigateway.RestApiProps{
//			RestApiName: jsii.String("StepFunctionTrigger"),
//			EndpointConfiguration: &awsapigateway.EndpointConfiguration{
//				Types: []awsapigateway.EndpointType{awsapigateway.EndpointType_REGIONAL},
//			},
//		})
//	stepFunctionIntegration := awsapigateway.NewIntegration(&awsapigateway.IntegrationProps{
//		Type:                  awsapigateway.IntegrationType_AWS,
//		IntegrationHttpMethod: jsii.String("POST"),
//		Uri: jsii.String(fmt.Sprintf("arn:aws:apigateway:%s:states:action/StartExecution",
//			*awscdk.Stack_Of(stack).Region())),
//		Options: &awsapigateway.IntegrationOptions{
//			CredentialsRole: awsiam.NewRole(stack, jsii.String("APIGatewayStepFunctionRole"), &awsiam.RoleProps{
//				AssumedBy: awsiam.NewServicePrincipal(jsii.String("apigateway.amazonaws.com"), nil),
//			}),
//			RequestTemplates: &map[string]*string{
//				"application/json": jsii.String(fmt.Sprintf(`{
//						"stateMachineArn": "%s",
//						"input": "$util.escapeJsonString($input.json('$'))"
//					}`, stateMachine.StateMachineArn())),
//			},
//			IntegrationResponses: []*awsapigateway.IntegrationResponse{
//				{StatusCode: jsii.String("200"), ResponseParameters: &map[string]*string{"method.response.header.Content-Type": "'application/json'"}},
//				{StatusCode: jsii.String("400"), SelectionPattern: jsii.String("4\\d{2}"), ResponseParameters: &map[string]*string{"method.response.header.Content-Type": "'application/json'"}},
//				{StatusCode: jsii.String("500"), SelectionPattern: jsii.String("5\\d{2}"), ResponseParameters: &map[string]*string{"method.response.header.Content-Type": "'application/json'"}},
//			},
//		},
//	})
//	stepFunctionMethod := stepFunctionApi.Root().AddMethod(jsii.String("POST"), stepFunctionIntegration, &awsapigateway.MethodOptions{
//		AuthorizationType: awsapigateway.AuthorizationType_IAM, // Secure with IAM roles
//		MethodResponses: []*awsapigateway.MethodResponse{
//			{StatusCode: jsii.String("200"), ResponseParameters: &map[string]*bool{"method.response.header.Content-Type": jsii.Bool(true)}},
//			{StatusCode: jsii.String("400"), ResponseParameters: &map[string]*bool{"method.response.header.Content-Type": jsii.Bool(true)}},
//			{StatusCode: jsii.String("500"), ResponseParameters: &map[string]*bool{"method.response.header.Content-Type": jsii.Bool(true)}},
//		},
//	})
//	stepFunctionIntegration.Options().CredentialsRole().AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
//		Actions:   &[]*string{jsii.String("states:StartExecution")},
//		Resources: &[]*string{stateMachine.StateMachineArn()},
//	}))
//
//	stepFunctionApiUrlOutput := awscdk.NewCfnOutput(stack, jsii.String("StepFunctionApiUrl"), &awscdk.CfnOutputProps{
//		ExportName: jsii.String("StepFunctionApiUrl"),
//		Value: stepFunctionApi.Url(),
//	})
//
//	return &StepFunctionStack{
//		NestedStack:              stack,
//		StepFunctionApiUrlOutput: stepFunctionApiUrlOutput.Value(),
//	}
//}
