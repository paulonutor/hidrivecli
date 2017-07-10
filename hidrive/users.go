package hidrive

import (
	"github.com/dghubble/sling"
)

type User struct {
	AccountId   string `json:account`
	Alias       string `json:alias`
	Description string `json:descr`
	Email       string `json:email`
	IsOwner     bool   `json:is_owner`
	IsAdmin     bool   `json:is_admin`
	Home        string `json:home`
}

type UserListParams struct {
	AccountId string   `url:"account,omitempty"`
	Scope     string   `url:"scope,omitempty"`
	Fields    []string `url:"fields,comma,omitempty"`
}

type UserService struct {
	client *sling.Sling
}

func newUserService(baseClient *sling.Sling) *UserService {
	return &UserService{
		client: baseClient.Path("user/"),
	}
}

func (s *UserService) List(params *UserListParams) ([]User, error) {
	users := new([]User)
	_, err := s.client.New().Get("").QueryStruct(params).ReceiveSuccess(users)

	return *users, err
}
