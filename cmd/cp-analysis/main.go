package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/quicksight"
)

func main() {
	accountID := flag.String("account-id", "", "AWS Account ID")
	region := flag.String("region", "", "AWS Region ID")
	srcID := flag.String("src-id", "", "Source DataSet ID")
	dstID := flag.String("dst-id", "", "Destination DataSet ID")

	flag.Parse()

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(aws.ToString(region)),
	)
	if err != nil {
		log.Printf("unable to load SDK config, %v", err)
		os.Exit(1)
	}

	client := quicksight.NewFromConfig(cfg)
	describeOutput, err := client.DescribeAnalysis(ctx, &quicksight.DescribeAnalysisInput{
		AwsAccountId: accountID,
		AnalysisId:   srcID,
	})
	if err != nil {
		log.Printf("unable to describe analysis, %v", err)
		os.Exit(1)
	}

	describePermissionOutput, err := client.DescribeAnalysisPermissions(ctx, &quicksight.DescribeAnalysisPermissionsInput{
		AwsAccountId: accountID,
		AnalysisId:   srcID,
	})
	if err != nil {
		log.Printf("unable to describe analysis permissions, %v", err)
		os.Exit(1)
	}

	describeAnalysisDefinitionOutput, err := client.DescribeAnalysisDefinition(ctx, &quicksight.DescribeAnalysisDefinitionInput{
		AwsAccountId: accountID,
		AnalysisId:   srcID,
	})
	if err != nil {
		log.Printf("unable to describe analysis defenitions, %v", err)
		os.Exit(1)
	}

	_, _, _ = describeOutput, describePermissionOutput, describeAnalysisDefinitionOutput

	if _, err := client.CreateAnalysis(ctx, &quicksight.CreateAnalysisInput{
		AwsAccountId: accountID,
		AnalysisId:   dstID,
		Name:         dstID,
		Definition:   describeAnalysisDefinitionOutput.Definition,
		Permissions:  describePermissionOutput.Permissions,
		ThemeArn:     describeOutput.Analysis.ThemeArn,
		// SourceEntity: TODO
		// Parameters: TODO
		// Tags: TODO
	}); err != nil {
		log.Printf("unable to create analysis, %v", err)
		os.Exit(1)
	}
}
