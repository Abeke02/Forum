package service

import (
	"fmt"

	"main.go/internal/storage"
)

type Reaction interface {
	CreateReaction(id int, object string, action string, username string) error
	GetCounts(id int, object string) (int, int, error)
}

type ReactionService struct {
	storage storage.Reaction
}

func newReactionService(storage storage.Reaction) *ReactionService {
	return &ReactionService{
		storage: storage,
	}
}

func (r *ReactionService) CreateReaction(id int, object string, action string, username string) error {
	allReactions, err := r.storage.GetReactionsById(id, object)
	if err != nil {
		return nil
	}
	if username == "" {
		return fmt.Errorf("service: reaction: username is empty")
	}
	for i := 0; i < len(allReactions); i++ {
		if allReactions[i].Username == username {
			react, _ := r.storage.GetReactionById(id, object, username)
			if react.Reaction == action {
				if err := r.storage.DeleteReaction(id, object, username); err != nil {
					return err
				}
				return nil
			}
			if err := r.storage.UpdateReaction(id, action, object, username); err != nil {
				return err
			}
			return nil
		}
	}
	if err := r.storage.CreateReaction(id, action, object, username); err != nil {
		return err
	}

	return nil
}

func (r *ReactionService) GetCounts(id int, object string) (int, int, error) {
	reactions, err := r.storage.GetReactionsById(id, object)
	if err != nil {
		return 0, 0, err
	}
	likes := 0
	dislikes := 0
	for i := 0; i < len(reactions); i++ {
		if reactions[i].Reaction == "like" {
			likes++
		} else {
			dislikes++
		}
	}
	return likes, dislikes, nil
}
