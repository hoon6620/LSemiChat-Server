package handler

import (
	"app/api/application/interactor"
	"app/api/infrastructure/lcontext"
	"app/api/llog"
	"app/api/presentation/response"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type SocketHandler interface {
	WebsocketConnect(w http.ResponseWriter, r *http.Request)
}

type socketHandler struct {
	messageInteractor interactor.MessageInteractor
	threadInteractor  interactor.ThreadInteractor
}

func NewSocketHandler(mi interactor.MessageInteractor, ti interactor.ThreadInteractor) SocketHandler {
	return &socketHandler{
		messageInteractor: mi,
		threadInteractor:  ti,
	}
}

type SocketData struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type SocketMessageResponse struct {
	AuthorID  string     `json:"author"`
	ThreadID  string     `json:"thread"`
	Grade     int        `json:"grade"`
	Message   string     `json:"message"`
	CreatedAt *time.Time `json:"created_at"`
}

type SocketMessageRequest struct {
	Message string `json:"message"`
	Grade   int    `json:"grade"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var connList = make(map[string]*websocket.Conn)

func (sh *socketHandler) WebsocketConnect(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	connect, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		llog.Error(err)
	}
	userID, err := lcontext.GetUserIDFromContext(r.Context())
	if err != nil {
		response.Unauthorized(w, errors.Wrap(err, "failed to authentication"), "failed to authentication. please login")
		return
	}

	thread, err := sh.CheckReadOnly(connect, userID)
	if err != nil {
		response.Unauthorized(w, errors.Wrap(err, "failed to authentication"), "failed to authentication. please login")
		return
	}

	connList[userID] = connect
	if thread != "" {
		SendNotices(userID, "Web Socket Connected (room: "+thread+")")
		sh.webSocketProccessing(connect, userID, thread)
	} else {
		SendNotices(userID, "Web Socket Connected")
	}
}

func (sh *socketHandler) webSocketProccessing(connect *websocket.Conn, userID string, thread string) {
	var sd SocketData
	var err error
	for {
		//read message from websocket
		mType, p, err := connect.ReadMessage()
		if err != nil {
			llog.Error(err)
			break
		}
		//connection closed
		if mType == -1 {
			err = nil
			break
		}
		//unmarshal message to socketdata
		if err = json.Unmarshal(p, &sd); err != nil {
			llog.Error(err)
			break
		}
		//check datatype
		if sd.Type == "Message" {
			var message SocketMessageRequest
			if err = json.Unmarshal([]byte(sd.Data), &message); err != nil {
				llog.Error(err)
				break
			}
			if err = sh.sendMessage(userID, thread, message); err != nil {
				llog.Error(err)
			}
		}
	}
	removeConnect(userID, err)
}

func SendNotices(userID string, notice string) error {
	sd := &SocketData{
		Type: "notice",
		Data: notice,
	}
	js, err := json.Marshal(sd)
	if err != nil {
		llog.Warn(err)
	}
	connList[userID].WriteMessage(1, js)
	return nil
}

func (sh *socketHandler) sendMessage(authorID string, threadID string, msg SocketMessageRequest) error {
	message, err := sh.messageInteractor.Create(msg.Message, msg.Grade, authorID, threadID)
	if err != nil {
		return err
	}

	smrs := &SocketMessageResponse{
		AuthorID:  authorID,
		ThreadID:  threadID,
		Message:   message.Message,
		Grade:     message.Grade,
		CreatedAt: message.CreatedAt,
	}

	str, err := json.Marshal(smrs)
	if err != nil {
		return err
	}

	sd := &SocketData{
		Type: "message",
		Data: string(str),
	}

	sdJson, err := json.Marshal(sd)
	if err != nil {
		return err
	}

	members, err := sh.threadInteractor.GetMembersByThreadID(threadID)
	if err != nil {
		return err
	}

	for _, user := range members {
		if connList[user.ID] != nil {
			connList[user.ID].WriteJSON(sdJson)
		}
	}

	return nil
}

func removeConnect(id string, err error) {
	if err != nil {
		connList[id].WriteMessage(websocket.CloseMessage, []byte(err.Error()))
	} else {
		connList[id].WriteMessage(websocket.CloseMessage, []byte("websocket connetion closed"))
	}

	connList[id].Close()
	delete(connList, id)
	log.Println(id + ":Client disconnected")
}

func (sh *socketHandler) CheckReadOnly(connect *websocket.Conn, userID string) (string, error) {
	_, data, err := connect.ReadMessage()
	if err != nil {
		return "", err
	}
	var sd SocketData
	if err = json.Unmarshal(data, &sd); err != nil {
		log.Print(err)
		return "", err
	}

	if sd.Type != "setup" {
		return "", errors.New("Socket data type is not 'setup'")
	}

	if sd.Data == "" {
		return "", nil
	}

	return sd.Data, nil //for test

	// //check authorization
	// if !sh.threadInteractor.IsParticipated(sd.Data, userID) {
	// 	return "", errors.New(userID + " are not participated in room " + sd.Data)
	// } else {
	// 	return sd.Data, nil
	// }
}
