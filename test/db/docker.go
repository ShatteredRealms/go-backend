package testdb

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/cenkalti/backoff/v4"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	username = "root"
	password = "test"
	dbName   = "test"
)

var (
	portOffset = 1
)

func SetupKeycloakWithDocker() (func(), string) {
	pool, err := dockertest.NewPool("")
	chk(err)

	fnConfig := func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.NeverRestart()
	}

	wd, err := os.Getwd()
	chk(err)
	realmExportFile, err := filepath.Abs(fmt.Sprintf("%s/../../test/db/keycloak-realm-export.json", wd))
	chk(err)
	keycloakRunDockerOpts := &dockertest.RunOptions{
		Repository:   "quay.io/keycloak/keycloak",
		Tag:          "latest",
		Env:          []string{"KEYCLOAK_ADMIN=admin", "KEYCLOAK_ADMIN_PASSWORD=admin"},
		Cmd:          []string{"start-dev --import-realm --features=declarative-user-profile"},
		Mounts:       []string{fmt.Sprintf("%s:/opt/keycloak/data/import/realm-export.json", realmExportFile)},
		ExposedPorts: []string{"8080/tcp"},
	}

	keycloakResource, err := pool.RunWithOptions(keycloakRunDockerOpts, fnConfig)
	chk(err)

	closeFunc := func() {
		chk(keycloakResource.Close())
	}

	host := fmt.Sprintf("http://%s", keycloakResource.GetHostPort("8080/tcp"))
	return closeFunc, host
}

func ConnectKeycloakDocker(host string) *gocloak.GoCloak {

	client := gocloak.NewClient(host)

	chk(Retry(func() error {
		_, err := http.Get(host + "/realms/default")
		return err
	}))

	return client
}

func SetupKafkaWithDocker() (func(), uint) {
	pool, err := dockertest.NewPool("")
	chk(err)

	net, err := pool.CreateNetwork(fmt.Sprintf("go-testing-%d", time.Now().UnixMilli()))
	chk(err)

	zooKeeperPort := strconv.Itoa(2182 + portOffset)
	kafkaPortUint := 29093 + portOffset
	kafkaPort := strconv.Itoa(kafkaPortUint)
	kafkaBrokerPort := strconv.Itoa(9093 + portOffset)
	portOffset++

	zookeeperRunDockerOpts := &dockertest.RunOptions{
		Hostname:   "gozookeeper",
		Repository: "confluentinc/cp-zookeeper",
		Tag:        "latest",
		Env:        []string{"ZOOKEEPER_CLIENT_PORT=" + zooKeeperPort},
		PortBindings: map[docker.Port][]docker.PortBinding{
			docker.Port(zooKeeperPort + "/tcp"): {{HostIP: "gozookeeper", HostPort: zooKeeperPort + "/tcp"}},
		},
		ExposedPorts: []string{zooKeeperPort + "/tcp", zooKeeperPort + "/tcp"},
		Networks:     []*dockertest.Network{net},
	}

	kafkaRunDockerOpts := &dockertest.RunOptions{
		Hostname:   "gokafka",
		Repository: "confluentinc/cp-kafka",
		Tag:        "latest",
		Env: []string{
			"KAFKA_BROKER_ID=1",
			"KAFKA_ZOOKEEPER_CONNECT=gozookeeper:" + zooKeeperPort,
			fmt.Sprintf("KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://localhost:%s,PLAINTEXT_HOST://localhost:%s", kafkaBrokerPort, kafkaPort),
			"KAFKA_LISTENER_SECURITY_PROTOCOL_MAP=PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT",
			"KAFKA_INTER_BROKER_LISTENER_NAME=PLAINTEXT",
		},
		PortBindings: map[docker.Port][]docker.PortBinding{
			docker.Port(kafkaPort + "/tcp"):       {{HostIP: "localhost", HostPort: kafkaPort + "/tcp"}},
			docker.Port(kafkaBrokerPort + "/tcp"): {{HostIP: "localhost", HostPort: kafkaBrokerPort + "/tcp"}},
		},
		ExposedPorts: []string{kafkaPort + "/tcp", kafkaBrokerPort + "/tcp"},
		Networks:     []*dockertest.Network{net},
	}

	fnConfig := func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.NeverRestart()
	}

	var zookeeperResource, kafkaResource *dockertest.Resource
	chk(err)

	chk(pool.Retry(func() error {
		zookeeperResource, err = pool.RunWithOptions(zookeeperRunDockerOpts, fnConfig)
		return err
	}))

	chk(pool.Retry(func() error {
		kafkaResource, err = pool.RunWithOptions(kafkaRunDockerOpts, fnConfig)
		return err
	}))

	fnCleanup := func() {
		err1 := kafkaResource.Close()
		err2 := zookeeperResource.Close()
		err3 := net.Close()
		chk(err1)
		chk(err2)
		chk(err3)
	}

	return fnCleanup, uint(kafkaPortUint)
}

func SetupMongoWithDocker() (func(), string) {
	pool, err := dockertest.NewPool("")
	chk(err)

	runDockerOpt := &dockertest.RunOptions{
		Repository: "mongo", // image
		Tag:        "6.0",   // version
		Env:        []string{"MONGO_INITDB_ROOT_USERNAME=root", "MONGO_INITDB_ROOT_PASSWORD=password"},
	}

	fnConfig := func(config *docker.HostConfig) {
		config.AutoRemove = true                     // set AutoRemove to true so that stopped container goes away by itself
		config.RestartPolicy = docker.NeverRestart() // don't restart container
	}

	resource, err := pool.RunWithOptions(runDockerOpt, fnConfig)
	chk(err)
	// call clean up function to release resource
	fnCleanup := func() {
		err := resource.Close()
		chk(err)
	}

	return fnCleanup, fmt.Sprintf("mongodb://root:password@localhost:%s", resource.GetPort("27017/tcp"))
}

func ConnectMongoDocker(host string) *mongo.Database {
	var mdb *mongo.Database

	chk(Retry(func() error {
		db, err := mongo.Connect(
			context.TODO(),
			options.Client().ApplyURI(
				host,
			),
		)
		if err != nil {
			return err
		}
		mdb = db.Database("testdb")
		return db.Ping(context.TODO(), nil)
	}))

	return mdb
}

func SetupGormWithDocker() (func(), string) {
	pool, err := dockertest.NewPool("")
	chk(err)

	runDockerOpt := &dockertest.RunOptions{
		Repository: "postgres", // image
		Tag:        "14",       // version
		Env:        []string{"POSTGRES_PASSWORD=" + password, "POSTGRES_DB=" + dbName},
	}

	fnConfig := func(config *docker.HostConfig) {
		config.AutoRemove = true                     // set AutoRemove to true so that stopped container goes away by itself
		config.RestartPolicy = docker.NeverRestart() // don't restart container
	}

	resource, err := pool.RunWithOptions(runDockerOpt, fnConfig)
	chk(err)
	// call clean up function to release resource
	fnCleanup := func() {
		err := resource.Close()
		chk(err)
	}

	connStr := fmt.Sprintf("host=localhost port=%s user=postgres dbname=%s password=%s sslmode=disable",
		resource.GetPort("5432/tcp"), // get port of localhost
		dbName,
		password,
	)

	// container is ready, return *gorm.Db for testing
	return fnCleanup, connStr
}

func ConnectGormDocker(connStr string) *gorm.DB {
	var gdb *gorm.DB
	// retry until db server is ready
	chk(Retry(func() (err error) {
		gdb, err = gorm.Open(postgres.Open(connStr), &gorm.Config{
			Logger: logger.New(
				log.Logger,
				logger.Config{
					SlowThreshold:             0,
					Colorful:                  true,
					IgnoreRecordNotFoundError: true,
					ParameterizedQueries:      true,
					LogLevel:                  logger.Info,
				},
			),
		})
		if err != nil {
			return err
		}
		db, err := gdb.DB()
		if err != nil {
			return err
		}
		return db.Ping()
	}))

	// container is ready, return *gorm.Db for testing
	return gdb
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func Retry(op func() error) error {
	bo := backoff.NewExponentialBackOff()
	bo.MaxInterval = time.Second
	bo.MaxElapsedTime = time.Second * 10
	if err := backoff.Retry(op, bo); err != nil {
		if bo.NextBackOff() == backoff.Stop {
			return errors.Wrap(err, "reached retry deadline")
		}

		return err
	}

	return nil
}
