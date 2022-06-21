package scraper

import (
	"errors"
	"time"

	"github.com/thisilike/ymls/scraper/data_collectors"
	"github.com/thisilike/ymls/scraper/download_client"
	"github.com/thisilike/ymls/scraper/storage_interfaces"
)

type Scraper struct {
	Name                   string
	URLs                   []string
	DownloadClient         download_client.DownloadClient
	Retrys                 int64
	OnFailure              string
	Workers                int64
	Delay                  int64
	DataCollectors         []data_collectors.DataCollector
	StorageInterfaces      []storage_interfaces.StorageInterface
	DataCollectorsRegister map[string]bool
}

type ScrapeResult struct {
	url      string
	errors   []error
	succeded bool
}

func (scraper *Scraper) StartScraping() error {
	log.Infof("Starting scraping on scraper: '%s'", scraper.Name)
	taskPool := make(chan ScrapeResult, len(scraper.URLs))
	defer close(taskPool)
	results := make(chan ScrapeResult, len(scraper.URLs))
	defer close(results)
	for _, url := range scraper.URLs {
		taskPool <- ScrapeResult{url: url}
	}
	go func() {
		scrp := func(task ScrapeResult) {
			log.Debugf("scraping: '%s'", task.url)
			res := scraper.Scrape(task)
			if res.succeded {
				results <- res
			} else if len(res.errors) >= int(scraper.Retrys) {
				log.Errorf(
					"failed to download: '%s' with error: '%s', try: '%s'",
					res.url,
					res.errors[len(res.errors)-1],
					len(res.errors),
				)
				results <- res
			} else {
				log.Warnf(
					"failed to download: '%s' with error: '%s', try: '%s'",
					res.url,
					res.errors[len(res.errors)-1],
					len(res.errors),
				)
				taskPool <- res
			}
		}
		for task := range taskPool {
			go scrp(task)
			time.Sleep(time.Duration(scraper.Delay) * time.Millisecond)
		}
	}()
	count := 0
	for result := range results {
		count++
		if !result.succeded && scraper.OnFailure == "abort" {
			log.Errorf("failure detected. aborting")
			for i, e := range result.errors {
				log.Errorf("Error %s: '%s'", i, e.Error())
			}
			break
		} else if count == len(scraper.URLs) {
			return nil
		}
	}
	return nil
}

func (scraper *Scraper) Scrape(task ScrapeResult) ScrapeResult {
	log.Debugf("downloading: '%s'", task.url)
	data, err := scraper.DownloadClient.Get(task.url)
	if err != nil {
		task.errors = append(task.errors, err)
		return task
	}
	log.Debugf("collecting: '%s'", task.url)
	scrapedData, err := scraper.Collect(data)
	if err != nil {
		task.errors = append(task.errors, err)
		return task
	}
	log.Debugf("saving: '%s'", task.url)
	err = scraper.Save(scrapedData)
	if err != nil {
		task.errors = append(task.errors, err)
		return task
	}
	task.succeded = true
	return task
}

func (scraper *Scraper) Collect(docData []byte) (map[string]interface{}, error) {
	scrapedData := make(map[string]interface{})
	for _, dataCollector := range scraper.DataCollectors {
		data, err := dataCollector.Collect(docData)
		if err != nil {
			log.Errorf("failed to collect data: '%s'", err.Error())
			return nil, errors.New("failed to collect data")
		}
		scrapedData[dataCollector.Name] = data
	}
	return scrapedData, nil
}

func (scraper *Scraper) Save(data map[string]interface{}) error {
	for _, storageInterface := range scraper.StorageInterfaces {
		err := storageInterface.Save(data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (scraper *Scraper) Download(url string) ([]byte, error) {
	return scraper.DownloadClient.Get(url)
}

func NewSraper(config map[string]interface{}) (scraper Scraper, err error) {
	scraper.Name = config["name"].(string)
	for _, url := range config["urls"].([]interface{}) {
		scraper.URLs = append(scraper.URLs, url.(string))
	}
	tmp, err := download_client.NewDownloadClient(config["downloadClient"].(map[string]interface{}))
	if err != nil {
		log.Errorf("failed to create DownloadClient: '%s'", err.Error())
		return
	}
	scraper.DownloadClient = tmp
	err = scraper.DownloadClient.Open()
	if err != nil {
		log.Errorf("failed to open DownloadClient: '%s'", err.Error())
		return
	}
	scraper.Retrys = int64(config["retrys"].(int))
	scraper.OnFailure = config["onFailure"].(string)
	scraper.Delay = int64(config["delay"].(int))
	scraper.Workers = int64(config["workers"].(int))
	scraper.DataCollectors, err = data_collectors.NewDataCollectors(
		scraper.DataCollectorsRegister,
		config["dataCollectors"].([]interface{}),
	)
	if err != nil {
		log.Errorf("failed to create DataCollectors: '%s'", err)
		return
	}
	scraper.StorageInterfaces, err = storage_interfaces.NewStorageInterfaces(
		config["storageInterfaces"].([]interface{}),
	)
	if err != nil {
		log.Errorf("failed to create StorageInterfaces: '%s'", err)
		return
	}
	for _, storageInterface := range scraper.StorageInterfaces {
		err = storageInterface.Open()
		if err != nil {
			log.Errorf("failed to open StorageInterface: '%s'", err)
			return
		}
	}
	return
}
