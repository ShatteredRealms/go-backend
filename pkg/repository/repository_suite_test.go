package repository_test

import (
	"bytes"
	"context"
	"encoding/gob"
	"testing"

	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
	testdb "github.com/ShatteredRealms/go-backend/test/db"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus/hooks/test"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type initializeData struct {
	GormConfig  config.DBConfig
	MdbConnStr  string
	RedisConfig config.DBPoolConfig
}

var (
	hook *test.Hook

	gdb            *gorm.DB
	gdbCloseFunc   func()
	mdb            *mongo.Database
	mdbCloseFunc   func()
	redisCloseFunc func()

	data initializeData

	characterRepo   repository.CharacterRepository
	chatRepo        repository.ChatRepository
	gamebackendRepo repository.GamebackendRepository
	invRepo         repository.InventoryRepository
)

func TestRepository(t *testing.T) {
	SynchronizedBeforeSuite(func() []byte {
		log.Logger, hook = test.NewNullLogger()

		var gormPort string
		gdbCloseFunc, gormPort = testdb.SetupGormWithDocker()
		Expect(gdbCloseFunc).NotTo(BeNil())

		mdbCloseFunc, data.MdbConnStr = testdb.SetupMongoWithDocker()
		Expect(mdbCloseFunc).NotTo(BeNil())

		redisCloseFunc, data.RedisConfig = testdb.SetupRedisWithDocker()

		data.GormConfig = config.DBConfig{
			ServerAddress: config.ServerAddress{
				Port: gormPort,
				Host: "localhost",
			},
			Name:     testdb.DbName,
			Username: testdb.Username,
			Password: testdb.Password,
		}
		gdb = testdb.ConnectGormDocker(data.GormConfig.PostgresDSN())
		Expect(gdb).NotTo(BeNil())
		mdb = testdb.ConnectMongoDocker(data.MdbConnStr)
		Expect(mdb).NotTo(BeNil())

		var err error
		characterRepo, err = repository.NewCharacterRepository(gdb)
		Expect(err).NotTo(HaveOccurred())
		Expect(characterRepo).NotTo(BeNil())
		Expect(characterRepo.Migrate(context.Background())).NotTo(HaveOccurred())

		chatRepo = repository.NewChatRepository(gdb)
		Expect(chatRepo).NotTo(BeNil())
		Expect(chatRepo.Migrate(context.Background())).NotTo(HaveOccurred())

		gamebackendRepo = repository.NewGamebackendRepository(gdb)
		Expect(gamebackendRepo).NotTo(BeNil())
		Expect(gamebackendRepo.Migrate(context.Background())).NotTo(HaveOccurred())

		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		Expect(enc.Encode(data)).To(Succeed())

		return buf.Bytes()
	}, func(inBytes []byte) {
		log.Logger, hook = test.NewNullLogger()

		dec := gob.NewDecoder(bytes.NewBuffer(inBytes))
		Expect(dec.Decode(&data)).To(Succeed())

		gdb = testdb.ConnectGormDocker(data.GormConfig.PostgresDSN())
		Expect(gdb).NotTo(BeNil())
		mdb = testdb.ConnectMongoDocker(data.MdbConnStr)
		Expect(mdb).NotTo(BeNil())

		var err error
		characterRepo, err = repository.NewCharacterRepository(gdb)
		Expect(err).NotTo(HaveOccurred())
		Expect(characterRepo).NotTo(BeNil())
		chatRepo = repository.NewChatRepository(gdb)
		Expect(chatRepo).NotTo(BeNil())
		gamebackendRepo = repository.NewGamebackendRepository(gdb)
		Expect(gamebackendRepo).NotTo(BeNil())
		invRepo = repository.NewInventoryRepository(mdb)
		Expect(invRepo).NotTo(BeNil())
	})

	BeforeEach(func() {
		log.Logger, hook = test.NewNullLogger()
	})

	SynchronizedAfterSuite(func() {
	}, func() {
		gdbCloseFunc()
		mdbCloseFunc()
		redisCloseFunc()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Repository Suite")
}
