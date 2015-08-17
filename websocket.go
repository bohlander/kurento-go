package kurento

import (
	"encoding/json"
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)

// Error that can be filled in response
type Error struct {
	Code    int64
	Message string
	Data    string
}

const ConnectionLost = -1

// Implements error built-in interface
func (e *Error) Error() string {
	return fmt.Sprintf("[%d] %s %s", e.Code, e.Message, e.Data)
}

// Response represents server response
type Response struct {
	Jsonrpc string
	Id      float64
	Result  map[string]string // should change if result has no several form
	Error   *Error
}

type Event struct {
	Jsonrpc string
	Method  string
	Params  map[string]interface{}
	Error   *Error
}

type Connection struct {
	clientId  float64
	eventId   float64
	clients   map[float64]chan Response
	host      string
	ws        *websocket.Conn
	SessionId string
	events    map[string]map[string]map[string]eventHandler // eventName -> objectId -> handlerId -> handler.
	Dead      chan bool
	IsDead    bool
}

var connections = make(map[string]*Connection)

func NewConnection(host string) *Connection {
	if connections[host] != nil {
		return connections[host]
	}

	c := new(Connection)
	connections[host] = c

	c.events = make(map[string]map[string]map[string]eventHandler)
	c.clients = make(map[float64]chan Response)
	c.Dead = make(chan bool, 1)

	var err error
	c.ws, err = websocket.Dial(host+"/kurento", "", "http://127.0.0.1")
	if err != nil {
		log.Fatal(err)
	}
	c.host = host
	go c.handleResponse()
	return c
}

func (c *Connection) Create(m IMediaObject, options map[string]interface{}) error {
	elem := &MediaObject{}
	elem.setConnection(c)
	return elem.Create(m, options)
}

func (c *Connection) handleResponse() {
	for { // run forever
		r := Response{}
		ev := Event{}
		var message string
		err := websocket.Message.Receive(c.ws, &message)
		if err != nil {
			log.Printf("Error receiving on websocket %s", err)
			c.IsDead = true
			c.Dead <- true
			break
		}

		if debug {
			log.Printf("RAW %s", message)
		}

		// Decode into both possible types. One should be valid
		json.Unmarshal([]byte(message), &r)
		json.Unmarshal([]byte(message), &ev)

		isResponse := r.Id > 0 && r.Result != nil
		isEvent := ev.Method == "onEvent"

		if isResponse {
			if r.Result["sessionId"] != "" {
				if debug {
					log.Println("sessionId returned ", r.Result["sessionId"])
				}
				c.SessionId = r.Result["sessionId"]
			}
			if debug {
				log.Printf("Response: %v", r)
			}
			// if webscocket client exists, send response to the chanel
			if c.clients[r.Id] != nil {
				c.clients[r.Id] <- r
				// chanel is read, we can delete it
				delete(c.clients, r.Id)
			} else if debug {
				log.Println("Dropped message because there is no client ", r.Id)
				log.Println(r)
			}
		} else if isEvent {

			val := ev.Params["value"].(map[string]interface{})
			if debug {
				log.Printf("Received event value %v", val)
			}

			t := val["type"].(string)
			objectId := val["object"].(string)

			data := val["data"].(map[string]interface{})

			if handlers, ok := c.events[t]; ok {
				if objHandlers, ok := handlers[objectId]; ok {
					for _, handler := range objHandlers {
						handler(data)
					}
				}
			}
		} else if debug {
			log.Println("Unsupported message from KMS: ", message)
		}
	}
}

func (c *Connection) Request(req map[string]interface{}) <-chan Response {
	if c.IsDead {
		errchan := make(chan Response, 1)
		errresp := Response{
			Id: req["id"].(float64),
			Error: &Error{
				Code:    ConnectionLost,
				Message: "No connection to Kurento server",
			},
		}
		errchan <- errresp
		return errchan
	}

	c.clientId++
	req["id"] = c.clientId
	if c.SessionId != "" {
		req["sessionId"] = c.SessionId
	}
	c.clients[c.clientId] = make(chan Response)
	if debug {
		j, _ := json.MarshalIndent(req, "", "    ")
		log.Println("json", string(j))
	}
	err := websocket.JSON.Send(c.ws, req)
	if err != nil {
		log.Printf("Error sending on websocket %s", err)
		c.Dead <- true
		c.IsDead = true

		delete(c.clients, c.clientId)

		errchan := make(chan Response, 1)
		errresp := Response{
			Id: req["id"].(float64),
			Error: &Error{
				Code:    ConnectionLost,
				Message: "No connection to Kurento server",
			},
		}

		errchan <- errresp
		return errchan
	}
	return c.clients[c.clientId]
}

func (c *Connection) Subscribe(event, objectId, handlerId string, handler eventHandler) {
	var oh map[string]map[string]eventHandler
	var ok bool

	if oh, ok = c.events[event]; !ok {
		c.events[event] = make(map[string]map[string]eventHandler)
		oh = c.events[event]
	}

	var he map[string]eventHandler
	if he, ok = oh[objectId]; !ok {
		oh[objectId] = make(map[string]eventHandler)
		he = oh[objectId]
	}

	he[handlerId] = handler
}

func (c *Connection) Unsubscribe(event, objectId, handlerId string) {
	var oh map[string]map[string]eventHandler
	var he map[string]eventHandler
	var ok bool
	if oh, ok = c.events[event]; !ok {
		return // not found
	}

	if he, ok = oh[objectId]; !ok {
		return // not found
	}

	delete(he, handlerId)
}
