package client

// type Loader inteface {
// 	func Load[T any](host string) (T, error)
// }

// import (
// 	"fmt"
// 	"io"
// 	"net/http"
// )

// func main() {
// 	client := &http.Client{}
// 	urls := []string{
// 		"http://localhost:8111/state",
// 		"http://localhost:8111/mission.json",
// 		"http://localhost:8111/map_obj.json",
// 		"http://localhost:8111/map_info.json",
// 		"http://localhost:8111/gamechat?lastId=0",
// 		"http://localhost:8111/hudmsg?lastEvt=0&lastDmg=0",
// 		"http://localhost:8111/indicators",
// 	}

// 	for _, url := range urls {
// 		resp, err := client.Get(url)
// 		if err != nil {
// 			fmt.Printf("Error fetching %s: %v\n", url, err)
// 			continue
// 		}
// 		defer resp.Body.Close()
// 		body, err := io.ReadAll(resp.Body)
// 		if err != nil {
// 			fmt.Printf("Error reading response body for %s: %v\n", url, err)
// 			continue
// 		}
// 		fmt.Println(string(body))
// 	}
// }
