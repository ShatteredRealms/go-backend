package model_test

import (
	"testing"
	"time"

	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

func TestModel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Model Suite")
}
func randomDimensionAndMap() (*model.Dimension, *model.Map) {
	uuid1, err := uuid.NewRandom()
	Expect(err).NotTo(HaveOccurred())
	uuid2, err := uuid.NewRandom()
	Expect(err).NotTo(HaveOccurred())
	m := &model.Map{
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
	dimension := &model.Dimension{
		Model: model.Model{
			Id:        &uuid1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: gorm.DeletedAt{},
		},
		Name:     faker.Username(),
		Location: "us-central",
		Version:  faker.Username(),
		Maps:     []*model.Map{m},
	}
	m.Dimensions = []*model.Dimension{dimension}

	return dimension, m
}
