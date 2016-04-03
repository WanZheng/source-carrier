package client

import (
	"fmt"
	"log"
	"net/http"
)

func (c *SyncClient) handleSync(w http.ResponseWriter, req *http.Request) {
	if err := c.sync(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "done")
}

func (c *SyncClient) runHttpServer() error {
	http.HandleFunc("/sync", c.handleSync)

	addr := fmt.Sprintf("127.0.0.1:%d", c.port)
	log.Print("listen at: ", addr)
	return http.ListenAndServe(addr, nil)
}
