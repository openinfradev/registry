package minio

// MinioImageName is minio image repository name
const MinioImageName string = "minio/minio"

// MinioImageTag is minio image tag
const MinioImageTag string = "latest"

// MinioDockerRunTemplate is minio docker run command
const MinioDockerRunTemplate string = "docker run -d --restart=always -p %d:%d --name taco-minio-%s -v %s:/data -e \"MINIO_ACCESS_KEY=%s\" -e \"MINIO_SECRET_KEY=%s\" %s:%s server /data"
