/*  model.go
*
* @Author:             Audy Vaksi <vaksipranata@gmail.com>
* @Date:               October 08, 2018
* @Last Modified by:   @vaksi
* @Last Modified time: 08/10/18 13:07 
 */

package healths

// HealthModel for model healthCheck
type HealthModel struct {
    Service string `json:"service"`
}
