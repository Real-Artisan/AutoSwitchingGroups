import boto3

def stop_auto_scaling_groups(region):
    ec2_client = boto3.client('ec2', region_name=region)
    autoscaling_client = boto3.client('autoscaling', region_name=region)

    # Get all Auto Scaling groups in the specified region
    response = autoscaling_client.describe_auto_scaling_groups()
    auto_scaling_groups = response['AutoScalingGroups']

    if not auto_scaling_groups:
        print("No Auto Scaling groups found in the specified region.")
        return

    for group in auto_scaling_groups:
        group_name = group['AutoScalingGroupName']
        
        # Get instance IDs from the Auto Scaling group
        instance_ids = [instance['InstanceId'] for instance in group['Instances']]
        ec2_client.terminate_instances(InstanceIds=instance_ids)

        autoscaling_client.update_auto_scaling_group(
            AutoScalingGroupName=group_name,
            DesiredCapacity=0,
            MaxSize=0,
            MinSize=0,
        )

        # Suspend scaling processes to prevent automatic instance launches
        autoscaling_client.suspend_processes(AutoScalingGroupName=group_name)

        print(f"All instances in Auto Scaling group {group_name} stopped and scaling processes suspended.")

def resume_auto_scaling_groups(region):
    autoscaling_client = boto3.client('autoscaling', region_name=region)

    # Get all Auto Scaling groups in the specified region
    response = autoscaling_client.describe_auto_scaling_groups()
    auto_scaling_groups = response['AutoScalingGroups']

    if not auto_scaling_groups:
        print("No Auto Scaling groups found in the specified region.")
        return

    for group in auto_scaling_groups:
        group_name = group['AutoScalingGroupName']

        autoscaling_client.update_auto_scaling_group(
            AutoScalingGroupName=group_name,
            DesiredCapacity=1,
            MaxSize=1,
            MinSize=1,
        )

        # Resume scaling processes
        autoscaling_client.resume_processes(AutoScalingGroupName=group_name)

        print(f"Scaling processes for Auto Scaling group {group_name} resumed.")

# Interactive CLI
def interactive_cli():
    region = input("Enter the AWS region: ")

    while True:
        action = input("Enter 'stop' to stop all instances, 'resume' to resume scaling processes, or 'q' to quit: ")

        if action == 'stop':
            stop_auto_scaling_groups(region)
        elif action == 'resume':
            resume_auto_scaling_groups(region)
        elif action == 'q':
            print("Exiting the program...")
            break
        else:
            print("Invalid input. Please try again.")

# Run
interactive_cli()
