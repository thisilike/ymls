package download_client

import (
	"github.com/sirupsen/logrus"
	"github.com/thisilike/ymls/config"
)

var log *logrus.Logger

// NewFuntion for every DownloadClient is Registerd here
var DownloadClientRegister map[string]func(map[string]interface{}) (DownloadClient, error)

func init() {
	log = config.Logger
	DownloadClientRegister = make(map[string]func(map[string]interface{}) (DownloadClient, error))
	DownloadClientRegister["go-http"] = NewGoHttpDownloadClient
}
