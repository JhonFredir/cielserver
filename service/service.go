//Author: Jhon Riascos

package service

import (
	"fmt"

	//"log"
	"errors"
	//"flag"
	"net/http"
	//"time"
	"bufio"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"io"
)

type Option func(*Server) error
type Server struct {
	//listener http.Handler
}

// =========== Config --improve with a json with the info
const server_addr = ":3000"
const client_addr, port string = "http://localhost:", "3000"

// --improve with a Db connection
var update_channels = true
var update_channels_ptr = &update_channels

func check(e error) {
	if e != nil {
		//panic(e)
	fmt.Println(e)
	}
}

// Files and folders
func Create_folder(name string) {
	//err :=
	os.Mkdir(name, 0766)
	//check(err)

}

func Delete_folder(name string) {

	os.Remove(name)
	// missing case when the folder have files and it don't erase,  but in the end all is erase,
	// but create problems if create the same channel
}

func Read_message(message *http.Response) {
	scanner := bufio.NewScanner(message.Body)

	for scanner.Scan() { // how to read n lines of the response

		fmt.Println(scanner.Text())

	}
	defer message.Body.Close()
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func Connect(addr, port string) {
	//+"channel/"
	url := addr + port
	res, err := http.Get(url)

	if err != nil {
		fmt.Printf("> Error: %s\n", err)
		fmt.Println("> Client is closed")
		os.Exit(0)
	}

	Read_message(res)

}
func Download_file_req(message *http.Request,folder , name string ) {
	
	data,err := io.ReadAll(message.Body)
    check(err)

	file, ferr := os.Create("Channels/"+folder+"/"+name) //like open a file
    check(ferr)

	defer file.Close()

	w_err := os.WriteFile("Channels/"+folder+"/"+name , data , 0644)
    check(w_err)
	
	
}
func Download_file_resp(message *http.Response,folder , name string ) {
	
	data,err := io.ReadAll(message.Body)
    check(err)

	file, ferr := os.Create(name) //like open a file
    check(ferr)

	defer file.Close()

	w_err := os.WriteFile(name , data , 0644)
    check(w_err)
	
	
}

func Send_file(recipient ...string) {
	//recipient[0]  mode "to_server" 
	//recipient[1]  url to send to server

	if recipient[0] == "to_server" {
		fmt.Println("> Input the file's name(it should be where you run the client) ")
		fmt.Println("> Upload: ")

		fmt.Print("   ")
		var file_name string
		fmt.Scan(&file_name)
		data, err := os.Open(file_name)
		check(err)

		content, err := data.Stat()
		check(err)
		fmt.Println("   content size: ", content.Size(), "bytes")

		//req, err:= http.NewRequest("POST",recipient[1],data)
		resp, err := http.Post(recipient[1], "all", data)
		//recipient[1]  url to send to server


		check(err)
		Read_message(resp)

	}

	//recipient[0]  mode "to_client" 
	//recipient[1]  client

	if recipient[0] == "to_client" {
		fmt.Println("> Sending files to: "+ recipient[1] )
   // print list of sub channel

	}

}

// Channels
//Example: [[channel1,user, user2] ,[user,channel3,channel2 ], [channel3, user2,user3, user] ]

var num_channels int = 12

var lst_channels = make([][]string, num_channels)

func Be_there(list []string, mode string, value string) bool {
	if mode == "client" {
		for i := 0; i < len(list); i++ {
			if list[i] == value {
				return true
			}
		}

	}

	if mode == "channel" {
		for i := 0; i < len(lst_channels); i++ {
			if len(lst_channels[i]) == 0 {
				continue
			}
			if lst_channels[i][0] == value {
				return true
			}
		}

	}
	return false

}

func Create_channel(name string) {
	status := true
	for status {
		if *update_channels_ptr {

			*update_channels_ptr = false // red light

			//lst_channels[i] is not used in the channel mode

			channel_exists := Be_there(lst_channels[0], "channel", name) //lst_channels[0] is not used in the channel mode
			if channel_exists {
				//lst_channels[i] is not used in the channel mode
				// if the channel exists
				// how to improve  to use without a parameter in mode channel ?
				fmt.Println("the channel already exists: " + name)
				*update_channels_ptr = true // green light
				status = false

				break
			}
			//assigment
			for i := 0; i < num_channels; i++ {

				if len(lst_channels[i]) != 0 {
					continue
				}

				if len(lst_channels[i]) == 0 {

					lst_channels[i] = append(lst_channels[i], name)
					Create_folder("Channels/" + name)
					break
				}

			}
			//missing the case where there are no more channels
			//i can use copy and create another array +1

			*update_channels_ptr = true // green light
			status = false
		}
	}

}
func Supr_channel(channel, user string) {
	status := true
	for status {
		if *update_channels_ptr {

			*update_channels_ptr = false // red light

			for i := 0; i < num_channels; i++ {
				if len(lst_channels[i]) > 0 {

					if lst_channels[i][0] == channel {

						lst_channels[i] = nil
						Delete_folder("Channels/" + channel)

						break

					}
				}
			}

			*update_channels_ptr = true // green light
			status = false

		}
	}
}

func Sub_channel(channel, user string) {
	status := true
	for status {
		if *update_channels_ptr {

			*update_channels_ptr = false // red light

			for i := 0; i < num_channels; i++ {
				if len(lst_channels[i]) > 0 {

					if lst_channels[i][0] == channel {
						if !Be_there(lst_channels[i], "client", user) { //== false
							lst_channels[i] = append(lst_channels[i], user)
							break
						}
						//missing the case where the channel don't exist

					}
				}
			}
			*update_channels_ptr = true // green light
			status = false
		}
	}
}

func Unsub_channel(channel, user string) {
	status := true
	for status {
		if *update_channels_ptr {

			*update_channels_ptr = false // red light

			for i := 0; i < num_channels; i++ {
				if len(lst_channels[i]) > 0 {

					if lst_channels[i][0] == channel {

						for clien := 1; clien < len(lst_channels[i]); clien++ {

							if len(lst_channels[i]) > 1 {
								if lst_channels[i][clien] == user {
									next := clien + 1

									lst_channels[i] = append(lst_channels[i][:clien], lst_channels[i][next:]...)

									break
								}

							}

						}

						break
					}
				}
			}
			*update_channels_ptr = true // green light
			status = false
		}
	}
}

//========= Server and client

func start_server() {
	Answers()

	err := http.ListenAndServe(server_addr, nil)
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("> Error: %s\n", err)
		}
	}
}

// finish func
func Stadistics() {
	fmt.Println("> Stadistics")

}

func Action(do_this ...string) http.Handler {

	//split commands

	fn_default := func(w http.ResponseWriter, r *http.Request) {

		instruction := strings.Split(r.URL.Path[4:], "/")

		//instruction[0] == user
		//instruction[1] == action
		//instruction[2] == channel
		//instruction[3] == file's name

		channel_exists := Be_there(lst_channels[0], "channel", instruction[2]) 
			//lst_channels[0] is not used in the channel mode

		switch instruction[1] {

		case "create":
			fmt.Println("> client:", instruction[0], ",create channel:", instruction[2])
			if channel_exists{
				fmt.Fprint(w, "> the channel does exist: "+instruction[2]+"\n")
				break
				}

			//create channel
			Create_channel(instruction[2])
			fmt.Println(lst_channels)
			fmt.Fprint(w, "> create channel:"+instruction[2]+"\n")

		case "supr":
			fmt.Println("> client:", instruction[0], ",erasing channel:", instruction[2])
			//supr channel
			if !channel_exists{
				fmt.Fprint(w, "> the channel doesn't exist: "+instruction[2]+"\n")
				break
				}
			Supr_channel(instruction[2], instruction[0])
			fmt.Println(lst_channels)
			fmt.Fprint(w, "> delete channel: "+instruction[2]+"\n")

		case "send":
			fmt.Println("> client:", instruction[0], ",sending file to channel:", instruction[2])
			
			fmt.Println("cuerpo")

			if !channel_exists{
			fmt.Fprint(w, "> the channel doesn't exist: "+instruction[2]+"\n")
			break
			}
			fmt.Println( instruction[2],instruction[3])

			Download_file_req( r ,instruction[2],instruction[3])
			// read and write the file on the channel

			fmt.Println("cuerpo")

			fmt.Fprint(w, "> sending a file to channel: "+instruction[2]+"\n")

		case "receive":
			fmt.Println("> client:", instruction[0], ",downloading files... ")
			
			Send_file("to_client",instruction[0] )
			//receive file
			//esta sub
			//lista de channel
			// listar archivos y luego enviar uno a uno
			//download complete

		case "sub":
			fmt.Println("> client:", instruction[0], ",subscribing to", instruction[2])
			//subscribe
			if !channel_exists{
				fmt.Fprint(w, "> the channel doesn't exist: "+instruction[2]+"\n")
				break
				}
			Sub_channel(instruction[2], instruction[0])
			fmt.Println(lst_channels)

		case "unsub":
			fmt.Println("> client:", instruction[0], ",unsubscribing to ", instruction[2])
			//unsubscribe
			if !channel_exists{
				fmt.Fprint(w, "> the channel doesn't exist: "+instruction[2]+"\n")
				break
				}
			Unsub_channel(instruction[2], instruction[0])
			fmt.Println(lst_channels)

		default:
			fmt.Fprint(w, "> Doing something(the above)...\n")
		}

	}
	return http.HandlerFunc(fn_default)

}

func Answers() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "> Connected...\n")

	})

	http.HandleFunc("/do/", Action().ServeHTTP)

}

func Cli(view string) {

	//Channel

	if view == "server" {
		Create_folder("Channels")
		//status_server := true

		sigChan := make(chan os.Signal, 1)

		fmt.Println("> Init server")
		go start_server()
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		fmt.Print("> Listening and answering\n")

		<-sigChan
		// ********* how  i can  do a interactive cli after the process is start ?
		//for status_server {

		//	var command string
		//	fmt.Scan(&command)

		//	if command == "close" {
		//	status_server = false

		// missing send the signal to interrupted the process and finished
		//supr channel folder
		//os.Exit(1)

		os.RemoveAll("Channels")
		fmt.Println("> Server is closed")
		os.Exit(0)
	}

	// funcs with flags

	//}

	//}

	if view == "client" {
		status_client := true

		fmt.Println("> Init Client")
		var username string

		fmt.Print("> Input username :\n> ")
		fmt.Scan(&username)

		fmt.Println("> Welcome,", username)
		Connect(client_addr, port)

		for status_client {
			var command string
			var value string

			fmt.Print("> ")

			fmt.Scan(&command, &value)

			if strings.Contains(command, "close") {
				status_client = false
				//Make Des-conexion
				//os.Kill()
				fmt.Println("> Client is closed")
			}
			// create funcs with flags
			if strings.Contains(command, "send") {
				url := client_addr + port + "/do/" + username + "/send/" + value
				Send_file("to_server", url)
				//here complete
				//Connect(client_addr, port+"/do/"+username+"/send/"+value)

				continue
			}
			if strings.Contains(command, "receive") {
				Connect(client_addr, port+"/do/"+username+"/receive/"+value) //here complete

				continue
			}

			if strings.Contains(command, "create") {
				Connect(client_addr, port+"/do/"+username+"/create/"+value) //here complete

				continue
			}
			if strings.Contains(command, "supr") {
				Connect(client_addr, port+"/do/"+username+"/supr/"+value) //here complete

				continue
			}

			if command == "sub" {
				Connect(client_addr, port+"/do/"+username+"/sub/"+value) //here complete

				continue
			}

			if command == "unsub" {
				Connect(client_addr, port+"/do/"+username+"/unsub/"+value) //here complete
				continue

			}

		}
		os.Exit(0)
	}

}

func New_server() {
	// handle , options(func) ...Option
	Cli("server")
}

//Client funcs

func New_client() {
	ch := make(chan bool)
	go Cli("client")

	ch <- true

}
