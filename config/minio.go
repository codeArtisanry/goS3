package config

type MinioConfig struct {
	Endpoint        string `envconfig:"MINIO_ENDPOINT"`
	AccessKey       string `envconfig:"MINIO_ACCESS_KEY"`
	SecretAccessKey string `envconfig:"MINIO_SECRET_KEY"`
	Region          string `envconfig:"MINIO_REGION"`
	Bucket          string `envconfig:"BUCKET"`
}
