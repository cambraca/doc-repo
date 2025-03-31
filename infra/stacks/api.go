package stacks

import (
	"fmt"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecrassets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecspatterns"
	elbv2 "github.com/aws/aws-cdk-go/awscdk/v2/awselasticloadbalancingv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"path/filepath"
)

type ApiStackProps struct {
	awscdk.NestedStackProps
	Vpc             awsec2.IVpc
	DocumentsBucket awss3.IBucket
}

type ApiStack struct {
	awscdk.NestedStack
	ApiUrlOutput awscdk.CfnOutput
}

func NewApiStack(scope constructs.Construct, id string, props *ApiStackProps) *ApiStack {
	stack := awscdk.NewNestedStack(scope, &id, &props.NestedStackProps)

	apiCluster := awsecs.NewCluster(stack, jsii.String("ApiCluster"), &awsecs.ClusterProps{
		Vpc: props.Vpc,
	})
	apiLogGroup := awslogs.NewLogGroup(stack, jsii.String("ApiLogGroup"), &awslogs.LogGroupProps{
		Retention:     awslogs.RetentionDays_ONE_WEEK,
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})

	// Build Docker image from local Dockerfile
	dockerImageAsset := awsecrassets.NewDockerImageAsset(stack, jsii.String("GoApiImage"), &awsecrassets.DockerImageAssetProps{
		Directory: jsii.String(filepath.Join("..", "api")),
	})

	apiTaskDefinition := awsecs.NewFargateTaskDefinition(stack, jsii.String("ApiTaskDefinition"), &awsecs.FargateTaskDefinitionProps{
		Cpu:            jsii.Number(256),
		MemoryLimitMiB: jsii.Number(512),
		ExecutionRole: awsiam.NewRole(stack, jsii.String("ApiTaskExecutionRole"), &awsiam.RoleProps{
			AssumedBy: awsiam.NewServicePrincipal(jsii.String("ecs-tasks.amazonaws.com"), nil),
		}),
		TaskRole: awsiam.NewRole(stack, jsii.String("ApiTaskRole"), &awsiam.RoleProps{
			AssumedBy: awsiam.NewServicePrincipal(jsii.String("ecs-tasks.amazonaws.com"), nil),
		}),
	})

	// Grant ECR pull permission to the EXECUTION ROLE
	apiTaskDefinition.ExecutionRole().AddToPrincipalPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   &[]*string{jsii.String("ecr:GetAuthorizationToken")},
		Resources: &[]*string{jsii.String("*")},
	}))
	apiTaskDefinition.ExecutionRole().AddToPrincipalPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   &[]*string{jsii.String("ecr:BatchCheckLayerAvailability"), jsii.String("ecr:GetDownloadUrlForLayer"), jsii.String("ecr:BatchGetImage")},
		Resources: &[]*string{dockerImageAsset.Repository().RepositoryArn()}, // Limit to the specific repository
	}))

	apiTaskDefinition.AddContainer(jsii.String("ApiContainer"), &awsecs.ContainerDefinitionOptions{
		Image: awsecs.ContainerImage_FromRegistry(dockerImageAsset.ImageUri(), nil),
		PortMappings: &[]*awsecs.PortMapping{
			{
				ContainerPort: jsii.Number(8080),
			},
		},
		Environment: &map[string]*string{
			"DOCUMENTS_BUCKET_NAME": props.DocumentsBucket.BucketName(),
			//"DATABASE_HOST":         props.DatabaseEndpointAddress,
			//"DATABASE_PORT":         props.DatabaseEndpointPort,
			//"DATABASE_USER":     jsii.String("admin"),                                                       // Consider Secrets Manager
			//"DATABASE_PASSWORD": awscdk.SecretValue_PlainText(jsii.String("yourStrongPassword")).ToString(), // Consider Secrets Manager
		},
		Logging: awsecs.LogDriver_AwsLogs(&awsecs.AwsLogDriverProps{
			LogGroup:     apiLogGroup,
			StreamPrefix: jsii.String("api"),
		}),
	})
	apiTaskDefinition.TaskRole().AddToPrincipalPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions: &[]*string{
			jsii.String("s3:GetObject"),
			jsii.String("s3:PutObject"),
		},
		Resources: &[]*string{jsii.String(fmt.Sprintf("%s/*", *props.DocumentsBucket.BucketArn()))},
	}))
	//apiTaskDefinition.TaskRole().AddToPrincipalPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
	//	Actions:   &[]*string{jsii.String("rds:Connect")},
	//	Resources: &[]*string{jsii.String(fmt.Sprintf("arn:aws:rds:*:%s:db:*", *awscdk.Stack_Of(stack).Account()))},
	//	Conditions: &map[string]interface{}{"ArnEquals": map[string]*string{"rds:db-id": jsii.String("postgresdb")}},
	//}))
	apiService := awsecspatterns.NewApplicationLoadBalancedFargateService(stack, jsii.String("ApiService"), &awsecspatterns.ApplicationLoadBalancedFargateServiceProps{
		Cluster:            apiCluster,
		TaskDefinition:     apiTaskDefinition,
		PublicLoadBalancer: jsii.Bool(true),
		DesiredCount:       jsii.Number(1), // How many instances of the api we want to run
		ListenerPort:       jsii.Number(80),
		// ContainerPort:      jsii.Number(8080), // Removed here
	})

	// Configure the health check for the target group
	apiService.TargetGroup().ConfigureHealthCheck(&elbv2.HealthCheck{
		Path:     jsii.String("/status"),
		Port:     jsii.String("8080"),
		Protocol: elbv2.Protocol_HTTP,
		//HealthyHttpCodes: jsii.String("200"),
	})

	//apiService.Service().Connections().AllowTo(
	//	awsec2.Peer_SecurityGroupId(props.DatabaseSecurityGroupId, jsii.String("RDS Access")),
	//	awsec2.Port_Tcp(jsii.Number(5432)),
	//	jsii.String("Allow API access to PostgreSQL"),
	//)

	apiUrlOutput := awscdk.NewCfnOutput(stack, jsii.String("ApiUrl"), &awscdk.CfnOutputProps{
		Value: apiService.LoadBalancer().LoadBalancerDnsName(),
	})

	return &ApiStack{
		NestedStack:  stack,
		ApiUrlOutput: apiUrlOutput,
	}
}
