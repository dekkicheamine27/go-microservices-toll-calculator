package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go/truck-toll-calculator/types"
)




type Client struct {
	EndPoint string
}

func NewClient(endPoint string) *Client {
    return &Client{EndPoint: endPoint}
}

func (c *Client) AggregateDistance(distance types.Distance) error{
   b , err := json.Marshal(distance)
   if err!= nil {
       return err
   }

   req, err := http.NewRequest("POST", c.EndPoint, bytes.NewReader(b))

   if err!= nil {
       return err
   }

   resp, err := http.DefaultClient.Do(req)
   if err!= nil {
       return err
   }

   if resp.StatusCode != http.StatusOK {
	  fmt.Errorf("the serve responsed with no 200 status code %d ", resp.StatusCode)
   }

   return nil

}

