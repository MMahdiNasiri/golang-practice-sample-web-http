package redis

import (
	"os"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	ClusterClient *redis.ClusterClient
}

func New() (*Client, error) {
	//masterPort := os.Getenv("MASTER_REDIS_PORT")
	//if masterPort == "" {
	//	masterPort = "6379"
	//}
	//masterRoute := os.Getenv("MASTER_REDIS_ROUTE")
	//
	//masterClient := redis.NewClient(&redis.Options{
	//	Addr:     fmt.Sprintf("%s:%s", masterRoute, masterPort),
	//	Password: "",
	//	DB:       0,
	//})
	//var ctx = context.Background()
	//
	//_, err := masterClient.Ping(ctx).Result()
	//if err != nil {
	//	return nil, err
	//}
	//
	//replicaPort := os.Getenv("REPLICA_REDIS_PORT")
	//if replicaPort == "" {
	//	replicaPort = "6379"
	//}
	//replicaRoute := os.Getenv("REPLICA_REDIS_ROUTE")
	//
	//replicaClient := redis.NewClient(&redis.Options{
	//	Addr:     fmt.Sprintf("%s:%s", replicaRoute, replicaPort),
	//	Password: "",
	//	DB:       0,
	//})
	//
	//_, err = replicaClient.Ping(ctx).Result()
	//
	//if err != nil {
	//	return nil, err
	//}
	//client := &Client{
	//	Master:  masterClient,
	//	Replica: replicaClient,
	//}

	clusterClient := redis.NewFailoverClusterClient(
		&redis.FailoverOptions{
			MasterName:     os.Getenv("REDIS_MASTER_NAME"),
			SentinelAddrs:  []string{"sentinel-1:26379", "sentinel-2:26379", "sentinel-3:26379"},
			DB:             0,
			RouteByLatency: true,
		})

	client := &Client{
		ClusterClient: clusterClient,
	}

	return client, nil
}

func (client *Client) Close() error {
	return client.ClusterClient.Close()
}
