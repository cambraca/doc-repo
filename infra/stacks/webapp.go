package stacks

import (
	"fmt"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	cloudfront "github.com/aws/aws-cdk-go/awscdk/v2/awscloudfront"
	cforigins "github.com/aws/aws-cdk-go/awscdk/v2/awscloudfrontorigins"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	ecrassets "github.com/aws/aws-cdk-go/awscdk/v2/awsecrassets"
	ecs "github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	ecspatterns "github.com/aws/aws-cdk-go/awscdk/v2/awsecspatterns"
	elbv2 "github.com/aws/aws-cdk-go/awscdk/v2/awselasticloadbalancingv2"
	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	logs "github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	rds "github.com/aws/aws-cdk-go/awscdk/v2/awsrds"
	s3 "github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	s3deployment "github.com/aws/aws-cdk-go/awscdk/v2/awss3deployment"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"path/filepath"
)

type WebAppStackProps struct {
	awscdk.NestedStackProps
	Vpc             ec2.IVpc
	DocumentsBucket s3.IBucket
	DbInstance      rds.IDatabaseInstance
	DbSecurityGroup ec2.ISecurityGroup
}

type WebAppStack struct {
	awscdk.NestedStack
	ApiUrlOutput      awscdk.CfnOutput
	FrontendUrlOutput awscdk.CfnOutput
}

func NewWebAppStack(scope constructs.Construct, id string, props *WebAppStackProps) *WebAppStack {
	stack := awscdk.NewNestedStack(scope, &id, &props.NestedStackProps)

	// API

	apiCluster := ecs.NewCluster(stack, jsii.String("ApiCluster"), &ecs.ClusterProps{
		Vpc: props.Vpc,
	})
	apiLogGroup := logs.NewLogGroup(stack, jsii.String("ApiLogGroup"), &logs.LogGroupProps{
		Retention:     logs.RetentionDays_ONE_WEEK,
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})

	// Build Docker image from local Dockerfile
	dockerImageAsset := ecrassets.NewDockerImageAsset(stack, jsii.String("GoApiImage"), &ecrassets.DockerImageAssetProps{
		Directory: jsii.String(filepath.Join("..", "api")),
	})

	apiTaskDefinition := ecs.NewFargateTaskDefinition(stack, jsii.String("ApiTaskDefinition"), &ecs.FargateTaskDefinitionProps{
		Cpu:            jsii.Number(256),
		MemoryLimitMiB: jsii.Number(512),
		ExecutionRole: iam.NewRole(stack, jsii.String("ApiTaskExecutionRole"), &iam.RoleProps{
			AssumedBy: iam.NewServicePrincipal(jsii.String("ecs-tasks.amazonaws.com"), nil),
		}),
		TaskRole: iam.NewRole(stack, jsii.String("ApiTaskRole"), &iam.RoleProps{
			AssumedBy: iam.NewServicePrincipal(jsii.String("ecs-tasks.amazonaws.com"), nil),
		}),
	})

	// Grant ECR pull permission to the EXECUTION ROLE
	apiTaskDefinition.ExecutionRole().AddToPrincipalPolicy(iam.NewPolicyStatement(&iam.PolicyStatementProps{
		Actions:   &[]*string{jsii.String("ecr:GetAuthorizationToken")},
		Resources: &[]*string{jsii.String("*")},
	}))
	apiTaskDefinition.ExecutionRole().AddToPrincipalPolicy(iam.NewPolicyStatement(&iam.PolicyStatementProps{
		Actions:   &[]*string{jsii.String("ecr:BatchCheckLayerAvailability"), jsii.String("ecr:GetDownloadUrlForLayer"), jsii.String("ecr:BatchGetImage")},
		Resources: &[]*string{dockerImageAsset.Repository().RepositoryArn()}, // Limit to the specific repository
	}))

	apiTaskDefinition.AddContainer(jsii.String("ApiContainer"), &ecs.ContainerDefinitionOptions{
		Image: ecs.ContainerImage_FromRegistry(dockerImageAsset.ImageUri(), nil),
		PortMappings: &[]*ecs.PortMapping{
			{
				ContainerPort: jsii.Number(8080),
			},
		},
		Environment: &map[string]*string{
			"DOCUMENTS_BUCKET_NAME": props.DocumentsBucket.BucketName(),
			"DB_HOST":               props.DbInstance.DbInstanceEndpointAddress(),
			"DB_PORT":               props.DbInstance.DbInstanceEndpointPort(),
			"DB_USER":               jsii.String("docrepouser"), // TODO: get from db stack
			"DB_PASSWORD":           jsii.String("abcd1234"),    // TODO: obviously
			"DB_NAME":               jsii.String("docrepo"),     // TODO: get from db stack
			"DB_SSLMODE":            jsii.String("require"),
			//"DATABASE_HOST":         props.DatabaseEndpointAddress,
			//"DATABASE_PORT":         props.DatabaseEndpointPort,
			//"DATABASE_USER":     jsii.String("admin"),                                                       // Consider Secrets Manager
			//"DATABASE_PASSWORD": awscdk.SecretValue_PlainText(jsii.String("yourStrongPassword")).ToString(), // Consider Secrets Manager
		},
		Logging: ecs.LogDriver_AwsLogs(&ecs.AwsLogDriverProps{
			LogGroup:     apiLogGroup,
			StreamPrefix: jsii.String("api"),
		}),
	})
	apiTaskDefinition.TaskRole().AddToPrincipalPolicy(iam.NewPolicyStatement(&iam.PolicyStatementProps{
		Actions: &[]*string{
			jsii.String("s3:GetObject"),
			jsii.String("s3:PutObject"),
		},
		Resources: &[]*string{jsii.String(fmt.Sprintf("%s/*", *props.DocumentsBucket.BucketArn()))},
	}))
	apiTaskDefinition.TaskRole().AddToPrincipalPolicy(iam.NewPolicyStatement(&iam.PolicyStatementProps{
		Actions:   &[]*string{jsii.String("rds:Connect")},
		Resources: &[]*string{props.DbInstance.InstanceArn()},
		//Resources: &[]*string{jsii.String(fmt.Sprintf("arn:aws:rds:*:%s:db:*", *awscdk.Stack_Of(stack).Account()))},
		//Conditions: &map[string]interface{}{
		//	"ArnEquals": map[string]*string{"rds:db-id": jsii.String("postgresdb")},
		//},
	}))

	apiService := ecspatterns.NewApplicationLoadBalancedFargateService(stack, jsii.String("ApiService"), &ecspatterns.ApplicationLoadBalancedFargateServiceProps{
		Cluster:            apiCluster,
		TaskDefinition:     apiTaskDefinition,
		PublicLoadBalancer: jsii.Bool(true), // TODO: restrict so only CloudFront can access this
		DesiredCount:       jsii.Number(1),  // How many instances of the api we want to run
		//EphemeralStorageGiB: jsii.Number(21), // GiB (min: 21)
	})

	// Configure the health check for the target group
	apiService.TargetGroup().ConfigureHealthCheck(&elbv2.HealthCheck{
		Path:     jsii.String("/status"),
		Port:     jsii.String("8080"),
		Protocol: elbv2.Protocol_HTTP,

		//// TODO: remove these
		//// (the default values should be fine; we're lowering them to try to make deploys faster in dev)
		//HealthyThresholdCount: jsii.Number(2),
		//Interval:              awscdk.Duration_Seconds(jsii.Number(5)),
		//Timeout:               awscdk.Duration_Seconds(jsii.Number(4)),
	})

	apiService.Service().Connections().AllowTo(
		ec2.Peer_SecurityGroupId(props.DbSecurityGroup.SecurityGroupId(), nil), // jsii.String("RDSSecurityGroup")), TODO: ?
		//ec2.Port_Tcp(props.DbInstance.DbInstanceEndpointPort()),
		ec2.Port_Tcp(jsii.Number(5432)),
		jsii.String("Allow API access to PostgreSQL"),
	)

	// --- CloudFront for API ---
	apiOrigin := cforigins.NewLoadBalancerV2Origin(apiService.LoadBalancer(), &cforigins.LoadBalancerV2OriginProps{
		ProtocolPolicy: cloudfront.OriginProtocolPolicy_HTTP_ONLY, // CloudFront talks to ALB over HTTP
	})

	apiDistribution := cloudfront.NewDistribution(stack, jsii.String("ApiDistribution"), &cloudfront.DistributionProps{
		DefaultBehavior: &cloudfront.BehaviorOptions{
			Origin:               apiOrigin,
			ViewerProtocolPolicy: cloudfront.ViewerProtocolPolicy_REDIRECT_TO_HTTPS,
			//AllowedMethods:       cloudfront.AllowedMethods_ALLOW_ALL(),                              // Or specific methods
			CachePolicy: cloudfront.CachePolicy_CACHING_DISABLED(), // API responses are often dynamic
			//OriginRequestPolicy:  cloudfront.OriginRequestPolicy_ALL_VIEWER_AND_CLOUDFRONT_2022(), // Forward headers
		},
		// You can configure error responses, logging, etc. here
	})

	apiUrl := fmt.Sprintf("https://%s", *apiDistribution.DomainName())

	apiUrlOutput := awscdk.NewCfnOutput(stack, jsii.String("ApiUrl"), &awscdk.CfnOutputProps{
		ExportName: jsii.String("ApiUrl"),
		Value:      jsii.String(apiUrl),
	})

	// FRONTEND

	// S3 bucket to host the frontend files
	frontendBucket := s3.NewBucket(stack, jsii.String("FrontendBucket"), &s3.BucketProps{
		PublicReadAccess:  jsii.Bool(false),
		RemovalPolicy:     awscdk.RemovalPolicy_DESTROY,
		AutoDeleteObjects: jsii.Bool(true),
	})

	// CloudFront distribution to serve the S3 bucket content
	s3Origin := cforigins.S3BucketOrigin_WithOriginAccessControl(frontendBucket, &cforigins.S3BucketOriginWithOACProps{
		OriginAccessLevels: &[]cloudfront.AccessLevel{
			cloudfront.AccessLevel_READ,
			cloudfront.AccessLevel_LIST,
		},
	})
	frontendDistribution := cloudfront.NewDistribution(stack, jsii.String("FrontendDistribution"), &cloudfront.DistributionProps{
		DefaultBehavior: &cloudfront.BehaviorOptions{
			Origin: s3Origin,
			//CachePolicy: cloudfront.CachePolicy_CACHING_DISABLED(), // TODO: in prod, caching should be enabled
			//, &cforigins.S3BucketOriginProps{
			//	OriginAccessIdentity: nil, // Or your OAI if you're using one
			//}),
			ViewerProtocolPolicy: cloudfront.ViewerProtocolPolicy_REDIRECT_TO_HTTPS,
		},
		DefaultRootObject: jsii.String("index.html"),
		ErrorResponses: &[]*cloudfront.ErrorResponse{
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
	})

	// Deploy the frontend files to the S3 bucket (and generate config.json)
	s3deployment.NewBucketDeployment(stack, jsii.String("DeployFrontend"), &s3deployment.BucketDeploymentProps{
		Sources: &[]s3deployment.ISource{
			s3deployment.Source_Asset(jsii.String(filepath.Join("..", "frontend", "dist")), nil),
			s3deployment.Source_JsonData(jsii.String("config.json"), struct {
				ApiUrl string `json:"api_url"`
			}{
				ApiUrl: apiUrl,
			}),
		},
		DestinationBucket: frontendBucket,

		Distribution: frontendDistribution,
		DistributionPaths: &[]*string{
			jsii.String("/*"),
		},
	})

	frontendUrlOutput := awscdk.NewCfnOutput(stack, jsii.String("FrontendUrl"), &awscdk.CfnOutputProps{
		ExportName: jsii.String("FrontendUrl"),
		Value:      jsii.String(fmt.Sprintf("https://%s", *frontendDistribution.DomainName())),
	})

	return &WebAppStack{
		NestedStack:       stack,
		ApiUrlOutput:      apiUrlOutput,
		FrontendUrlOutput: frontendUrlOutput,
	}
}
