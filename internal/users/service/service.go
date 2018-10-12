/*  service.go
*
* @Author:             Audy Vaksi <vaksipranata@gmail.com>
* @Date:               October 09, 2018
* @Last Modified by:   @vaksi
* @Last Modified time: 09/10/18 16:01 
 */

package service

import (
    "github.com/sirupsen/logrus"
    "github.com/vaksi/user_management/internal/users"
    "github.com/vaksi/user_management/internal/users/repository"
)

// UserService of struct service user
type UserService struct {
    RepoUser repository.UserRepositoryFactory
}

// UserServiceFactory of user serviceßß
type UserServiceFactory interface {
    CreateUser(email, fullName string, gender int8) (error)
}

// CreateUser of function for create user
func (u *UserService) CreateUser(email, fullName string, gender int8) (err error) {
    // store data user
    err = u.RepoUser.Store(&users.UserModel{
        Email:    email,
        FullName: fullName,
        Gender:   gender,
    })
    if err != nil {
        logrus.Error(err)
        return
    }
    return
}
