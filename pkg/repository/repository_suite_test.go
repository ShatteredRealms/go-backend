package repository_test

import (
	"context"
	"strings"
	"testing"

	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/repository"
	testdb "github.com/ShatteredRealms/go-backend/test/db"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus/hooks/test"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var (
	hook *test.Hook

	gdb          *gorm.DB
	gdbCloseFunc func()

	mdb          *mongo.Database
	mdbCloseFunc func()

	characterRepo   repository.CharacterRepository
	chatRepo        repository.ChatRepository
	gamebackendRepo repository.GamebackendRepository
	invRepo         repository.InventoryRepository
)

func TestRepository(t *testing.T) {
	splitter := "\n"
	SynchronizedBeforeSuite(func() []byte {
		log.Logger, hook = test.NewNullLogger()
		var gdbConnStr, mdbConnStr string

		gdbCloseFunc, gdbConnStr = testdb.SetupGormWithDocker()
		Expect(gdbCloseFunc).NotTo(BeNil())

		mdbCloseFunc, mdbConnStr = testdb.SetupMongoWithDocker()
		Expect(mdbCloseFunc).NotTo(BeNil())

		gdb = testdb.ConnectGormDocker(gdbConnStr)
		Expect(gdb).NotTo(BeNil())
		mdb = testdb.ConnectMongoDocker(mdbConnStr)
		Expect(mdb).NotTo(BeNil())

		characterRepo = repository.NewCharacterRepository(gdb)
		Expect(characterRepo).NotTo(BeNil())
		Expect(characterRepo.Migrate(context.Background())).NotTo(HaveOccurred())

		chatRepo = repository.NewChatRepository(gdb)
		Expect(chatRepo).NotTo(BeNil())
		Expect(chatRepo.Migrate(context.Background())).NotTo(HaveOccurred())

		gamebackendRepo = repository.NewGamebackendRepository(gdb)
		Expect(gamebackendRepo).NotTo(BeNil())
		Expect(gamebackendRepo.Migrate(context.Background())).NotTo(HaveOccurred())

		return []byte(gdbConnStr + splitter + mdbConnStr)
	}, func(hostsBytes []byte) {
		log.Logger, hook = test.NewNullLogger()

		hosts := strings.Split(string(hostsBytes), splitter)
		Expect(hosts).To(HaveLen(2))
		gdb = testdb.ConnectGormDocker(hosts[0])
		Expect(gdb).NotTo(BeNil())
		mdb = testdb.ConnectMongoDocker(hosts[1])
		Expect(mdb).NotTo(BeNil())

		characterRepo = repository.NewCharacterRepository(gdb)
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
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Repository Suite")
}
