package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func stopAutoScalingGroup(autoScalingGroupName string) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		fmt.Println("Error loading AWS config:", err)
		return
	}

	ec2Svc := ec2.New(cfg)
	autoscalingSvc := autoscaling.New(cfg)

	// Get instance IDs from the Auto Scaling group
	asgInput := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []string{autoScalingGroupName},
	}
	asgOutput, err := autoscalingSvc.DescribeAutoScalingGroups(asgInput)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var instanceIDs []string
	for _, instance := range asgOutput.AutoScalingGroups[0].Instances {
		instanceIDs = append(instanceIDs, aws.ToString(instance.InstanceId))
	}

	// Terminate instances
	terminateInput := &ec2.TerminateInstancesInput{
		InstanceIds: instanceIDs,
	}
	_, err = ec2Svc.TerminateInstancesRequest(terminateInput).Send()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Set desired capacity to 0 to stop all instances
	updateInput := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(autoScalingGroupName),
		DesiredCapacity:      aws.Int64(0),
		MaxSize:              aws.Int64(0),
		MinSize:              aws.Int64(0),
	}
	_, err = autoscalingSvc.UpdateAutoScalingGroupRequest(updateInput).Send()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Suspend scaling processes to prevent automatic instance launches
	suspendInput := &autoscaling.SuspendProcessesInput{
		AutoScalingGroupName: aws.String(autoScalingGroupName),
	}
	_, err = autoscalingSvc.SuspendProcessesRequest(suspendInput).Send()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("All instances in Auto Scaling group", autoScalingGroupName, "stopped and scaling processes suspended.")
}

func resumeAutoScalingGroup(autoScalingGroupName string) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		fmt.Println("Error loading AWS config:", err)
		return
	}

	autoscalingSvc := autoscaling.New(cfg)

	// Set max size to the total number of instances to allow instance launches
	updateInput := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: aws.String(autoScalingGroupName),
		DesiredCapacity:      aws.Int64(1),
		MaxSize:              aws.Int64(1),
		MinSize:              aws.Int64(1),
	}
	_, err = autoscalingSvc.UpdateAutoScalingGroupRequest(updateInput).Send()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Resume scaling processes to allow automatic instance launches
	resumeInput := &autoscaling.ResumeProcessesInput{
		AutoScalingGroupName: aws.String(autoScalingGroupName),
	}
	_, err = autoscalingSvc.ResumeProcessesRequest(resumeInput).Send()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Scaling processes for Auto Scaling group", autoScalingGroupName, "resumed.")
}


func interactiveCLI() {
	var action string
	for {
		fmt.Print("Enter 'stop' to stop all instances, 'resume' to resume scaling processes, or 'q' to quit: ")
		fmt.Scanln(&action)

		if action == "stop" {
			var autoScalingGroupName string
			fmt.Print("Enter the Auto Scaling group name: ")
			fmt.Scanln(&autoScalingGroupName)
			stopAutoScalingGroup(autoScalingGroupName)
		} else if action == "resume" {
			var autoScalingGroupName string
			fmt.Print("Enter the Auto Scaling group name: ")
			fmt.Scanln(&autoScalingGroupName)
			resumeAutoScalingGroup(autoScalingGroupName)
		} else if action == "q" {
			fmt.Println("Exiting the program...")
			break
		} else {
			fmt.Println("Invalid input. Please try again.")
		}
	}
}

func main() {
	interactiveCLI()
}
