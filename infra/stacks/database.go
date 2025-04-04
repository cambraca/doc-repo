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
	Vpc                       ec2.IVpc
	AllowAccessFromEverywhere bool
}

type DatabaseStack struct {
	awscdk.NestedStack
	Instance      rds.IDatabaseInstance
	SecurityGroup ec2.ISecurityGroup
}

func NewDatabaseStack(scope constructs.Construct, id string, props *DatabaseStackProps) *DatabaseStack {
	stack := awscdk.NewNestedStack(scope, &id, &props.NestedStackProps)

	dbSecurityGroup := ec2.NewSecurityGroup(stack, jsii.String("RDSSecurityGroup"), &ec2.SecurityGroupProps{
		Vpc:              props.Vpc,
		Description:      jsii.String("Allow access to RDS PostgreSQL"),
		AllowAllOutbound: jsii.Bool(props.AllowAccessFromEverywhere),
	})
	var peer ec2.IPeer
	if props.AllowAccessFromEverywhere {
		// TODO: this allows access from anywhere
		peer = ec2.Peer_Ipv4(jsii.String("0.0.0.0/0"))
	} else {
		peer = ec2.Peer_Ipv4(props.Vpc.VpcCidrBlock())
	}
	dbSecurityGroup.AddIngressRule(
		peer,
		ec2.Port_Tcp(jsii.Number(5432)),
		jsii.String("Allow PostgreSQL access from within the VPC"),
		jsii.Bool(false), // TODO: is this doing anything?
	)
	var subnets *[]ec2.ISubnet
	var subnetId string
	if props.AllowAccessFromEverywhere {
		subnets = props.Vpc.PublicSubnets()
		subnetId = "RDSPublicSubnetGroup"
	} else {
		subnets = props.Vpc.PrivateSubnets()
		subnetId = "RDSSubnetGroup"
	}
	dbSubnetGroup := rds.NewSubnetGroup(stack, jsii.String(subnetId), &rds.SubnetGroupProps{
		Description: jsii.String("RDSSubnetGroup"),
		Vpc:         props.Vpc,
		VpcSubnets:  &ec2.SubnetSelection{Subnets: subnets},
	})
	dbInstance := rds.NewDatabaseInstance(stack, jsii.String("PostgresDB"), &rds.DatabaseInstanceProps{
		Engine: rds.DatabaseInstanceEngine_Postgres(&rds.PostgresInstanceEngineProps{
			Version: rds.PostgresEngineVersion_VER_17_4(), // Choose your version
		}),
		InstanceType:     ec2.InstanceType_Of(ec2.InstanceClass_T3, ec2.InstanceSize_MICRO), // TODO is it too small?
		Vpc:              props.Vpc,
		SubnetGroup:      dbSubnetGroup,
		SecurityGroups:   &[]ec2.ISecurityGroup{dbSecurityGroup},
		StorageType:      rds.StorageType_GP3,
		AllocatedStorage: jsii.Number(20),        // GiB (minimum: 20)
		DatabaseName:     jsii.String("docrepo"), // TODO: pass this to webapp stack
		Credentials: rds.Credentials_FromUsername( // TODO: good password and pass these to the webapp stack
			jsii.String("docrepouser"),
			&rds.CredentialsFromUsernameOptions{
				Password: awscdk.SecretValue_UnsafePlainText(jsii.String("abcd1234")),
			},
		),
		BackupRetention:    awscdk.Duration_Days(jsii.Number(7)),
		RemovalPolicy:      awscdk.RemovalPolicy_SNAPSHOT, // Consider RETAIN
		PubliclyAccessible: jsii.Bool(props.AllowAccessFromEverywhere),
	})

	awscdk.NewCfnOutput(stack, jsii.String("DatabaseEndpointAddress"), &awscdk.CfnOutputProps{
		ExportName: jsii.String("DatabaseEndpointAddress"),
		Value:      dbInstance.DbInstanceEndpointAddress(),
	})
	awscdk.NewCfnOutput(stack, jsii.String("DatabaseEndpointPort"), &awscdk.CfnOutputProps{
		ExportName: jsii.String("DatabaseEndpointPort"),
		Value:      dbInstance.DbInstanceEndpointPort(),
	})
	awscdk.NewCfnOutput(stack, jsii.String("DatabaseSecurityGroupId"), &awscdk.CfnOutputProps{
		ExportName: jsii.String("DatabaseSecurityGroupId"),
		Value:      dbSecurityGroup.SecurityGroupId(),
	})

	return &DatabaseStack{
		NestedStack:   stack,
		Instance:      dbInstance,
		SecurityGroup: dbSecurityGroup,
	}
}
