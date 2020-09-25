### What is ecrcp
ecrcp aims to mimic cp command in Linux systems as closely as possible in its implementation. Consider ecrcp to be the cp equivalent to copy container images from docker hub to
ECR. 

Docker have recently made some changes to their policy regarding public docker images. Docker have removed the ability to freely pull images from docker hub
unrestricted. This means anonymous users cannot pull images beyond a certain rate limit. Systems that were configured to pull public images from docker hub will 
stop working. The CLI tool I am going to talk about in this post will help to get around this limitation if you happen to have a private repository. In particular 
I am going to focus on AWS's ECR in this post. Many AWS customers already have ECR as their private repository but might be pulling lot of public images from dockerhub
as this is the defult for docker deamon. If you can proxy or cache these public images onto your private repository you can avoid hitting the rate limit that docker has set.

### Prerequisites
1. ecrcp relies on docker client . It uses the same docker client to pull and push images to docker hub and AWS's ECR registry. 
2. AWS credentials that have permissions to create a repository and push to a repsotory in AWS ECR.

### Required AWS permissions

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": [
                "ecr:PutLifecyclePolicy",
                "ecr:PutImageTagMutability",
                "ecr:DescribeImageScanFindings",
                "ecr:StartImageScan",
                "ecr:GetLifecyclePolicyPreview",
                "ecr:CreateRepository",
                "ecr:GetDownloadUrlForLayer",
                "ecr:PutImageScanningConfiguration",
                "ecr:GetAuthorizationToken",
                "ecr:ListTagsForResource",
                "ecr:UploadLayerPart",
                "ecr:PutImage",
                "ecr:BatchGetImage",
                "ecr:CompleteLayerUpload",
                "ecr:DescribeImages",
                "ecr:TagResource",
                "ecr:InitiateLayerUpload",
                "ecr:BatchCheckLayerAvailability",
                "ecr:GetLifecyclePolicy"
            ],
            "Resource": "*"
        }
    ]
}
```

### How to get it
``` 
wget https://github.com/bit-cloner/ecrcp/releases/download/0.9/ecrcp
```
https://github.com/bit-cloner/ecrcp

### Usage
``` 
chmod +x ./ecrcp
ecrcp nginx:latest XXXXXXXXXXXX.dkr.ecr.<REGION>.amazonaws.com/<REPO-NAME>/<IMAGE-NAME>:<TAG> 
```
### How is this useful 

When there are a number of public images that your container runtime/orchestration systems rely on. Getting a copy of these images into a local container repository can be tedious and manual.
you would need to 
1. Do a docker pull 
2. Create a repository in AWS ECR
3. Tag it apporpriately 
4. Push it to target repository

This tool you can do all of the above from a single terminal. Also you can use it in a script/CI/CD to do the same.

### What's next
1. I plan to build a similar tool for GCR and ACR
2. Make these as linux packages so they can be installed using for ex: ```apt install ecrcp ``` 