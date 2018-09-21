package main

import (
	"fmt"
	"regexp"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"strconv"
	"time"
	"github.com/prometheus/common/log"
)

func extractErrorRate(reader io.Reader, config HTTPProbe) int {
	var re = regexp.MustCompile(`(\d+)]]$`)
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Errorf("Error reading HTTP body: %s", err)
		return 0
	}
	var str = string(body)
	matches := re.FindStringSubmatch(str)
	value, err := strconv.Atoi(matches[1])
	if err == nil {
		return value
	}
	return 0
}

func printRespBody(reader io.Reader) string {
	body, err:= ioutil.ReadAll(reader)
	if err != nil {
		return "Error reading HTTP body"
	}
	var str = string(body)
	return str
}
func probeHTTP(target string, w http.ResponseWriter, module Module) (success bool) {
	config := module.HTTP

	client := &http.Client{
		Timeout: module.Timeout,
	}
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	log.Infof(timestamp)
 	requestURL := config.Prefix + target + "/stats/?stat=received&since=1537531898&until=" + timestamp
	log.Infof(requestURL)
	log.Infof("URL should be https://sentry.io/api/0/projects/{organization}/%s", target)
	log.Infof("I believe that the endpoint we are hitting requires additional info.")
	request, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		log.Errorf("Error creating request for target %s: %s", target, err)
		return
	}

	for key, value := range config.Headers {
		if strings.Title(key) == "Host" {
			request.Host = value
			continue
		}
		request.Header.Set(key, value)
	}

	resp, err := client.Do(request)
	// Err won't be nil if redirects were turned off. See https://github.com/golang/go/issues/3795
	if err != nil && resp == nil {
		log.Warnf("Error for HTTP request to %s: %s", target, err)
	} else {
		status := strconv.Itoa(resp.StatusCode)
		log.Infof(status)
		log.Infof(printRespBody(resp.Body))
		defer resp.Body.Close()
		if len(config.ValidStatusCodes) != 0 {
			for _, code := range config.ValidStatusCodes {
				if resp.StatusCode == code {
					success = true
					break
				}
			}
		} else if 200 <= resp.StatusCode && resp.StatusCode < 300 {
			success = true
		}
		if success {
			fmt.Fprintf(w, "probe_sentry_error_received %d\n", extractErrorRate(resp.Body, config))
		}
	}
	if resp == nil {
		resp = &http.Response{}
	}

	fmt.Fprintf(w, "probe_sentry_status_code %d\n", resp.StatusCode)
	fmt.Fprintf(w, "probe_sentry_content_length %d\n", resp.ContentLength)

	return
}
