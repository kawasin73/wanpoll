package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/pkg/errors"
)

var (
	routerUrl    = flag.String("url", "http://172.16.255.254/GetWanIP.html", "routerIpPage url which shows global wan ip address")
	user         = flag.String("user", "", "basic auth username for routerIpPage")
	password     = flag.String("password", "", "basic auth password for routerIpPage")
	interval     = flag.Int("interval", 1, "interval time polling to routerIpPage (second)")
	region       = flag.String("region", "ap-northeast-1", "aws region where route53 is placed")
	hostedZoneId = flag.String("hz", "", "hosted zone Id in route53")
	recordNames  = flag.String("name", "", "recode name separated with ',' comma")
	ttl          = flag.Int64("ttl", 60, "dns record ttl (second)")
)

func main() {
	flag.Parse()
	if *interval < 1 {
		*interval = 1
	}
	if *hostedZoneId == "" {
		log.Println("hz (hosted zone id) is required")
		os.Exit(1)
		return
	}
	records := strings.Split(*recordNames, ",")
	if len(records) == 0 {
		log.Println("record (record name) is required")
		os.Exit(1)
		return
	}
	if *ttl < 1 {
		log.Println("ttl must be greater than 0")
		os.Exit(1)
		return
	}

	duration := time.Duration(*interval) * time.Second
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(*region),
	})
	if err != nil {
		log.Println("failed to create aws session :", err)
		os.Exit(1)
		return
	}
	rt := route53.New(sess)

	var currentIp string

	for {
		ip, err := fetchIp(*routerUrl, *user, *password)
		if err != nil {
			log.Println("failed to fetch ip address :", err)
			time.Sleep(duration)
			continue
		}
		if ip == currentIp {
			time.Sleep(duration)
			continue
		}

		log.Printf("detect new ip address \"%v\" -> \"%v\"\n", currentIp, ip)

		err = setDNS(rt, *hostedZoneId, records, *ttl, ip)
		if err != nil {
			log.Println("failed to set dns :", err)
			time.Sleep(duration)
			continue
		}

		for _, r := range records {
			log.Printf("update record \"%v\" == \"%v\"\n", r, ip)
		}

		currentIp = ip
		time.Sleep(duration)
	}
}

var (
	ipMatcher     = regexp.MustCompile(`([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3})`)
	errIpNotFound = errors.New("ip address not found in response")
)

func fetchIp(url, user, password string) (string, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.Wrap(err, "build request")
	}
	req.SetBasicAuth(user, password)
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "send get request")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", errors.New(resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "read body")
	}
	ip := ipMatcher.Find(body)
	if ip == nil {
		return "", errIpNotFound
	}
	return string(ip), nil
}

func setDNS(rt *route53.Route53, zoneId string, records []string, ttl int64, ip string) error {
	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: make([]*route53.Change, len(records)),
			Comment: aws.String(fmt.Sprintf("changed at %v", time.Now())),
		},
		HostedZoneId: aws.String(zoneId),
	}
	for i, record := range records {
		params.ChangeBatch.Changes[i] = &route53.Change{
			Action: aws.String(route53.ChangeActionUpsert),
			ResourceRecordSet: &route53.ResourceRecordSet{
				Name: aws.String(record),
				ResourceRecords: []*route53.ResourceRecord{
					{
						Value: aws.String(ip),
					},
				},
				TTL:  aws.Int64(ttl),
				Type: aws.String(route53.RRTypeA),
			},
		}
	}
	_, err := rt.ChangeResourceRecordSets(params)
	if err != nil {
		return err
	}
	return nil
}
