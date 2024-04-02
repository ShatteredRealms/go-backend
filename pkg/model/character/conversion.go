package character

import (
	"context"
	"fmt"

	"github.com/ShatteredRealms/go-backend/pkg/common"
	"github.com/ShatteredRealms/go-backend/pkg/log"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
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
		return 0, common.ErrHandleRequest.Err()
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
		return "", common.ErrHandleRequest.Err()

	}

	return targetCharacterName, nil
}
