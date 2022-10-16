//Author: Jhon Riascos

package service

import (
	"fmt"

	//"log"
	"errors"
	"net/http"
   //
	//"time"
	"bufio"
	"os"
)

type Option func(*Server) error
type Server struct {
	//listener http.Handler
}
// =========== Files --handle files




// =========== Config --improve with a json with the info
const server_addr = ":3000"
const client_addr, port string = "http://localhost:", "3000"

//Check Server funcs to improve 
/*
func WithTimeout(timeout time.Duration) Option {
    return func(server *Server) error {

        // timeout checks & assignment here
        return &server
    }
}


func service_id(id string) (){
    return id
}
*/


//======= Channels
//
//Example: [[channel1,user, user2] ,[channel2,user,user3 ], [channel3, user2,user3, user] ]

/*var num_channels int =3

s := make([][]string,num_channels)

func Create_channel(name string) {
	
	
}
*/

func Connect(addr, port string) {
	//+"channel/"
	url := addr + port
	res, err := http.Get(url)

	if err != nil {
		fmt.Printf("> Error: %s\n", err)
		fmt.Println("> Client is closed")
		os.Exit(0)
	}

	scanner := bufio.NewScanner(res.Body)
	scanner.Scan()
	fmt.Println(scanner.Text())

}

func start_server() {
	Answers()

	err := http.ListenAndServe(server_addr, nil)
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("> Error: %s\n", err)
		}
	}
}

func Cli(view string) {
	
	status := true

	if view == "server" {

		fmt.Println("> Init server")
		fmt.Print("> Listening and answering\n> ")
		start_server()
		for status {

			var command string

			fmt.Scan(&command)
			if command == "close" {
				status = false

				os.Exit(1)
				fmt.Println("> Server is closed")
			}

			// funcs with flags

		}
	}
	if view == "client" {
		fmt.Println("> Init Client")
		Connect(client_addr, port)

		for status {
			var command string
			fmt.Print("> ")
			fmt.Scan(&command)
			if command == "close" {
				status = false
				//os.Kill()
				fmt.Println("> Client is closed")
			}
			// create funcs with flags

			if command == "subscribe" {
				Connect(client_addr, port+"/channel/")
			}
			if command == "unsubscribe" {
				Connect(client_addr, port+"/channel/")
			}

		}
	}

	os.Exit(0)
}

func Action(do_this ...string) http.Handler {

	//verificar error

	if len(do_this) <= 0 {

		fn_default := func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "> Initialize some channels...")

		}
		return http.HandlerFunc(fn_default)

	}

	if do_this[0] == "sub" {
		fn := func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "> Subscribing to a channel...")

		}
		return http.HandlerFunc(fn)

	}
	if do_this[0] == "unsub" {
		fn := func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "> Unsubscribing to a channel...")
		}
		return http.HandlerFunc(fn)

	}
	
	fn_default := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Trying to Subscribing to a channel... :D Error")

	}
	return http.HandlerFunc(fn_default)

}

func Answers() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "> Connected...")
	})

	http.HandleFunc("/channel/", Action().ServeHTTP)

	

}

func New_server() {
	// handle , options(func) ...Option
	Cli("server")
}

//Client funcs

func New_client() {

	Cli("client")

}
