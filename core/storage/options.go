package storage

// Option for setting params used by chunk manager client.
type config struct {
	address           string
	bucketName        string
	accessKeyID       string
	secretAccessKeyID string
	useSSL            bool
	createBucket      bool
	rootPath          string
	useIAM            bool
	iamEndpoint       string

	// deprecated
	cloudProvider string
	storageType   string

	backupAccessKeyID       string
	backupSecretAccessKeyID string
	backupBucketName        string
	backupRootPath          string
}

func newDefaultConfig() *config {
	return &config{}
}

// Option is used to config the retry function.
type Option func(*config)

func Address(addr string) Option {
	return func(c *config) {
		c.address = addr
	}
}

func BucketName(bucketName string) Option {
	return func(c *config) {
		c.bucketName = bucketName
	}
}

func AccessKeyID(accessKeyID string) Option {
	return func(c *config) {
		c.accessKeyID = accessKeyID
	}
}
func SecretAccessKeyID(secretAccessKeyID string) Option {
	return func(c *config) {
		c.secretAccessKeyID = secretAccessKeyID
	}
}

func UseSSL(useSSL bool) Option {
	return func(c *config) {
		c.useSSL = useSSL
	}
}

func CreateBucket(createBucket bool) Option {
	return func(c *config) {
		c.createBucket = createBucket
	}
}

func RootPath(rootPath string) Option {
	return func(c *config) {
		c.rootPath = rootPath
	}
}

func UseIAM(useIAM bool) Option {
	return func(c *config) {
		c.useIAM = useIAM
	}
}

func IAMEndpoint(iamEndpoint string) Option {
	return func(c *config) {
		c.iamEndpoint = iamEndpoint
	}
}
