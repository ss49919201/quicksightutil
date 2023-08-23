package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/quicksight"
	"github.com/aws/aws-sdk-go-v2/service/quicksight/types"
	"github.com/aws/smithy-go"
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
	describeOutput, err := client.DescribeDataSet(ctx, &quicksight.DescribeDataSetInput{
		AwsAccountId: accountID,
		DataSetId:    srcID,
	})
	if err != nil {
		log.Printf("unable to describe data set, %v", err)
		os.Exit(1)
	}

	describePermissionOutput, err := client.DescribeDataSetPermissions(ctx, &quicksight.DescribeDataSetPermissionsInput{
		AwsAccountId: accountID,
		DataSetId:    srcID,
	})
	if err != nil {
		log.Printf("unable to describe data set permissions, %v", err)
		os.Exit(1)
	}

	describeDataSetRefreshPropertiesOutput, err := client.DescribeDataSetRefreshProperties(ctx, &quicksight.DescribeDataSetRefreshPropertiesInput{
		AwsAccountId: accountID,
		DataSetId:    srcID,
	})
	if err != nil {
		// cf. https://github.com/aws/aws-sdk-go-v2/issues/1110
		isResourceNotFoundException := func() bool {
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) {
				if _, ok := apiErr.(*types.ResourceNotFoundException); ok {
					return true
				}
			}
			return false
		}()

		if !isResourceNotFoundException {
			log.Printf("unable to describe data set refresh properties, %v", err)
			os.Exit(1)
		}
	}

	if _, err := client.CreateDataSet(ctx, &quicksight.CreateDataSetInput{
		AwsAccountId:                       accountID,
		DataSetId:                          dstID,
		ImportMode:                         describeOutput.DataSet.ImportMode,
		Name:                               dstID,
		PhysicalTableMap:                   describeOutput.DataSet.PhysicalTableMap,
		ColumnGroups:                       describeOutput.DataSet.ColumnGroups,
		ColumnLevelPermissionRules:         describeOutput.DataSet.ColumnLevelPermissionRules,
		DataSetUsageConfiguration:          describeOutput.DataSet.DataSetUsageConfiguration,
		DatasetParameters:                  describeOutput.DataSet.DatasetParameters,
		FieldFolders:                       describeOutput.DataSet.FieldFolders,
		LogicalTableMap:                    describeOutput.DataSet.LogicalTableMap,
		Permissions:                        describePermissionOutput.Permissions,
		RowLevelPermissionDataSet:          describeOutput.DataSet.RowLevelPermissionDataSet,
		RowLevelPermissionTagConfiguration: describeOutput.DataSet.RowLevelPermissionTagConfiguration,
		// Tags: TODO,
	}); err != nil {
		log.Printf("unable to describe data set refresh properties, %v", err)
		os.Exit(1)
	}

	if describeDataSetRefreshPropertiesOutput != nil {
		if _, err := client.PutDataSetRefreshProperties(ctx, &quicksight.PutDataSetRefreshPropertiesInput{
			AwsAccountId:             accountID,
			DataSetId:                dstID,
			DataSetRefreshProperties: describeDataSetRefreshPropertiesOutput.DataSetRefreshProperties,
		}); err != nil {
			log.Printf("unable to describe data set refresh properties, %v", err)
			os.Exit(1)
		}
	}
}
