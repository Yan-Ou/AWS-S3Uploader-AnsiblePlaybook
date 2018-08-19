# Prerequisite
1. Python 2.7.x is installed.
2. Ansible 2.6.x is installed.
3. bot0, botocore and boto3 are installed. 
4. AWS account with necessary permissions. 
5. AWSAccessKey and AWSSecretKey ready for running the playbook.

# Playbook tasks:
- Create 1 x VPC (1 x public subnet, 1 x Internet gateway, and 1 public route table)
- Create 1 x S3 bucket
- Create 1 x Key pair
- Provision 1 x EC2 instance 
- App role
  - Install Go packages
  - Install Nginx
  - Install Git
  - Copy the web app (main.go, view.html) to EC2 instance
  - Startup the web app
  
  # How to run the playbook
  1. Issue the following command:
    ansible-playbook site.yml -i hosts
  2. After issuing the command, input your AWSAccessKey, AWSSecretKey and region name which you want to deploy to. Default region name is ap-southeast2 #Sydney.
  
  # File structure:
  - site.yml: contains all the tasks
  - hosts: contains all inventories
  - ansible.cfg: ansible configuration 
  - src/niginx.conf: nginx service configuration
  - src/main.go: web app 
  - src/view.html: web app template
  - app: ansible role to deploy web app
  - aws: aws tasks to provision VPC, Keypair,EC2 and S3
  
  # Web app:
  1. Implemented by Go along with AWS SDK.
  2. Call AWS REST APIs to upload files to S3 bucket.
  3. Fire up web service and listen to port 8001(root path handler), 8002(upload path handler).
  4. Nginx wraps two path and bind them to port 80 for public access. 
  
