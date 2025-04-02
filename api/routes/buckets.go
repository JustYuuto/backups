package routes

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func GetBuckets(c *gin.Context) {
	var accessKeyId = os.Getenv("S3_ACCESS_KEY_ID")
	var accessKeySecret = os.Getenv("S3_SECRET_ACCESS_KEY")
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion("auto"),
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Unable to load AWS SDK config"})
		log.Fatalf("unable to load SDK config, %v", err)
		return
	}
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(os.Getenv("S3_ENDPOINT"))
	})

	buckets, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to load your buckets"})
		log.Println(fmt.Sprintf("failed to load your buckets, %v\n", err))
		return
	}

	var bucketsFinal []string
	for _, object := range buckets.Buckets {
		name := *object.Name
		bucketsFinal = append(bucketsFinal, name)
	}

	c.JSON(200, bucketsFinal)
}
