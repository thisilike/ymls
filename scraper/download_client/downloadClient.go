package download_client

// TODO
type DownloadClient interface {
	Open() error
	Close() error
	Get(string) ([]byte, error)
}

func NewDownloadClient(config map[string]interface{}) (DownloadClient, error) {
	return DownloadClientRegister[config["type"].(string)](config)
}
