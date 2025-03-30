package stacks

//import (
//	"github.com/aws/aws-cdk-go/awscdk/v2"
//	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
//	"github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
//	"github.com/aws/aws-cdk-go/awscdk/v2/awsecspatterns"
//	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
//	"github.com/aws/constructs-go/constructs/v10"
//	"github.com/aws/jsii-runtime-go"
//)
//
//type FrontendStackProps struct {
//	awscdk.NestedStackProps
//	Vpc    awsec2.IVpc
//	ApiUrl *string
//}
//
//type FrontendStack struct {
//	awscdk.NestedStack
//	FrontendUrlOutput *string
//}
//
//func NewFrontendStack(scope constructs.Construct, id string, props *FrontendStackProps) *FrontendStack {
//	stack := awscdk.NewNestedStack(scope, &id, &props.NestedStackProps)
//
//	frontendCluster := awsecs.NewCluster(stack, jsii.String("FrontendCluster"), &awsecs.ClusterProps{
//		Vpc: props.Vpc,
//	})
//	frontendLogGroup := awslogs.NewLogGroup(stack, jsii.String("FrontendLogGroup"), &awslogs.LogGroupProps{
//		Retention:     awslogs.RetentionDays_ONE_WEEK,
//		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
//	})
//	frontendTaskDefinition := awsecs.NewFargateTaskDefinition(stack, jsii.String("FrontendTaskDefinition"), &awsecs.FargateTaskDefinitionProps{
//		Cpu:            jsii.Number(256),
//		MemoryLimitMiB: jsii.Number(512),
//		LogDriver: awsecs.LogDriver_AwsLogs(&awsecs.AwsLogDriverProps{
//			LogGroup:     frontendLogGroup,
//			StreamPrefix: jsii.String("frontend"),
//		}),
//	})
//	frontendContainer := frontendTaskDefinition.AddContainer(jsii.String("FrontendContainer"), &awsecs.ContainerDefinitionOptions{
//		Image: awsecs.ContainerImage_FromRegistry(jsii.String("your-javascript-app-image:latest")), // Replace
//		PortMappings: &[]*awsecs.PortMapping{
//			{
//				ContainerPort: jsii.Number(3000),
//			},
//		},
//		Environment: &map[string]*string{
//			"API_URL": props.ApiUrl,
//		},
//	})
//	frontendService := awsecspatterns.NewApplicationLoadBalancedFargateService(stack, jsii.String("FrontendService"), &awsecspatterns.ApplicationLoadBalancedFargateServiceProps{
//		Cluster:            frontendCluster,
//		TaskDefinition:     frontendTaskDefinition,
//		PublicLoadBalancer: jsii.Bool(true),
//		DesiredCount:       jsii.Number(1), // Adjust
//		ListenerPort:       jsii.Number(80),
//		ContainerPort:      jsii.Number(3000),
//	})
//
//	frontendUrlOutput := awscdk.NewCfnOutput(stack, jsii.String("FrontendUrl"), &awscdk.CfnOutputProps{
//		Value: frontendService.LoadBalancer().LoadBalancerDnsName(),
//	})
//
//	return &FrontendStack{
//		NestedStack:       stack,
//		FrontendUrlOutput: frontendUrlOutput.Value(),
//	}
//}
