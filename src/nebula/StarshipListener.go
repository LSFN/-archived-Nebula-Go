package nebula

import (
	"net"

	"code.google.com/p/goprotobuf/proto"
	"code.google.com/p/uuid"

	"lsfn/common"
)

type StarshipListener struct {
	conn     *net.TCPConn
	Messages chan STSup
}

func (listener *StarshipListener) Listen() {
	var lengthToRead uint64
	var lengthVariant := new(common.Variant)
	var currentReadBuffer []byte
	for {
		var data []byte
		bytes, err := listener.conn.Read(data)
		if err != nil {
			listener.Messages.Close()
			break
		}
		for bytes > 0 {
			if lengthToRead == 0 {
				for _, b := range data {
					lengthVariant.ConnectByte(b)
					if lengthVariant.IsComplete() {
						lengthToRead = lengthVariant.ToUint64()
						lengthVariant.Reset()
						break
					}
				}
			}
			if lengthToRead > 0 {
				if bytes >= lengthToRead {
					currentReadBuffer = append(currentReadBuffer, data[:lengthToRead]...)
					data = data[lengthToRead:]
					bytes -= lengthToRead
					lengthToRead = 0
					var upMessage common.STSup = new(common.STSup)
					proto.Unmarshal(currentReadBuffer, upMessage)
					listener.Messages <- upMessage
				} else {
					currentReadBuffer = append(currentReadBuffer, data...)
					data = nil
					bytes -= lengthToRead
					lengthToRead -= bytes
					bytes = 0
				}
			}
		}
	}
}

func (listener *StarshipListener) readSingleMessage() (common.STSup, error) {
	// Read off the length of the message into a variant
	lengthVariant := new(common.Variant)
	singleByte := make([]byte, 1)
	for !lengthVariant.isComplete() {
		bytes, err := listener.conn.Read(singleByte)
		if err != nil {
			return nil, err
		}
		if bytes == 1 {
			lengthVariant.ConnectByte(singleByte[0])
		}
	}

	// Receive a message of the stated length
	receiverSlice := make([]byte, lengthVariant.ToUint64)
	var bytes int
	for bytes < len(receiverSlice) {
		x, err := listener.conn.Read(receiverSlice[bytes:])
		if err != nil {
			return nil, err
		}
		bytes += x
	}

	// Unmarshal the message into a protobuf structure
	result := new(common.STSup)
	err := proto.Unmarshal(receiverSlice, result)
	if err != nil {
		result = nil
	}

	return result, err
}

func (listener *StarshipListener) SendMessage(downMessage common.STSdown) {
	err := listener.sendSingleMessage(downMessage)
	if err != nil {
		listener.Disconnect()
	}
}

func (listener *StarshipListener) sendSingleMessage(downMessage common.STSdown) error {
	rawMessage, err := proto.Marshal(downMessage)
	if err != nil {
		return err
	}

	var bytes int
	variantLength := new(common.Variant)
	variantLength.FromUint64(len(rawMessage))
	rawLength := variantLength.ToBytes()
	for bytes < len(rawLength) {
		x, err := listener.conn.Write(rawLength[bytes:])
		if err != nil {
			return err
		}
		bytes += x
	}

	bytes = 0
	for bytes < len(rawMessage) {
		x, err := listener.conn.Write(rawMessage[bytes:])
		if err != nil {
			return err
		}
		bytes += x
	}

	return nil
}

func (listener *StarshipListener) Disconnect() {
	listener.conn.Close()
	listener.Messages.Close()
}

// TODO possibly disconnect
func (listener *StarshipListener) Handshake(gameID string, allowJoin bool, rejoinIDs []string) bool {
	// Send the JoinInfo message
	joinInfoMessage := &common.STSdown{
		JoinInfo := &common.STSdown_JoinInfo{
			AllowJoin: allowJoin
			GameIDToken: gameID
		}
	}
	err := listener.sendSingleMessage(joinInfoMessage)
	// if joins are not allowed then when allowJoin is false, handshakes will end here
	if err != nil {
		return false
	}

	// Receive the JoinRequest message
	joinRequestMessage, err := listener.receiveSingleMessage()
	if err != nil {
		return false
	}
	if !joinRequestMessage.HasJoinRequest() {
		return false
	}

	// Send the JoinResponse message
	joinType := joinRequestMessage.GetJoinRequest.GetJoinType()
	if(joinType == common.STSup_JoinRequest_JOIN) {
		if allowJoin {}
			id string := uuid.New()
			listener.sendSingleMessage(&common.STSdown{
				JoinResponse: &common.STSdown_JoinResponse{
					Type: common.STSdown_JoinResponse_JOIN_ACCEPTED
					RejoinToken: id
				}
			})	
			return true		
		} else {}
			listener.sendSingleMessage(&common.STSdown{
				JoinResponse: &common.STSdown_JoinResponse{
					Type: common.STSdown_JoinResponse_JOIN_REJECTED
				}
			})
			return false
		}
	} else {
		successfulRejoin bool = false;
		if joinRequestMessage.HasRejoinID() {
			rejoinID = joinRequestMessage.GetRejoinID()
			for _, orphanID := range rejoinIDs {
				if rejoinID == orphanID {
					successfulRejoin = true
				}
			}
		}
		if successfulRejoin {
			listener.sendSingleMessage(&common.STSdown{
				JoinResponse: &common.STSdown_JoinResponse{
					Type: common.STSdown_JoinResponse_REJOIN_ACCEPTED
					RejoinToken: joinRequestMessage.GetRejoinID()
				}
			})
			return true
		} else {
			listener.sendSingleMessage(&common.STSdown{
				JoinResponse: &common.STSdown_JoinResponse{
					Type: common.STSdown_JoinResponse_JOIN_REJECTED
				}
			})
			return false
		}
	}
}