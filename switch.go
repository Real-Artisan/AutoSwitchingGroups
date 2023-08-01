package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

func main() {
	minSize := flag.Int64("min", 1, "Minimum desired capacity for Autoscaling groups")
	maxSize := flag.Int64("max", 1, "Maximum desired capacity for Autoscaling groups")
	desiredCap := flag.Int64("desired", 1, "Desired capacity for Autoscaling groups")

	flag.Parse()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("us-east-2"),
		},
	}))

	var err error

	asg := autoscaling.New(sess)

	instances := &autoscaling.DescribeAutoScalingGroupsInput{
		MaxRecords: aws.Int64(100),
	}

	fmt.Println(instances)

	result, err := asg.DescribeAutoScalingGroups(instances)
	if err != nil {
		panic(err)
	}

	for _, group := range result.AutoScalingGroups {
		for _, tag := range group.Tags {
			if *tag.Key == "Name" && strings.HasPrefix(*tag.Value, "Dev") {
				fmt.Println(*group.AutoScalingGroupName)
				fmt.Printf("Updating autoscaling group %s\n", *group.AutoScalingGroupName)

				// Your custom function to set desired capacity, min size, and max size
				updateAutoScalingGroupCapacity(asg, group.AutoScalingGroupName, *minSize, *maxSize, *desiredCap)

				break // no need to check remaining tags for this group
			}
		}
	}
}

// Custom function to update Auto Scaling group capacity
func updateAutoScalingGroupCapacity(asg *autoscaling.AutoScaling, autoScalingGroupName *string, minSize, maxSize, desiredCap int64) {
	_, err := asg.UpdateAutoScalingGroup(&autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: autoScalingGroupName,
		MinSize:              aws.Int64(minSize),
		MaxSize:              aws.Int64(maxSize),
		DesiredCapacity:      aws.Int64(desiredCap),
	})

	if err != nil {
		fmt.Printf("Failed to update autoscaling group %s: %v\n", *autoScalingGroupName, err)
	} else {
		fmt.Printf("Successfully updated autoscaling group %s\n", *autoScalingGroupName)
	}
}
