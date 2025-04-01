package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	rds "github.com/aws/aws-cdk-go/awscdk/v2/awsrds"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type DatabaseStackProps struct {
	awscdk.NestedStackProps
	Vpc ec2.IVpc
}

type DatabaseStack struct {
	awscdk.NestedStack
	EndpointAddressOutput *string
	EndpointPortOutput    *string
	SecurityGroupIdOutput *string
	Instance              rds.IDatabaseInstance
	SecurityGroup         ec2.ISecurityGroup
}

func NewDatabaseStack(scope constructs.Construct, id string, props *DatabaseStackProps) *DatabaseStack {
	stack := awscdk.NewNestedStack(scope, &id, &props.NestedStackProps)

	dbSecurityGroup := ec2.NewSecurityGroup(stack, jsii.String("RDSSecurityGroup"), &ec2.SecurityGroupProps{
		Vpc:              props.Vpc,
		Description:      jsii.String("Allow access to RDS PostgreSQL"),
		AllowAllOutbound: jsii.Bool(false),
	})
	dbSecurityGroup.AddIngressRule(
		ec2.Peer_Ipv4(props.Vpc.VpcCidrBlock()),
		ec2.Port_Tcp(jsii.Number(5432)),
		jsii.String("Allow PostgreSQL access from within the VPC"),
		jsii.Bool(false), // TODO is false right?
	)
	dbSubnetGroup := rds.NewSubnetGroup(stack, jsii.String("RDSSubnetGroup"), &rds.SubnetGroupProps{
		Description: jsii.String("asd"), // TODO
		Vpc:         props.Vpc,
		//VpcSubnets:         &props.Vpc.PrivateSubnets(),
		SubnetGroupName: jsii.String("rds-private-subnets"),
	})
	dbInstance := rds.NewDatabaseInstance(stack, jsii.String("PostgresDB"), &rds.DatabaseInstanceProps{
		Engine: rds.DatabaseInstanceEngine_Postgres(&rds.PostgresInstanceEngineProps{
			Version: rds.PostgresEngineVersion_VER_17_4(), // Choose your version
		}),
		InstanceType:     ec2.InstanceType_Of(ec2.InstanceClass_T3, ec2.InstanceSize_MICRO), // TODO is it too small?
		Vpc:              props.Vpc,
		SubnetGroup:      dbSubnetGroup,
		SecurityGroups:   &[]ec2.ISecurityGroup{dbSecurityGroup},
		AllocatedStorage: jsii.Number(1), // GiB
		//MasterUsername:   jsii.String("admin"),                                            // Replace
		//MasterPassword:   awscdk.SecretValue_PlainText(jsii.String("yourStrongPassword")), // Replace (consider Secrets Manager)
		BackupRetention: awscdk.Duration_Days(jsii.Number(7)),
		RemovalPolicy:   awscdk.RemovalPolicy_SNAPSHOT, // Consider RETAIN
	})

	endpointAddressOutput := awscdk.NewCfnOutput(stack, jsii.String("DatabaseEndpointAddress"), &awscdk.CfnOutputProps{
		ExportName: jsii.String("DatabaseEndpointAddress"),
		Value:      dbInstance.DbInstanceEndpointAddress(),
	})
	endpointPortOutput := awscdk.NewCfnOutput(stack, jsii.String("DatabaseEndpointPort"), &awscdk.CfnOutputProps{
		ExportName: jsii.String("DatabaseEndpointPort"),
		Value:      dbInstance.DbInstanceEndpointPort(),
	})
	securityGroupIdOutput := awscdk.NewCfnOutput(stack, jsii.String("DatabaseSecurityGroupId"), &awscdk.CfnOutputProps{
		ExportName: jsii.String("DatabaseSecurityGroupId"),
		Value:      dbSecurityGroup.SecurityGroupId(),
	})

	return &DatabaseStack{
		NestedStack:           stack,
		EndpointAddressOutput: endpointAddressOutput.ImportValue(),
		EndpointPortOutput:    endpointPortOutput.ImportValue(),
		SecurityGroupIdOutput: securityGroupIdOutput.ImportValue(),
		Instance:              dbInstance,
		SecurityGroup:         dbSecurityGroup,
	}
}
