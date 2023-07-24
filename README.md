### AutoSwitchingGroups
A tool created using Python 3 for stopping or resuming scaling processes in an autoscaling group. 

### How to use
Ensure That the minimum capacity of the auto scaling group is set to 0.

Clone this tool.

Install boto3
```
pip install boto3
```
Make sure your AWS credentials are set correctly with the appropriate permissions to interact with EC2 instances.

```
python3 switch.py
```
Enter 'stop' to stop all instances, 'resume' to resume scaling processes, or 'q' to quit:


### Developer
This tool was created by Daniel Pereowei Iwenya. <a href="mailto:iwenyadaniel12@gmail.com">Contact Developer.</a>
