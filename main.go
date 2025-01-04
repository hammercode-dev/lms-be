package main

import (
	"github.com/hammer-code/lms-be/cmd"
)

// @title Golang API (LMS-BE)
// @version 1.0
// @description This is the API documentation for this service.
// @contact.name API Support
// @contact.url http://hammercode.org/support
// @contact.email hammercode.org
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cmd.Execute()
}
