package game_test

import (
	"testing"
	"time"

	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/model/game"
	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

func TestModel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Game Model Suite")
}

func randomDimensionAndMap() (*game.Dimension, *game.Map) {
	uuid1, err := uuid.NewRandom()
	Expect(err).NotTo(HaveOccurred())
	uuid2, err := uuid.NewRandom()
	Expect(err).NotTo(HaveOccurred())
	m := &game.Map{
		Model: model.Model{
			Id:        &uuid2,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: gorm.DeletedAt{},
		},
		Name:       faker.Username(),
		Path:       faker.Username(),
		MaxPlayers: 40,
		Instanced:  false,
	}
	dimension := &game.Dimension{
		Model: model.Model{
			Id:        &uuid1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: gorm.DeletedAt{},
		},
		Name:     faker.Username(),
		Location: "us-central",
		Version:  faker.Username(),
		Maps:     []*game.Map{m},
	}
	m.Dimensions = []*game.Dimension{dimension}

	return dimension, m
}
