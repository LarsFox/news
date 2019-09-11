package queues

import (
	"errors"
	"log"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats"

	pb "github.com/LarsFox/newsgrpc"
)

var (
	// ErrInternal — внутренняя ошибка.
	ErrInternal = errors.New("internal error")
)

// Manager ...
type Manager struct {
	getNewsSubj    string
	natsConnection *nats.Conn
}

// NewManager — конструктор менеджера работы с очередями.
func NewManager(addr, getNewsSubj string) (*Manager, error) {
	// Подключение к серверу nats.io.
	natsConnection, err := nats.Connect(addr)
	if err != nil {
		return nil, err
	}

	return &Manager{
		getNewsSubj:    getNewsSubj,
		natsConnection: natsConnection,
	}, nil
}

func (m *Manager) RequestNewsPiece(newsID string) (*pb.GetNewsPieceReply, error) {
	b, err := proto.Marshal(&pb.GetNewsPieceRequest{NewsID: newsID})
	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}

	msg, err := m.natsConnection.Request(m.getNewsSubj, b, time.Second*5)
	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}

	reply := &pb.GetNewsPieceReply{}
	if err := proto.Unmarshal(msg.Data, reply); err != nil {
		log.Println(err)
		return nil, ErrInternal
	}
	return reply, nil
}
