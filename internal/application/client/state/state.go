package state

import "github.com/bubaew95/yandex-diplom-2/internal/model"

type State struct {
	Token string
	User  model.User
}

func NewState(token string, user model.User) *State {
	return &State{Token: token, User: user}
}

func (state *State) SetToken(token string) {
	state.Token = token
}

func (state *State) GetToken() string {
	return state.Token
}

func (state *State) GetUser() model.User {
	return state.User
}

func (state *State) SetUser(user model.User) {
	state.User = user
}
