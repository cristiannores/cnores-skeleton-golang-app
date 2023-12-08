package settings_connection

type SecureConnection struct {
	IsEnabled           bool
	KeyLocation         string
	KeyPassword         string
	CaLocation          string
	CertificateLocation string
}

type ConsumerConnection struct {
	FetchMessageMaxBytes int
	FetchMaxBytes        int
	MaxPollInterval      int
	SessionTimeout       int
	IntervalHeartbeat    int
}

type ProducerConnection struct {
	RequestTimeout int
}
