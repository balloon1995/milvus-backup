package core

import (
	"context"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zilliztech/milvus-backup/core/storage/tencent"

	"github.com/zilliztech/milvus-backup/core/paramtable"
	"github.com/zilliztech/milvus-backup/core/storage"
	"github.com/zilliztech/milvus-backup/internal/log"
	"go.uber.org/zap"
)

var Params paramtable.BackupParams

func newMinioChunkManager(ctx context.Context, bucketName string) (*storage.MinioChunkManager, error) {
	endPoint := getMinioAddress()
	accessKeyID, _ := Params.Load("minio.accessKeyID")
	secretAccessKey, _ := Params.Load("minio.secretAccessKey")
	useSSLStr, _ := Params.Load("minio.useSSL")
	useSSL, _ := strconv.ParseBool(useSSLStr)

	storageConfig := &storage.StorageConfig{
		StorageType:       "minio",
		Address:           endPoint,
		AccessKeyID:       accessKeyID,
		SecretAccessKeyID: secretAccessKey,
		UseSSL:            useSSL,
		CreateBucket:      true,
		UseIAM:            false,
		IAMEndpoint:       "",
	}

	client, err := storage.NewMinioChunkManagerWithConfig(ctx, storageConfig)
	return client, err
}

func getMinioAddress() string {
	minioHost := Params.LoadWithDefault("minio.address", paramtable.DefaultMinioAddress)
	if strings.Contains(minioHost, ":") {
		return minioHost
	}
	port := Params.LoadWithDefault("minio.port", paramtable.DefaultMinioPort)
	return minioHost + ":" + port
}

func TestWriteAEmptyBackupFile(t *testing.T) {

	Params.Init()
	testBucket, err := Params.Load("minio.bucketName")
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testCM, err := newMinioChunkManager(ctx, testBucket)
	err = testCM.Write(ctx, testBucket, "backup/test_backup6", nil)
	assert.NoError(t, err)
}

func TestReadBackupFiles(t *testing.T) {

	Params.Init()
	testBucket, err := Params.Load("minio.bucketName")
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testCM, err := newMinioChunkManager(ctx, testBucket)
	files, _, err := testCM.ListWithPrefix(ctx, testBucket, "/backup", true)
	assert.NoError(t, err)

	for _, file := range files {
		log.Info("BackupFiles", zap.String("path", file))
	}

}

func TestReadMilvusData(t *testing.T) {
	var params paramtable.BackupParams
	params.GlobalInitWithYaml("backup.yaml")
	params.Init()

	context := context.Background()
	//backupContext := CreateBackupContext(context, params)

	client, err := createStorageClient(context, params)
	assert.NoError(t, err)
	paths, _, err := client.ListWithPrefix(context, params.MinioCfg.BucketName, "file/insert_log/437296492118216229/437296492118216230/", true)
	assert.NoError(t, err)
	for _, path := range paths {
		if strings.Contains(path, "index_files") {
			continue
		}
		if strings.Contains(path, "437296588890839162") ||
			strings.Contains(path, "437296588890833721") ||
			strings.Contains(path, "437296584963129351") ||
			strings.Contains(path, "437296581056135434") ||
			strings.Contains(path, "437296588890833719") {
			log.Info(path)
			bytes, err := client.Read(context, params.MinioCfg.BucketName, path)
			os.MkdirAll(path, os.ModePerm)
			os.Remove(path)
			err = os.WriteFile(path, bytes, 0666)
			assert.NoError(t, err)
			//log.Info("paths", zap.Strings("paths", paths))
		}

		//log.Info(path)
	}

}

func TestTencentOSS(t *testing.T) {
	var creds *credentials.Credentials
	bucketLookupType := minio.BucketLookupDNS
	minioOpts := &minio.Options{
		BucketLookup: bucketLookupType,
		Creds:        creds,
		Secure:       true,
	}
	client, err := tencent.NewMinioClient("cos.ap-nanjing.myqcloud.com", minioOpts)
	assert.NoError(t, err)
	exist, err := client.BucketExists(context.Background(), "")
	assert.NoError(t, err)
	assert.True(t, exist)
}
