package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"
)

func main() {
	res, err := http.Get("https://api.ipify.org")

	if err != nil {
		log.Fatal(err)
	}

	ip, _ := ioutil.ReadAll(res.Body)

	ipAddr := net.ParseIP(string(ip))

	if ipAddr == nil {
		log.Fatal(err)
	}

	awsSession, err := session.NewSession()

	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	svc := route53.New(awsSession)

	var workGroup sync.WaitGroup

	err = svc.ListHostedZonesPagesWithContext(ctx, &route53.ListHostedZonesInput{}, func(output *route53.ListHostedZonesOutput, lastPage bool) bool {
		for _, v := range output.HostedZones {
			workGroup.Add(1)
			go updateHostedZone(&ctx, svc, v, &workGroup, &ipAddr)
		}
		return lastPage
	})

	if err != nil {
		panic(err)
	}

	workGroup.Wait()
}

func updateHostedZone(ctx *context.Context, svc *route53.Route53, hostedZone *route53.HostedZone, group *sync.WaitGroup, newIP *net.IP) {
	var changeBatch = &route53.ChangeBatch{}

	svc.ListResourceRecordSetsPagesWithContext(*ctx, &route53.ListResourceRecordSetsInput{
		HostedZoneId:          hostedZone.Id,
	}, func(records *route53.ListResourceRecordSetsOutput, lastPage bool) bool {
		for _, v := range records.ResourceRecordSets {
			if *v.Type != "A" {
				continue
			}

			//ignore multi values
			if len(v.ResourceRecords) != 1  {
				continue
			}

			record := v.ResourceRecords[0]
			recordIp := net.ParseIP(*record.Value)

			// if not valid ip skip
			if recordIp == nil {
				continue
			}

			// if same ip skip
			if newIP.Equal(recordIp) {
				continue
			}


			record.SetValue(newIP.String())

			changeBatch.Changes = append(changeBatch.Changes, &route53.Change{
				Action:            aws.String("UPSERT"),
				ResourceRecordSet: v,
			})
		}

		if lastPage {
			if len(changeBatch.Changes) > 0 {
				fmt.Println("IP Change found updating records...")
				_, err := svc.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
					ChangeBatch:  changeBatch,
					HostedZoneId: hostedZone.Id,
				})

				if err != nil {
					log.Fatal(err)
				}
			} else {
				fmt.Println("No changes needed finishing....")
			}

			group.Done()
		}

		return lastPage
	})
}
