package testdb

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ShatteredRealms/go-backend/pkg/config"
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
	Username = "postgres"
	Password = "password"
	DbName   = "test"
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
		Repository: "quay.io/keycloak/keycloak",
		Tag:        "21.0.0",
		Env:        []string{"KEYCLOAK_ADMIN=admin", "KEYCLOAK_ADMIN_PASSWORD=admin"},
		Cmd: []string{
			"start-dev",
			"--import-realm",
			"--health-enabled=true",
			"--features=declarative-user-profile",
		},
		Mounts:       []string{fmt.Sprintf("%s:/opt/keycloak/data/import/realm-export.json", realmExportFile)},
		ExposedPorts: []string{"8080/tcp"},
	}

	keycloakResource, err := pool.RunWithOptions(keycloakRunDockerOpts, fnConfig)
	chk(err)

	// Uncomment to see docker log
	// go func() {
	// 	pool.Client.Logs(docker.LogsOptions{
	// 		Container:    keycloakResource.Container.ID,
	// 		OutputStream: log.Logger.Out,
	// 		ErrorStream:  log.Logger.Out,
	// 		Follow:       true,
	// 		Stdout:       true,
	// 		Stderr:       true,
	// 		Timestamps:   false,
	// 		RawTerminal:  true,
	// 	})
	// }()

	closeFunc := func() {
		chk(keycloakResource.Close())
	}

	host := "http://127.0.0.1:" + keycloakResource.GetPort("8080/tcp")
	err = Retry(func() error {
		_, err := http.Get(host + "/health/ready")
		return err
	}, time.Second*60)

	err = Retry(func() error {
		_, err = http.Get(host + "/realms/default")
		return err
	}, time.Second*60)
	chk(err)

	return closeFunc, host
}

func SetupKafkaWithDocker() (func(), string) {
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
			fmt.Sprintf("KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://127.0.0.1:%s,PLAINTEXT_HOST://127.0.0.1:%s", kafkaBrokerPort, kafkaPort),
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

	return fnCleanup, strconv.Itoa(kafkaPortUint)
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
	}, time.Second*30))

	return mdb
}

func SetupGormWithDocker() (func(), string) {
	pool, err := dockertest.NewPool("")
	chk(err)

	runDockerOpt := &dockertest.RunOptions{
		Repository: "postgres", // image
		Tag:        "14",       // version
		Env:        []string{"POSTGRES_PASSWORD=" + Password, "POSTGRES_DB=" + DbName},
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

	// container is ready, return *gorm.Db for testing
	return fnCleanup, resource.GetPort("5432/tcp")
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
	}, time.Second*10))

	// container is ready, return *gorm.Db for testing
	return gdb
}

func SetupRedisWithDocker() (fnCleanup func(), redisPoolConfig config.DBPoolConfig) {
	pool, err := dockertest.NewPool("")
	chk(err)

	runDockerOpt := &dockertest.RunOptions{
		Repository: "grokzen/redis-cluster", // image
		Tag:        "7.0.10",                // version
		Env: []string{
			"INITIAL_PORT=7000",
			"MASTERS=3",
			"SLAVES_PER_MASTER=1",
		},
	}

	fnConfig := func(config *docker.HostConfig) {
		config.AutoRemove = true                     // set AutoRemove to true so that stopped container goes away by itself
		config.RestartPolicy = docker.NeverRestart() // don't restart container
	}

	resource, err := pool.RunWithOptions(runDockerOpt, fnConfig)
	chk(err)
	// call clean up function to release resource
	fnCleanup = func() {
		err := resource.Close()
		chk(err)
	}

	// container is ready, return *gorm.Db for testing
	return fnCleanup, config.DBPoolConfig{
		Master: config.DBConfig{
			ServerAddress: config.ServerAddress{
				Port: resource.GetPort("7000/tcp"),
				Host: "localhost",
			},
		},
		Slaves: []config.DBConfig{
			{
				ServerAddress: config.ServerAddress{
					Port: resource.GetPort("7001/tcp"),
					Host: "localhost",
				},
			},
			{
				ServerAddress: config.ServerAddress{
					Port: resource.GetPort("7002/tcp"),
					Host: "localhost",
				},
			},
			{
				ServerAddress: config.ServerAddress{
					Port: resource.GetPort("7003/tcp"),
					Host: "localhost",
				},
			},
			{
				ServerAddress: config.ServerAddress{
					Port: resource.GetPort("7004/tcp"),
					Host: "localhost",
				},
			},
			{
				ServerAddress: config.ServerAddress{
					Port: resource.GetPort("7005/tcp"),
					Host: "localhost",
				},
			},
		},
	}
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func Retry(op func() error, maxTime time.Duration) error {
	bo := backoff.NewExponentialBackOff()
	bo.MaxInterval = time.Second
	bo.MaxElapsedTime = maxTime
	if err := backoff.Retry(op, bo); err != nil {
		if bo.NextBackOff() == backoff.Stop {
			return errors.Wrap(err, "reached retry deadline")
		}

		return err
	}

	return nil
}
