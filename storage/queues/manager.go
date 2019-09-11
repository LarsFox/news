package queues

import (
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats"

	pb "github.com/LarsFox/newsgrpc"

	"github.com/LarsFox/news/storage/dbs"
)

// Manager ...
type Manager struct {
	getNewsSubj    string
	natsConnection *nats.Conn
	dbm            *dbs.Manager
}

// NewManager — конструктор менеджера работы с очередями.
func NewManager(addr, getNewsSubj string, dbm *dbs.Manager) (*Manager, error) {
	// Подключение к серверу nats.io.
	natsConnection, err := nats.Connect(addr)
	if err != nil {
		return nil, err
	}

	return &Manager{
		natsConnection: natsConnection,
		getNewsSubj:    getNewsSubj,
		dbm:            dbm,
	}, nil
}

func (m *Manager) Listen() error {
	ch := make(chan *nats.Msg, 64)
	if _, err := m.natsConnection.ChanSubscribe(m.getNewsSubj, ch); err != nil {
		return err
	}

	log.Println("Subscribed to nats")
	for {
		m.ReplyNewsPiece(<-ch)
	}
}

func (m *Manager) ReplyNewsPiece(msg *nats.Msg) {
	req := &pb.GetNewsPieceRequest{}
	if err := proto.Unmarshal(msg.Data, req); err != nil {
		log.Println(err)
		return
	}

	var resp *pb.GetNewsPieceReply
	news, err := m.dbm.GetNewsPiece(req.NewsID)
	switch err {
	case nil:
		resp = &pb.GetNewsPieceReply{Header: news.Header, Date: news.Date}
	case dbs.ErrNotFound:
		resp = &pb.GetNewsPieceReply{ErrorCode: 1001}
	case dbs.ErrInternal:
		resp = &pb.GetNewsPieceReply{ErrorCode: 1}
	default:
		log.Println(err)
		resp = &pb.GetNewsPieceReply{ErrorCode: 1}
	}

	b, err := proto.Marshal(resp)
	if err != nil {
		return
	}

	if err := m.natsConnection.Publish(msg.Reply, b); err != nil {
		log.Println(err)
	}
}
