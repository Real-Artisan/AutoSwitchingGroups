import boto3

def stop_auto_scaling_group(auto_scaling_group_name):
    ec2_client = boto3.client('ec2')
    autoscaling_client = boto3.client('autoscaling')

    # Get instance IDs from the Auto Scaling group
    response = autoscaling_client.describe_auto_scaling_groups(AutoScalingGroupNames=[auto_scaling_group_name])
    instance_ids = [instance['InstanceId'] for instance in response['AutoScalingGroups'][0]['Instances']]

    # Terminate instances
    ec2_client.terminate_instances(InstanceIds=instance_ids)


    # Set desired capacity to 0 to stop all instances
    autoscaling_client.update_auto_scaling_group(
        AutoScalingGroupName=auto_scaling_group_name,
        DesiredCapacity=0,
        MaxSize=0,
        MinSize=0,
    )

    # Suspend scaling processes to prevent automatic instance launches
    autoscaling_client.suspend_processes(AutoScalingGroupName=auto_scaling_group_name)

    print("All instances in Auto Scaling group " + auto_scaling_group_name + " stopped and scaling processes suspended.")

def resume_auto_scaling_group(auto_scaling_group_name):
    autoscaling_client = boto3.client('autoscaling')

    # Set max size to the total number of instances to allow instance launches
    autoscaling_client.update_auto_scaling_group(
        AutoScalingGroupName=auto_scaling_group_name,
        DesiredCapacity=1,
        MaxSize=1,
        MinSize=1,
    )

    # Resume scaling processes to allow automatic instance launches
    autoscaling_client.resume_processes(AutoScalingGroupName=auto_scaling_group_name)

    print("Scaling processes for Auto Scaling group " + auto_scaling_group_name + " resumed.")

# Interactive CLI
def interactive_cli():

    while True:
        action = input("Enter 'stop' to stop all instances, 'resume' to resume scaling processes, or 'q' to quit: ")

        if action == 'stop':
            auto_scaling_group_name = input("Enter the Auto Scaling group name: ")
            stop_auto_scaling_group(auto_scaling_group_name)
        elif action == 'resume':
            auto_scaling_group_name = input("Enter the Auto Scaling group name: ")
            resume_auto_scaling_group(auto_scaling_group_name)
        elif action == 'q':
            print("Exiting the program...")
            break
        else:
            print("Invalid input. Please try again.")

# Run the interactive CLI
interactive_cli()
