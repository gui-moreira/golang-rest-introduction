package user

import (
	"reflect"
	"testing"
)

func TestUserRepoSave(t *testing.T) {
	i := InmemUserRepo{Users: make(map[int]User)}

	tests := []struct {
		name string
		u    *User
		err  error
	}{
		{"should save user", &User{Name: "bruno", Age: 24}, nil},
		{"should not save invalid user", &User{}, ErrValidateUser},
		{"should not save nil user", nil, ErrValidateUser},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := i.Save(tt.u); err != tt.err {
				t.Errorf("InmemUserRepo.Save() error = %v, wantErr %v", err, tt.err)
			}
		})
	}

	u, ok := i.Users[1]
	if !ok {
		t.Errorf("The user was not saved")
	}

	if u.Name != "bruno" || u.Age != 24 {
		t.Errorf("The user was saved incorectly")
	}
}

func TestUserRepoGet(t *testing.T) {
	u1 := User{ID: 1, Name: "gui", Age: 27}

	i := InmemUserRepo{
		Users: map[int]User{
			1: u1,
		},
	}

	tests := []struct {
		name string
		id   int
		want *User
		err  error
	}{
		{"should get user", 1, &u1, nil},
		{"should not get user when user is not found", 2, nil, ErrUserNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := i.Get(tt.id)

			if err != tt.err {
				t.Errorf("InmemUserRepo.Get() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InmemUserRepo.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
