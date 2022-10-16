//Author: Jhon Riascos
package main

import (
	
	"cielserver.com/service"
)

func init_server() {

	//Create a func with more options
	service.New_server() 
	
	/* ,
	/		service.service_id("Server_session"),
			service.WithTimeout(time.Minute),
	*/

}

func main() {

	init_server()

	

}
