package main-1

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mitchellh/mapstructure"
)

func main() {

	// response, err := http.Get("http://metevents.herokuapp.com/events")
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	data, _ := ioutil.ReadAll(response.Body)

	// 	ioutil.WriteFile("metevents.txt", data, 0777)

	dat, err := ioutil.ReadFile("/Users/rli233/myApps/rexGo/metevents.txt")
	if err != nil {
		fmt.Println(err)
	} else {
		var myevents map[string]interface{}
		err := json.Unmarshal(dat, &myevents)
		if err != nil {
			fmt.Println(err)
		}

		// eventlist := make([]MetEventList, 0)

		var temp MetEventList

		// fmt.Println(string(*myevents["data"]))

		mapstructure.Decode(myevents, temp)

		fmt.Println("%v", temp)
	}

	// // var f interface{}
	// var myevents MetList
	// err := json.Unmarshal(data, &myevents)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// // mylist := myevents["Data"]
	// var mydata interface{}
	// myevents.(map[string][]MetEvent)
	// describe(mydata["data"])
	// mydata2 := mydata["data"].([]interface{})
	// for _, v := range mydata2 {
	// 	describe(v)
	// 	fmt.Println(v)
	// }
	// for _, v := range mylist {
	// 	fmt.Println(v)
	// }
	// fmt.Println(mydata2[0])
	// describe(mydata2)

	// for k, v := range mydata2 {
	// 	fmt.Println("%s plays on %s")
	// }

	// fmt.Println("%T\n", mydata)
	// fmt.Println(mydata["data"])
	// }

	// event := MyCustomResource{"test", "10112019"}

	// b, _ := json.Marshal(event)

	// fmt.Println(b)

	// var eventReceived MyCustomResource
	// json.Unmarshal(b, &eventReceived)

	// if len(os.Args) < 2 {
	// 	printInput()
	// } else {
	// 	fmt.Println(os.Args[1], os.Args[2])
	// }
	// fmt.Println(eventReceived)

}

func printInput() {
	fmt.Println("begin typing ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if scanner.Err() != nil {
		fmt.Println("error")
	}
}

//MyCustomResource
type MyCustomResource struct {
	Name string
	Date string
}

type MetEvent struct {
	Date string `json:"Date"`
	Name string `json:"Name"`
}

type MetEventList struct {
	Data   string
	Events []MetEvent
}

func (r MyCustomResource) String() string {
	return fmt.Sprintf("%s %s", r.Name, r.Date)
}
func describe(i interface{}) {
	fmt.Printf("(%v, %T)\n", i, i)
}
