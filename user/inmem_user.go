package user

type InmemUserRepo struct {
	Users map[int]User
}

func (i InmemUserRepo) Get(id int) (*User, error) {
	u, ok := i.Users[id]
	if !ok {
		return nil, ErrUserNotFound
	}
	return &u, nil
}

func (i InmemUserRepo) Save(u *User) error {
	if u == nil || u.Name == "" {
		return ErrValidateUser
	}

	u.ID = len(i.Users) + 1
	i.Users[u.ID] = *u

	return nil
}
