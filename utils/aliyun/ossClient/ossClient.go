package ossClient

import (
	"log"
	"mime/multipart"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/credentials-go/credentials"
)

type Credentials struct {
	AccessKeyId     string
	AccessKeySecret string
	SecurityToken   string
}

type CredentialsProvider struct {
	cred credentials.Credential
}

func (credentials *Credentials) GetAccessKeyID() string {
	return credentials.AccessKeyId
}

func (credentials *Credentials) GetAccessKeySecret() string {
	return credentials.AccessKeySecret
}

func (credentials *Credentials) GetSecurityToken() string {
	return credentials.SecurityToken
}

func (defBuild CredentialsProvider) GetCredentials() oss.Credentials {
	id, _ := defBuild.cred.GetAccessKeyId()
	secret, _ := defBuild.cred.GetAccessKeySecret()
	token, _ := defBuild.cred.GetSecurityToken()

	return &Credentials{
		AccessKeyId:     *id,
		AccessKeySecret: *secret,
		SecurityToken:   *token,
	}
}

func NewStaticCredentialsProvider(credential credentials.Credential) CredentialsProvider {
	return CredentialsProvider{
		cred: credential,
	}
}

type MyBucket struct {
	*oss.Bucket
}

func NewBucket() (*MyBucket, error) {
	config := new(credentials.Config).
		// 指定Credential类型，固定值为ecs_ram_role。
		SetType("ecs_ram_role").
		// （可选项）指定角色名称。如果不指定，OSS会自动获取角色。强烈建议指定角色名称，以降低请求次数。
		SetRoleName(os.Getenv("AliyunRoleName"))

	ecsCredential, err := credentials.NewCredential(config)
	if err != nil {
		return nil, err
	}
	provider := NewStaticCredentialsProvider(ecsCredential)
	client, err := oss.New(os.Getenv("EndPoint"), "", "", oss.SetCredentialsProvider(&provider))
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}
	// 存储空间名称
	bucket, err := client.Bucket(os.Getenv("BucketName"))
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}
	return &MyBucket{bucket}, nil
	// 填写本地文件的完整路径，例如D:\\localpath\\examplefile.txt。
	// fd, err := os.Open("go.sum")
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	os.Exit(-1)
	// }
	// defer fd.Close()

}

func (mb *MyBucket) UploadVideo(file *multipart.FileHeader, internalURL string) error {
	// 将文件流上传至exampledir目录下的exampleobject.txt文件。
	fileStrem, err := file.Open()
	if err != nil {
		return err
	}
	err = mb.PutObject(internalURL, fileStrem)
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	return nil
}

func (mb *MyBucket) ObjectExternalURL(internalURL string) (signedURL string, err error) {
	// 生成用于下载的签名URL，并指定签名URL的有效时间为60秒。
	for i := 0; i < 5; i++ {
		signedURL, err = mb.SignURL(internalURL, oss.HTTPGet, 600)
		if err != nil {
			log.Println("Error when get video URL:", err)
			continue
		}
		return signedURL, nil
	}
	return "", err
}
