package informer

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Phoenix1504e/mini-control-plane/pkg/api"
)

func Watch(url string, handler func(api.WatchEvent)) {
	for {
		resp, err := http.Get(url)
		if err != nil {
			log.Println("watch connect failed:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		decoder := json.NewDecoder(resp.Body)

		for {
			var event api.WatchEvent
			err := decoder.Decode(&event)
			if err != nil {
				if err == io.EOF {
					log.Println("watch closed, reconnecting")
				} else {
					log.Println("watch decode error:", err)
				}
				resp.Body.Close()
				break
			}

			handler(event)
		}

		time.Sleep(1 * time.Second)
	}
}
