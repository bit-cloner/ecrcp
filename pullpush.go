package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func pullpush(dockerimg string, ecrurl string) {

	// check if cli has aws credentials
	if os.Getenv("AWS_SECRET_ACCESS_KEY") == "" || os.Getenv("AWS_ACCESS_KEY_ID") == "" || os.Getenv("AWS_DEFAULT_REGION") == "" {
		fmt.Println("AWS credentails not found. Set these credentials uisng aws configure or through environment variables")
		fmt.Println("Make sure these variables are set and have appropriate permission for route53 AWS_SECRET_ACCESS_KEY AWS_ACCESS_KEY_ID AWS_DEFAULT_REGION")
		fmt.Println("aborting..")
		os.Exit(1)
	}

	// pull docker hub image
	ctx := context.Background()
	clnt, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	fmt.Println("Pulling image: ", dockerimg)
	out, err := clnt.ImagePull(ctx, dockerimg, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	defer out.Close()

	io.Copy(os.Stdout, out)

	// Push to ecr location

	//get credentials to login docker deamon into ecr repository
	fmt.Println("Found AWS credentials")
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := ecr.New(sess)
	input := &ecr.GetAuthorizationTokenInput{} // you can give ids of registry you want to login to . defaults to loging to default registry if not present
	dockerconfig, err := svc.GetAuthorizationToken(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	reftoken := dockerconfig.AuthorizationData[0].AuthorizationToken
	token := *reftoken
	decodedtoken, err := base64.StdEncoding.DecodeString(token)
	password := string(decodedtoken)
	password = strings.Replace(password, "AWS:", "", -1)
	refendpoint := dockerconfig.AuthorizationData[0].ProxyEndpoint
	endpoint := *refendpoint
	fmt.Println("Got values for docker username and password to logon to ecr registry")
	fmt.Println("Creating destination repository just in case it doesn't exist already")
	// create destination repository
	// extract repository name from ecr url
	repo := between(ecrurl, "amazonaws.com/", ":")
	fmt.Println("creating repo: ", repo)
	repoinput := &ecr.CreateRepositoryInput{
		RepositoryName:     aws.String(repo),
		ImageTagMutability: aws.String("IMMUTABLE"),
		ImageScanningConfiguration: &ecr.ImageScanningConfiguration{
			ScanOnPush: aws.Bool(true)},
	}
	result, err := svc.CreateRepository(repoinput)
	if err != nil {
		fmt.Println(err)

	}
	fmt.Println(result)

	fmt.Println("Tagging pulled image as: ", ecrurl)
	// tag local image to remote repository
	err1 := clnt.ImageTag(ctx, dockerimg, ecrurl)
	if err1 != nil {
		panic(err1)
	}

	// push image to ECR
	fmt.Println("Pushing image to ECR repository")
	auth := types.AuthConfig{
		Username:      "AWS",
		Password:      password,
		ServerAddress: endpoint,
	}
	encodedJSON, err := json.Marshal(auth)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	pushresult, err := clnt.ImagePush(ctx, ecrurl, types.ImagePushOptions{RegistryAuth: authStr})
	if err != nil {
		panic(err)
	}
	//Parse the responses of image push
	responses, err := ioutil.ReadAll(pushresult)
	fmt.Println(string(responses))

	defer out.Close()

}

func between(value string, a string, b string) string {
	// Get substring between two strings.
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}
