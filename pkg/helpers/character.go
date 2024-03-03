package helpers

import (
	"context"
	"fmt"

	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	"github.com/ShatteredRealms/go-backend/pkg/log"
)

func GetCharacterIdFromTarget(
	ctx context.Context,
	charactersServiceClient pb.CharacterServiceClient,
	target *pb.CharacterTarget,
) (uint, error) {
	if target == nil {
		return 0, fmt.Errorf("target cannot be nil")
	}

	targetCharacterId := uint(0)
	switch t := target.Type.(type) {
	case *pb.CharacterTarget_Name:
		targetChar, err := charactersServiceClient.GetCharacter(ctx, target)
		if err != nil {
			return 0, err
		}
		targetCharacterId = uint(targetChar.Id)

	case *pb.CharacterTarget_Id:
		targetCharacterId = uint(t.Id)

	default:
		log.Logger.WithContext(ctx).Errorf("target type unknown: %+v", target)
		return 0, model.ErrHandleRequest
	}

	return targetCharacterId, nil
}

func GetCharacterNameFromTarget(
	ctx context.Context,
	charactersServiceClient pb.CharacterServiceClient,
	target *pb.CharacterTarget,
) (string, error) {
	if target == nil {
		return "", fmt.Errorf("target cannot be nil")
	}

	targetCharacterName := ""
	switch t := target.Type.(type) {
	case *pb.CharacterTarget_Name:
		targetCharacterName = t.Name

	case *pb.CharacterTarget_Id:
		targetChar, err := charactersServiceClient.GetCharacter(ctx, target)
		if err != nil {
			return "", err
		}
		targetCharacterName = targetChar.Name

	default:
		log.Logger.WithContext(ctx).Errorf("target type unknown: %+v", target)
		return "", model.ErrHandleRequest
	}

	return targetCharacterName, nil
}
