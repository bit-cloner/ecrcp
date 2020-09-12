Docker have recently made some changes to their policy regarding public docker images. Docker have removed the ability to freely pull images from docker hub
unrestricted. This means anonymous users cannot pull images beyond a certain rate limit.Systems that were configured to pull public images from docker hub will 
stop working. The CLI tool I am going to talk about in this post will help to get around this limitation if you happen to have a private repository. In particular 
I am going to focus on AWS's ECR in this post. Many AWS customers already have ECR as their private repository but might be pulling lot of public images from dockerhub
as this is the defult for docker deamon. If you can proxy or cache these public images onto your private repository you can avoid hitting the rate limit that docker has set.

### What is ecrcp
ecrcp aims to mimic cp command in Linux systems as closely as possible in its implementation. Consider ecrcp to be the cp equivalent to copy container images from docker hub to
ECR. 

### Prerequisites
1. ecrcp relies on docker client . It uses the same docker client to pull and push images to docker hub and AWS's ECR registry. 
2. AWS credentials that have permissions to create a repository and push to a repsotory in AWS ECR.

### How to get it
``` 
wget https://github.com/bit-cloner/ecrcp/releases/download/0.9/ecrcp
```
https://github.com/bit-cloner/ecrcp

### Usage
``` 
chmod +x ./ecrcp
ecrcp nginx:latest XXXXXXXXXXXX.dkr.ecr.eu-west-1.amazonaws.com/<REPO-NAME>/<IMAGE-NAME>:<TAG> 
```
### How is this useful 

Imagine there are a number of public images that your container runtime/ orichestration systems rely on. Getting a copy of these images into a local container repository can be tedious and manual to
do a docker pull , create a repository in AWS ECR, Tag it apporpriately and push it to target repository. Using this tool you can do it from a single terminal all of the above steps or you can
use it in a script to do the same. 