package helpers

import (
	"context"
	"fmt"
	"reflect"

	"github.com/ShatteredRealms/go-backend/pkg/model"
	"github.com/ShatteredRealms/go-backend/pkg/pb"
	log "github.com/sirupsen/logrus"
)

func GetCharacterIdFromTarget(
	ctx context.Context,
	charactersServiceClient pb.CharactersServiceClient,
	target *pb.CharacterTarget,
) (uint, error) {
	if target == nil {
		return 0, fmt.Errorf("target cannot be nil")
	}

	targetCharacterId := uint(0)
	switch t := target.Target.(type) {
	case *pb.CharacterTarget_Name:
		targetChar, err := charactersServiceClient.GetCharacter(ctx, target)
		if err != nil {
			return 0, err
		}
		targetCharacterId = uint(targetChar.Id)

	case *pb.CharacterTarget_Id:
		targetCharacterId = uint(t.Id)

	default:
		log.WithContext(ctx).Errorf("target type unknown: %s", reflect.TypeOf(target.Target).Name())
		return 0, model.ErrHandleRequest
	}

	return targetCharacterId, nil
}

func GetCharacterNameFromTarget(
	ctx context.Context,
	charactersServiceClient pb.CharactersServiceClient,
	target *pb.CharacterTarget,
) (string, error) {
	targetCharacterName := ""
	switch t := target.Target.(type) {
	case *pb.CharacterTarget_Name:
		targetCharacterName = t.Name

	case *pb.CharacterTarget_Id:
		targetChar, err := charactersServiceClient.GetCharacter(ctx, target)
		if err != nil {
			return "", err
		}
		targetCharacterName = targetChar.Name

	default:
		log.WithContext(ctx).Errorf("target type unknown: %s", reflect.TypeOf(target.Target).Name())
		return "", model.ErrHandleRequest
	}

	return targetCharacterName, nil
}
