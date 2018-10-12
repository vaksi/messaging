/*  service.go
*
* @Author:             Audy Vaksi <vaksipranata@gmail.com>
* @Date:               October 08, 2018
* @Last Modified by:   @vaksi
* @Last Modified time: 08/10/18 13:07
 */

package healths

import (
        "testing"

    . "github.com/smartystreets/goconvey/convey"
)

func TestHealthService_HealthCheckStatus(t *testing.T) {
    svcHealth := &HealthService{}
    Convey("Health Check Status Should Return Health Model", t, func(){
        act := svcHealth.HealthCheckStatus()
        So(act.Service, ShouldEqual, HealthModel{ Service: "user_management"}.Service)
    })
}
