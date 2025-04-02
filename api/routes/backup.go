package routes

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"
)

type RequestBody struct {
	Path       string `json:"path"`
	Bucket     string `json:"bucket"`
	BucketPath string `json:"bucket_path"`
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func addToZip(zipWriter *zip.Writer, basePath, filePath string) error {
	relPath, err := filepath.Rel(basePath, filePath)
	if err != nil {
		return err
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		_, err := zipWriter.Create(relPath + "/")
		return err
	}

	// Ajouter le fichier
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	zipFileWriter, err := zipWriter.Create(relPath)
	if err != nil {
		return err
	}

	_, err = io.Copy(zipFileWriter, file)
	return err
}

func zipFolder(sourceDir, zipFileName string) error {
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return addToZip(zipWriter, sourceDir, path)
	})
}

func StartBackup(c *gin.Context) {
	var body RequestBody
	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	fileExists, err := exists(body.Path)
	if err != nil || !fileExists {
		c.JSON(400, gin.H{"error": "File or directory does not exists"})
		return
	}

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

	path := fmt.Sprintf("%s/backup-%v.zip", os.TempDir(), time.Now().UnixMilli())
	err = zipFolder(body.Path, path)
	if err != nil {
		c.JSON(500, gin.H{"error": "There was an error during zip creation"})
		return
	}

	file, err := os.Open(path)
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Key:         aws.String(body.BucketPath),
		Bucket:      aws.String(body.Bucket),
		ContentType: aws.String("application/zip"),
		Body:        file,
	})
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "There was an error during zip upload"})
		return
	}

	c.Writer.WriteHeader(204)
}
