/*  model.go
*
* @Author:             Audy Vaksi <vaksipranata@gmail.com>
* @Date:               October 09, 2018
* @Last Modified by:   @vaksi
* @Last Modified time: 09/10/18 16:01 
 */

package users

// UserModel of user
type UserModel struct {
    ID       uint   `json:"id"`
    Email    string `json:"email"`
    FullName string `json:"full_name"`
    Gender   int8   `json:"gender"`
}
