package spider

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
)

type TSpiderConfig struct {
	PullOnStartup bool `xml:"task>startup"`
	Schedule string `xml:"task>schedule"`
	RedisServerAddress string `xml:"redis_server"`
}

var GTSpiderConfig *TSpiderConfig

func ParseXmlConfig(path string) (*TSpiderConfig, error) {
	if len(path) == 0 {
		return nil, errors.New("not found configure xml file")
	}

	n, err := GetFileSize(path)
	if err != nil || n == 0 {
		return nil, errors.New("not found configure xml file")
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	GTSpiderConfig = &TSpiderConfig{}

	data := make([]byte, n)

	m, err := f.Read(data)
	if err != nil {
		return nil, err
	}

	if int64(m) != n {
		return nil, errors.New(fmt.Sprintf("expect read configure xml file size %d but result is %d", n, m))
	}

	err = xml.Unmarshal(data, &GTSpiderConfig)
	if err != nil {
		return nil, err
	}

	Logger.Printf("read config %+v", GTSpiderConfig)

	return GTSpiderConfig, nil
}
