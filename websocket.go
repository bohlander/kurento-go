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
	events    map[string]map[float64]eventHandler // eventName -> handlerId -> handler.
}

var connections = make(map[string]*Connection)

func NewConnection(host string) *Connection {
	if connections[host] != nil {
		return connections[host]
	}

	c := new(Connection)
	connections[host] = c

	c.events = make(map[string]map[float64]eventHandler)
	c.clients = make(map[float64]chan Response)
	var err error
	c.ws, err = websocket.Dial(host+"/kurento", "", "http://127.0.0.1")
	if err != nil {
		log.Fatal(err)
	}
	c.host = host
	go c.handleResponse()
	return c
}

func (c *Connection) Create(m IMediaObject, options map[string]interface{}) {
	elem := &MediaObject{}
	elem.setConnection(c)
	elem.Create(m, options)
}

func (c *Connection) handleResponse() {
	for { // run forever
		r := Response{}
		ev := Event{}
		var message string
		websocket.Message.Receive(c.ws, &message)

		log.Printf("RAW %s", message)

		// Decode into both possible types. One should be valid
		json.Unmarshal([]byte(message), &r)
		json.Unmarshal([]byte(message), &ev)

		isResponse := r.Id > 0 && r.Result != nil
		isEvent := ev.Method == "onEvent"

		//websocket.JSON.Receive(c.ws, &r)
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
			data := val["data"].(map[string]interface{})

			if handlers, ok := c.events[t]; ok {
				for _, handler := range handlers {
					handler(data)
				}
			}
		} else {
			log.Println("Unsupported message from KMS: ", message)
		}

	}
}

func (c *Connection) Request(req map[string]interface{}) <-chan Response {
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
	websocket.JSON.Send(c.ws, req)
	return c.clients[c.clientId]
}

func (c *Connection) Subscribe(event string, handler eventHandler) float64 {
	var registered map[float64]eventHandler
	var ok bool

	if registered, ok = c.events[event]; !ok {
		c.events[event] = make(map[float64]eventHandler)
		registered = c.events[event]
	}

	eventId := c.eventId
	c.eventId += 1
	registered[eventId] = handler
	return eventId
}

func (c *Connection) Unsubscribe(event string, eventId float64) {
	var registered map[float64]eventHandler
	var ok bool
	if registered, ok = c.events[event]; !ok {
		return // not found
	}

	delete(registered, eventId)
}
