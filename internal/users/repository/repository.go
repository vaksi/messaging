/*  userRepository.go.go
*
* @Author:             Audy Vaksi <vaksipranata@gmail.com>
* @Date:               October 10, 2018
* @Last Modified by:   @vaksi
* @Last Modified time: 10/10/18 11:18 
 */

package repository

import (
    "github.com/vaksi/user_management/internal/users"
)

// UserRepository of user
type UserRepository struct {

}

// UserRepositoryFactory of User
type UserRepositoryFactory interface {
    Store(*users.UserModel) error
}
