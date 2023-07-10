package model

import "github.com/ShatteredRealms/go-backend/pkg/pb"

type ChatTemplate struct {
	Model
	Name string `gorm:"unique"`
}

type ChatTemplates []*ChatTemplate

func (c *ChatTemplate) ToPb() *pb.ChatTemplate {
	return &pb.ChatTemplate{
		Id:   c.Id.String(),
		Name: c.Name,
	}
}

func (chatTemplates ChatTemplates) ToPb() []*pb.ChatTemplate {
	out := make([]*pb.ChatTemplate, len(chatTemplates))
	for idx, c := range chatTemplates {
		out[idx] = c.ToPb()
	}

	return out
}
