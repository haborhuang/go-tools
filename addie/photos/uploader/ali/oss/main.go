package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"github.com/spf13/cobra"
)

var (
	Cmd = cobra.Command{
		Use: "Files uploader",
		Run: func(cmd *cobra.Command, args []string) {
			akData, err := ioutil.ReadFile(Opts.accessKeyFile)
			if nil != err {
				panic(fmt.Errorf("Failed to open access key file: %v", err))
			}

			var ak accessKey
			if err := yaml.Unmarshal(akData, &ak); nil != err {
				panic(fmt.Errorf("Failed to decode access key file: %v", err))
			}

			client, err := oss.New(Opts.domain, ak.Id, ak.Secret)
			if nil != err {
				panic(fmt.Errorf("New OSS client error: %v", err))
			}

			bucket, err := client.Bucket(Opts.bucket)
			if nil != err {
				panic(fmt.Errorf("Get bucket error: %v", err))
			}

			if err := upload(bucket); nil != err {
				panic(fmt.Errorf("Upload error: %v", err))
			}
		},
	}

	Opts struct {
		filePath      string
		keyPrefix     string
		accessKeyFile string
		domain        string
		bucket        string
	}
)

func init() {
	Cmd.Flags().StringVarP(&Opts.filePath, "file", "f", "upload/", "Path of file(s) or direcotry to be uploaded")
	Cmd.Flags().StringVarP(&Opts.keyPrefix, "prefix", "p", "", "Prefix of object key")
	Cmd.Flags().StringVarP(&Opts.domain, "domain", "d", "oss-cn-beijing.aliyuncs.com", "EndPoint domain")
	Cmd.Flags().StringVarP(&Opts.bucket, "bucket", "b", "addie", "Bucket name")
	Cmd.Flags().StringVarP(&Opts.accessKeyFile, "access-key", "k", "key.yaml", "File path of access key id and secret")
}

type accessKey struct {
	Id     string `yaml:"id"`
	Secret string `yaml:"secret"`
}

func main() {
	Cmd.Execute()
}

func upload(bucket *oss.Bucket) error {
	return filepath.Walk(Opts.filePath, func(fpath string, info os.FileInfo, err error) error {
		if nil != err {
			return err
		}

		if info.IsDir() {
			return nil
		}

		objKey := path.Join(Opts.keyPrefix, info.Name())
		log.Println("Uploading:", objKey)

		if err := bucket.PutObjectFromFile(objKey, fpath); nil != err {
			return fmt.Errorf("Put object error: %v", err)
		}

		return nil
	})
}
