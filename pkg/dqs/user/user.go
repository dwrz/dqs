package user

import (
	"time"

	"github.com/dwrz/dqs/pkg/dqs/diet"
)

type User struct {
	Birthday     time.Time `json:"birthday"`
	Diet         diet.Diet `json:"diet"`
	Gender       string    `json:"gender"`
	Height       int       `json:"height"`
	Name         string    `json:"name"`
	TargetDQS    int       `json:"targetDQS"`
	TargetWeight int       `json:"targetWeight"`
	Weight       int       `json:"weight"`
}

var DefaultUser = User{
	Diet: diet.VEGETARIAN,
}
