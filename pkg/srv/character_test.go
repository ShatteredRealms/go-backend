package srv_test

import (
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/model"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/pb"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/srv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Character", func() {
	var (
		serviceSrv pb.CharactersServiceServer
	)

	BeforeEach(func() {
		serviceSrv = srv.NewCharacterServiceServer(newServiceMock(), nil)
		Expect(serviceSrv).NotTo(BeNil())
	})
})

type serviceMock struct {
}

// AddPlayTime implements service.CharacterService
func (serviceMock) AddPlayTime(characterId uint64, amount uint64) (uint64, error) {
	return 1, nil
}

// Create implements service.CharacterService
func (serviceMock) Create(ownerId uint64, name string, genderId uint64, realmId uint64) (*model.Character, error) {
	return nil, nil
}

// Delete implements service.CharacterService
func (serviceMock) Delete(id uint64) error {
	return nil
}

// Edit implements service.CharacterService
func (serviceMock) Edit(character *pb.Character) (*model.Character, error) {
	return nil, nil
}

// FindAll implements service.CharacterService
func (serviceMock) FindAll() ([]*model.Character, error) {
	return nil, nil
}

// FindAllByOwner implements service.CharacterService
func (serviceMock) FindAllByOwner(ownerId uint64) ([]*model.Character, error) {
	return nil, nil
}

// FindById implements service.CharacterService
func (serviceMock) FindById(id uint64) (*model.Character, error) {
	return nil, nil
}

func newServiceMock() service.CharacterService {
	return serviceMock{}
}
