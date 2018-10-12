/*  service.go
*
* @Author:             Audy Vaksi <vaksipranata@gmail.com>
* @Date:               October 09, 2018
* @Last Modified by:   @vaksi
* @Last Modified time: 09/10/18 16:01
 */

package service

import (
    "fmt"
    "testing"

    . "github.com/smartystreets/goconvey/convey"
    "github.com/stretchr/testify/mock"
    "github.com/vaksi/user_management/internal/users"
    "github.com/vaksi/user_management/internal/users/repository/mocks"
)

func TestUserService_CreateUser(t *testing.T) {
    Convey("CreateUser Should Be Success", t, func() {
        // Parameter
        email := "testing@gmail.com"
        fullName := "Test 1"
        gender := int8(1)

        // Mocking Repo User
        repoUser := new(mocks.UserRepositoryFactory)
        repoUser.On("Store", mock.MatchedBy(func(u *users.UserModel) bool {
            return u.Email == email && u.FullName == u.FullName && u.Gender == u.Gender
        })).Return(nil)

        // Get Service User
        svcUser := UserService{RepoUser: repoUser}
        exp := svcUser.CreateUser(email, fullName, gender)

        // Assertion
        So(exp, ShouldBeNil)
    })

    Convey("CreateUser Should Be Error", t, func() {
        // Parameter
        email := "testing@gmail.com"
        fullName := "Test 1"
        gender := int8(1)

        // Mocking Repo User
        repoUser := new(mocks.UserRepositoryFactory)
        repoUser.On("Store", &users.UserModel{Email: email, FullName: fullName, Gender: gender}).Return(fmt.Errorf("something error"))

        // Get Service User
        svcUser := UserService{RepoUser: repoUser}
        exp := svcUser.CreateUser(email, fullName, gender)

        // Assertion
        So(exp, ShouldBeError)
    })
}
