package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/domainr/whois"
	"gopkg.in/yaml.v2"
)

type TLDConfig struct {
	Expiration       string `yaml:"expiration"`
	ExpirationFormat string `yaml:"expiration_format"`
	Availability     string `yaml:"availability"`
}

type DomainInfo struct {
	Domain      string
	Expiration  time.Time
	IsAvailable bool
	FullWhois   string
}

var tldConfigs map[string]TLDConfig

func init() {
	var err error

	tldConfigs, err = loadTLDConfigs()
	if err != nil {
		panic(err)
	}
}

func QueryDomainInfo(domain string) (*DomainInfo, error) {
	req, err := whois.NewRequest(domain)

	if err != nil {
		return nil, err
	}

	res, err := whois.DefaultClient.Fetch(req)

	if err != nil {
		return nil, err
	}

	body := string(res.Body)
	info := &DomainInfo{}

	suffix := strings.TrimPrefix(filepath.Ext(domain), ".")

	tld, ok := tldConfigs[strings.ToLower(suffix)]
	if !ok {
		err = fmt.Errorf("unsupported TLD")
		return nil, err
	}

	info, err = parseTLDInfo(body, tld)
	if err != nil {
		return nil, err
	}

	info.Domain = domain
	return info, nil
}

func loadTLDConfigs() (map[string]TLDConfig, error) {
	tldConfigs := make(map[string]TLDConfig)

	// Read the YAML files in the tlds directory
	dirContents, err := filepath.Glob("tlds/*.yml")
	if err != nil {
		return nil, err
	}

	// Load the TLD configuration from each YAML file
	for _, filename := range dirContents {
		configFile, err := readFile(filename)
		if err != nil {
			return nil, err
		}

		var tldConfig TLDConfig
		err = yaml.Unmarshal(configFile, &tldConfig)
		if err != nil {
			return nil, err
		}

		suffix := strings.TrimSuffix(filepath.Base(filename), ".yml")

		tldConfigs[suffix] = tldConfig
	}

	return tldConfigs, nil
}

func readFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	data := make([]byte, fileInfo.Size())
	_, err = file.Read(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func parseTLDInfo(body string, config TLDConfig) (*DomainInfo, error) {
	info := &DomainInfo{
		FullWhois: body,
	}

	statusRegex := regexp.MustCompile(config.Availability)
	statusMatch := statusRegex.FindString(body)

	if statusMatch != "" {
		info.IsAvailable = true
		return info, nil
	}

	expirationRegex := regexp.MustCompile(config.Expiration)
	expirationMatch := expirationRegex.FindStringSubmatch(body)

	if len(expirationMatch) > 1 {
		expiration := expirationMatch[1]

		expirationTime, err := time.Parse(config.ExpirationFormat, expiration)
		if err != nil {
			fmt.Println("error parsing datetime:", err)
			return info, nil
		}

		info.Expiration = expirationTime
		return info, nil
	}

	return info, nil
}
