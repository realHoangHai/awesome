package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/realHoangHai/awesome/pkg/log"
	"github.com/realHoangHai/awesomeProject/client"
	"github.com/realHoangHai/awesomeProject/config"
	pb "github.com/realHoangHai/awesomeProject/examples/helloworld/helloworld"
	"github.com/realHoangHai/awesomeProject/health"
	log2 "github.com/realHoangHai/awesomeProject/pkg/log"
	"google.golang.org/grpc/health/grpc_health_v1"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	correlationID := flag.String("correlation_id", "", "Use local correlation-id")
	name := flag.String("name", "Jack", "Name for greeting")
	flag.Parse()

	log.Init(log2.FromEnv(config.WithFileNoError(".env")))

	conf := client.ReadConfigFromEnv(config.WithFileNoError(".env"))
	conn := client.Must(client.Dial("", client.DialOptionsFromConfig(conf)...))
	c := pb.NewGreeterClient(conn)
	rep, err := c.SayHello(client.NewTracingContext(context.Background(), *correlationID), &pb.HelloRequest{
		Name: *name,
	})
	if err != nil {
		log2.Panic(err)
	}
	log2.Info("RESPONSE GRPC:", rep.Message)

	// Health check.
	rs, err := health.NewClient(conn).Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		log2.Panic("health check failed", err)
	}
	if rs.Status != health.StatusServing {
		log2.Panicf("got health status=%d, want status=%d", rs.Status, health.StatusServing)
	}
	log2.Info("HEALTH CHECK GRPC:", rs.Status)

	// HTTP
	host := fmt.Sprintf("http://%s", conf.Address)
	if conf.TLSCertFile != "" {
		host = fmt.Sprintf("https://%s", conf.Address)
	}
	body := bytes.NewBuffer([]byte(fmt.Sprintf(`{"name":"%s"}`, *name)))
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/hello", host), body)
	if err != nil {
		log2.Panic(err)
	}
	if *correlationID != "" {
		req.Header.Set("X-Correlation-Id", *correlationID)
	}
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	if conf.TLSCertFile != "" {
		caCert, err := ioutil.ReadFile(conf.TLSCertFile)
		if err != nil {
			log2.Panic(err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		}
	}
	res, err := httpClient.Do(req)
	if err != nil {
		log2.Panic(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log2.Panicf("got status_code=%d, want status_code=%d", res.StatusCode, http.StatusOK)
	}
	v := &pb.HelloReply{}
	if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
		log2.Panic(err)
	}
	log2.Info("RESPONSE HTTP:", v.Message)

	// internal apis
	log2.Info("HEALTH CHECK HTTP:", getString(httpClient, fmt.Sprintf("%s/internal/health", host)))
	log2.Infof("METRICS: \n%s\n", getString(httpClient, fmt.Sprintf("%s/internal/metrics", host)))
}

func getString(client *http.Client, url string) string {
	res, err := client.Get(url)
	if err != nil {
		log2.Panic(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log2.Panic("status code not ok, status_code=", res.StatusCode)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log2.Panic(err)
	}
	return string(b)
}
