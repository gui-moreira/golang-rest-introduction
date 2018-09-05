# Passo a passo da apresentação

## Criando camada de persistencia

1. `go get github.com/gui-moreira/golang-rest-introduction`
2. Abrir projeto
3. Criar pasta user e arquivo user.go

```golang
package user

import "errors"

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Repo interface {
	Get(id int) (*User, error)
	Save(u *User) error
}


var (
	ErrUserNotFound = errors.New("User not found")
	ErrValidateUser = errors.New("Invalid user")
)
```


4. Criar arquivo inmem_user.go

```golang
package user

type InmemUserRepo struct {
	Users map[int]User
}

func (i InmemUserRepo) Get(ID int) (*User, error) {
	return nil, nil
}

func (i InmemUserRepo) Save(u *User) error {
	if u == nil || u.Name == "" {
		return ErrValidateUser // errors.New("Invalid user")
	}

	u.ID = len(i.Users) + 1
	i.Users[u.ID] = *u

	return nil
}
```

5. Criar arquivo inmem_user_test.go (mesmo package)

```golang
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
```

6. Criar func Get em inmem_user.go

```golang
func (i InmemUserRepo) Get(id int) (*User, error) {
	u, ok := i.Users[id]
	if !ok {
		return nil, ErrUserNotFound // errors.New("User not found")V
	}
	return &u, nil
}
```

7. Criar test para Get

```golang
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
```

---

## Criando camada REST

1. Criar main.go

```golang
package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gui-moreira/golang-rest-introduction/user"
)

var repo user.Repo

func main() {
	repo = user.InmemUserRepo{Users: make(map[int]user.User)}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.Run(":8080")
}
```

2. Executar curl de ping

``` 
curl localhost:8080/ping
```

3. Criar endpoint POST de user

```golang
r.POST("/users", func(c *gin.Context) {
    u := user.User{}
    err := c.ShouldBindJSON(&u)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
        return
    }

    if err := repo.Save(&u); err != nil {
        if err == user.ErrValidateUser {
            c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{"id": u.ID})
})
```

4. Executar curl de POST

```
curl -X POST localhost:8080/users -d '{"name": "Marquinhos","age": 30}'
```

5. Criar endpoint GET de user

```golang
r.GET("/users/:id", func(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	u, err := repo.Get(id)
	if err != nil {
		if err == user.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, u)
})
```

6. Executar curl de GET

```
curl localhost:8080/users/1
```