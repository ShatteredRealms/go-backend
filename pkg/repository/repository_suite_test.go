package repository_test

import (
	"context"
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
	BeforeSuite(func() {
		log.Logger, hook = test.NewNullLogger()

		gdb, gdbCloseFunc = testdb.SetupGormWithDocker()
		Expect(gdb).NotTo(BeNil())

		characterRepo = repository.NewCharacterRepository(gdb)
		Expect(characterRepo).NotTo(BeNil())
		Expect(characterRepo.Migrate(context.Background())).To(Succeed())

		chatRepo = repository.NewChatRepository(gdb)
		Expect(chatRepo).NotTo(BeNil())
		Expect(chatRepo.Migrate(context.Background())).To(Succeed())

		gamebackendRepo = repository.NewGamebackendRepository(gdb)
		Expect(gamebackendRepo).NotTo(BeNil())
		Expect(gamebackendRepo.Migrate(context.Background())).To(Succeed())

		mdb, mdbCloseFunc = testdb.SetupMongoWithDocker()
		Expect(mdb).NotTo(BeNil())

		invRepo = repository.NewInventoryRepository(mdb)
		Expect(invRepo).NotTo(BeNil())
	})

	BeforeEach(func() {
		log.Logger, hook = test.NewNullLogger()
	})

	AfterSuite(func() {
		gdbCloseFunc()
		mdbCloseFunc()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Repository Suite")
}
