package minio

// MinioImageName is minio image repository name
const MinioImageName string = "minio/minio"

// MinioImageTag is minio image tag
const MinioImageTag string = "latest"

// MinioContainerNamePrefix is minio container name prefix
const MinioContainerNamePrefix string = "taco-minio-"

// MinioDockerRunTemplate is minio docker run command
const MinioDockerRunTemplate string = "docker run -d --restart=always -p %v:9000 --name taco-minio-%s -v %s:/data -e 'MINIO_ACCESS_KEY=%s' -e 'MINIO_SECRET_KEY=%s' %s:%s server /data"

// MinioContainerExistsTemplate is minio container exists checking command
const MinioContainerExistsTemplate string = "docker ps | grep taco-minio-%s | awk '{print $1}'"

// MinioRemoveContainerTemplate is minio container removing
const MinioRemoveContainerTemplate string = "docker rm -f taco-minio-%s"

// MinioGetContainerPortTemplate is getting port number
const MinioGetContainerPortTemplate string = "docker inspect taco-minio-%s | grep -Ei '\"HostPort\":.+?\"([0-9]+)\"' | grep -oEi '([0-9]+)' | head -1"

// MinioMinPort is minium port
const MinioMinPort int = 9001
