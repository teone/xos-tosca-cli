// TODO

// [ ] execute a post with file content
// [ ] run cli in a container

package main

import (
	"github.com/abiosoft/ishell"
	"net/http"
	"log"
	"io/ioutil"
	"fmt"
	"os"
	"encoding/json"
	"strings"
)


func printBanner(c *ishell.Shell) {
	var banner = `
$$\   $$\  $$$$$$\   $$$$$$\      $$$$$$$$\  $$$$$$\   $$$$$$\   $$$$$$\   $$$$$$\           $$$$$$\  $$\       $$$$$$\
$$ |  $$ |$$  __$$\ $$  __$$\     \__$$  __|$$  __$$\ $$  __$$\ $$  __$$\ $$  __$$\         $$  __$$\ $$ |      \_$$  _|
\$$\ $$  |$$ /  $$ |$$ /  \__|       $$ |   $$ /  $$ |$$ /  \__|$$ /  \__|$$ /  $$ |        $$ /  \__|$$ |        $$ |
 \$$$$  / $$ |  $$ |\$$$$$$\ $$$$$$\ $$ |   $$ |  $$ |\$$$$$$\  $$ |      $$$$$$$$ |$$$$$$\ $$ |      $$ |        $$ |
 $$  $$<  $$ |  $$ | \____$$\\______|$$ |   $$ |  $$ | \____$$\ $$ |      $$  __$$ |\______|$$ |      $$ |        $$ |
$$  /\$$\ $$ |  $$ |$$\   $$ |       $$ |   $$ |  $$ |$$\   $$ |$$ |  $$\ $$ |  $$ |        $$ |  $$\ $$ |        $$ |
$$ /  $$ | $$$$$$  |\$$$$$$  |       $$ |    $$$$$$  |\$$$$$$  |\$$$$$$  |$$ |  $$ |        \$$$$$$  |$$$$$$$$\ $$$$$$\
\__|  \__| \______/  \______/        \__|    \______/  \______/  \______/ \__|  \__|         \______/ \________|\______|
`
	c.Print(banner)
}


// -------------------------------- //
//            Configuration         //
// -------------------------------- //

//var url = "http://10.90.2.10/xos-tosca"
var url = "http://127.0.0.1:9102"

var recipe_folder = "/Users/teone/Sites/cord_profile/"

var xos_username = "xosadmin@opencord.org"

var xos_password = "ORR13pBZ8yrAZ42QiKhc"

// -------------------------------- //
//          End Configuration       //
// -------------------------------- //


func readFile(file string) string {
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	return string(dat)
}

func main(){
	// create new shell.
	// by default, new shell includes 'exit', 'help' and 'clear' commands.
	shell := ishell.New()

	printBanner(shell)
	shell.AddCmd(&ishell.Cmd{
		Name: "list-tosca",
		Help: "List the available TOSCA definition in XOS",
		Func: func(c *ishell.Context) {

			response, err := http.Get(url)

			if err != nil {
				fmt.Print(err.Error())
				os.Exit(1)
			}

			responseData, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}

			type Response map[string]interface{}

			var responseObject Response
			json.Unmarshal(responseData, &responseObject)

			recipeList := []string{}

			for k := range responseObject {
				recipeList = append(recipeList, k)
			}

			fmt.Println(strings.Join(recipeList,", "))
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "get-tosca",
		Help: "Get a TOSCA definition from XOS",
		Func: func(c *ishell.Context) {
			// get username
			c.Print("model: ")
			recipe := c.ReadLine()

			url := url + "/custom_type/" + recipe
			c.Println(url)

			response, err := http.Get(url)

			if err != nil {
				fmt.Print(err.Error())
				os.Exit(1)
			}

			responseData, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(responseData))

		},
	})

	// /Users/teone/Sites/cord_profile/public-net.yaml
	shell.AddCmd(&ishell.Cmd{
		Name: "run-tosca",
		Help: "Run a tosca recipe",
		Func: func(c *ishell.Context) {

			// read all the files in the TOSCA recipe folder
			recipeList := []string{}
			files, err := ioutil.ReadDir(recipe_folder)
			if err != nil {
				log.Fatal(err)
			}

			for _, f := range files {
				if strings.Contains(f.Name(), "yaml") {
					recipeList = append(recipeList, f.Name())
				}
			}

			choice := c.MultiChoice(recipeList, "Select a recipe to send:")

			recipe := readFile(recipe_folder + recipeList[choice])

			url := url + "/run"

			req, err := http.NewRequest("POST", url, strings.NewReader(recipe))
			req.Header.Set("xos-username", xos_username)
			req.Header.Set("xos-password", xos_password)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			fmt.Println("response Status:", resp.Status)
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println("response Body:", string(body))

		},
	})

	// run shell
	shell.Run()
}