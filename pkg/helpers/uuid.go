package helpers

import (
	"fmt"

	"github.com/google/uuid"
)

func ParseUUIDs(stringIds []string) ([]*uuid.UUID, error) {
	ids := make([]*uuid.UUID, len(stringIds))
	for idx, stringId := range stringIds {
		id, err := uuid.Parse(stringId)
		if err != nil {
			return nil, fmt.Errorf("invalid id: %s", stringId)
		}
		ids[idx] = &id
	}

	return ids, nil
}
