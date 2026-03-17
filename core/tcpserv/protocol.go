package tcpserv

import (
	"encoding/binary"
	"errors"

	"github.com/ZSLTChenXiYin/custproto"
)

const (
	GSERV_PROTOCOL_VERSION_UNKNOWN = iota
	GSERV_PROTOCOL_VERSION_FIRST
)

const (
	GSERV_PROTOCOL_TYPE_UNKNOWN = iota
	GSERV_PROTOCOL_TYPE_HEADER_ERROR
	GSERV_PROTOCOL_TYPE_UNIVERSAL_RESPONSE_ERROR
	GSERV_PROTOCOL_TYPE_AUTH_PLAYER_LOGIN
	GSERV_PROTOCOL_TYPE_AUTH_PLAYER_LOGOUT
	GSERV_PROTOCOL_TYPE_JOIN_ROOM
	GSERV_PROTOCOL_TYPE_ROOM_BROADCAST_DATA
	GSERV_PROTOCOL_TYPE_ROOM_UNICAST_DATA
	GSERV_PROTOCOL_TYPE_ROOM_MULTICAST_DATA
)

type GSERVProtocolHeader struct {
	ProtocolVersion uint8
	ProtocolType    uint8
}

type GSERVProtocolHeaderError struct {
	GSERVProtocolHeader
	ErrorLength uint32
	Error       []byte `custproto:"ErrorLength"`
}

func NewGSERVProtocolHeaderError(err string) *GSERVProtocolHeaderError {
	return &GSERVProtocolHeaderError{
		GSERVProtocolHeader: GSERVProtocolHeader{
			ProtocolVersion: GSERV_PROTOCOL_VERSION_FIRST,
			ProtocolType:    GSERV_PROTOCOL_TYPE_HEADER_ERROR,
		},
		ErrorLength: uint32(len(err)),
		Error:       []byte(err),
	}
}

const (
	GSERV_PROTOCOL_MESSAGE_TYPE_UNKNOWN = iota
	GSERV_PROTOCOL_MESSAGE_TYPE_TEXT
	GSERV_PROTOCOL_MESSAGE_TYPE_ERROR
)

type GSERVProtocolMessage struct {
	MessageType   uint8
	MessageLength uint32
	Message       []byte `custproto:"MessageLength"`
}

type GSERVProtocolUniversalResponseError struct {
	GSERVProtocolHeader
	Status uint8
	GSERVProtocolMessage
}

func NewGSERVProtocolUniversalResponseError(status uint8, message string) *GSERVProtocolUniversalResponseError {
	return &GSERVProtocolUniversalResponseError{
		GSERVProtocolHeader: GSERVProtocolHeader{
			ProtocolVersion: GSERV_PROTOCOL_VERSION_FIRST,
			ProtocolType:    GSERV_PROTOCOL_TYPE_UNIVERSAL_RESPONSE_ERROR,
		},
		Status: status,
		GSERVProtocolMessage: GSERVProtocolMessage{
			MessageType:   GSERV_PROTOCOL_MESSAGE_TYPE_ERROR,
			MessageLength: uint32(len(message)),
			Message:       []byte(message),
		},
	}
}

type GSERVProtocolToken struct {
	TokenLength uint16
	Token       []byte `custproto:"TokenLength"`
}

type GSERVProtocolAuthPlayerLoginRequest struct {
	GSERVProtocolToken
	GameID uint32
}

const (
	GSERV_PROTOCOL_AUTH_PLAYER_LOGIN_RESPONSE_FAILURE = iota
	GSERV_PROTOCOL_AUTH_PLAYER_LOGIN_RESPONSE_SUCCESS
	// 登录成功，但是该账号还用于启动其他游戏，即同账户启动了多个游戏
	GSERV_PROTOCOL_AUTH_PLAYER_LOGIN_RESPONSE_ALREADY_LOGIN_OTHER_GAME
	// 登录成功，但是该账号所登录的游戏正于异地在线状态
	GSERV_PROTOCOL_AUTH_PLAYER_LOGIN_RESPONSE_ALREADY_LOGIN_OTHER_PLACE
)

type GSERVProtocolAuthPlayerLoginResponse struct {
	GSERVProtocolHeader
	Status uint8
	GSERVProtocolMessage
}

func NewGSERVProtocolAuthPlayerLoginResponse(status uint8, message_type uint8, message string) *GSERVProtocolAuthPlayerLoginResponse {
	return &GSERVProtocolAuthPlayerLoginResponse{
		GSERVProtocolHeader: GSERVProtocolHeader{
			ProtocolVersion: GSERV_PROTOCOL_VERSION_FIRST,
			ProtocolType:    GSERV_PROTOCOL_TYPE_AUTH_PLAYER_LOGIN,
		},
		Status: status,
		GSERVProtocolMessage: GSERVProtocolMessage{
			MessageType:   message_type,
			MessageLength: uint32(len(message)),
			Message:       []byte(message),
		},
	}
}

type GSERVProtocolAuthPlayerLogoutRequest struct {
	GSERVProtocolToken
	GameID uint32
}

const (
	GSERV_PROTOCOL_AUTH_PLAYER_LOGOUT_RESPONSE_FAILURE = iota
	GSERV_PROTOCOL_AUTH_PLAYER_LOGOUT_RESPONSE_SUCCESS
	GSERV_PROTOCOL_AUTH_PLAYER_LOGOUT_RESPONSE_NOT_LOGIN
)

type GSERVProtocolAuthPlayerLogoutResponse struct {
	GSERVProtocolHeader
	Status uint8
	GSERVProtocolMessage
}

func NewGSERVProtocolAuthPlayerLogoutResponse(status uint8, message_type uint8, message string) *GSERVProtocolAuthPlayerLogoutResponse {
	return &GSERVProtocolAuthPlayerLogoutResponse{
		GSERVProtocolHeader: GSERVProtocolHeader{
			ProtocolVersion: GSERV_PROTOCOL_VERSION_FIRST,
			ProtocolType:    GSERV_PROTOCOL_TYPE_AUTH_PLAYER_LOGOUT,
		},
		Status: status,
		GSERVProtocolMessage: GSERVProtocolMessage{
			MessageType:   message_type,
			MessageLength: uint32(len(message)),
			Message:       []byte(message),
		},
	}
}

type GSERVProtocolJoinRoomRequest struct {
	GSERVProtocolToken
	GameID uint32
	RoomID uint64
}

const (
	GSERV_PROTOCOL_JOIN_ROOM_RESPONSE_FAILURE = iota
	GSERV_PROTOCOL_JOIN_ROOM_RESPONSE_SUCCESS
)

type GSERVProtocolJoinRoomResponse struct {
	GSERVProtocolHeader
	Status uint8
	GSERVProtocolMessage
}

func NewGSERVProtocolJoinRoomResponse(status uint8, message_type uint8, message string) *GSERVProtocolJoinRoomResponse {
	return &GSERVProtocolJoinRoomResponse{
		GSERVProtocolHeader: GSERVProtocolHeader{
			ProtocolVersion: GSERV_PROTOCOL_VERSION_FIRST,
			ProtocolType:    GSERV_PROTOCOL_TYPE_JOIN_ROOM,
		},
		Status: status,
		GSERVProtocolMessage: GSERVProtocolMessage{
			MessageType:   message_type,
			MessageLength: uint32(len(message)),
			Message:       []byte(message),
		},
	}
}

type GSERVProtocolRoomBroadcastDataRequest struct {
	GSERVProtocolToken
	GameID     uint32
	RoomID     uint64
	DataLength uint32
	Data       []byte `custproto:"DataLength"`
}

const (
	GSERV_PROTOCOL_ROOM_BROADCAST_DATA_RESPONSE_FAILURE = iota
)

type GSERVProtocolRoomBroadcastDataResponse struct {
	GSERVProtocolHeader
	GameID         uint32
	RoomID         uint64
	SourcePlayerID uint32
	DataLength     uint32
	Data           []byte `custproto:"DataLength"`
}

func NewGSERVProtocolRoomBroadcastDataResponse(game_id uint32, room_id uint64, source_player_id uint32, data []byte) *GSERVProtocolRoomBroadcastDataResponse {
	return &GSERVProtocolRoomBroadcastDataResponse{
		GSERVProtocolHeader: GSERVProtocolHeader{
			ProtocolVersion: GSERV_PROTOCOL_VERSION_FIRST,
			ProtocolType:    GSERV_PROTOCOL_TYPE_ROOM_BROADCAST_DATA,
		},
		GameID:         game_id,
		RoomID:         room_id,
		SourcePlayerID: source_player_id,
		DataLength:     uint32(len(data)),
		Data:           data,
	}
}

type GSERVProtocolRoomUnicastDataRequest struct {
	GSERVProtocolToken
	GameID     uint32
	RoomID     uint64
	PlayerID   uint32
	DataLength uint32
	Data       []byte `custproto:"DataLength"`
}

type GSERVProtocolRoomUnicastDataResponse struct {
	GSERVProtocolHeader
	GameID         uint32
	RoomID         uint64
	SourcePlayerID uint32
	DataLength     uint32
	Data           []byte `custproto:"DataLength"`
}

type GSERVProtocolRoomMulticastDataRequest struct {
	GSERVProtocolToken
	GameID        uint32
	RoomID        uint64
	PlayerIDCount uint32
	PlayerIDs     []uint32 `custproto:"PlayerIDCount"`
	DataLength    uint32
	Data          []byte `custproto:"DataLength"`
}

type GSERVProtocolRoomMulticastDataResponse struct {
	GSERVProtocolHeader
	GameID         uint32
	RoomID         uint64
	SourcePlayerID uint32
	DataLength     uint32
	Data           []byte `custproto:"DataLength"`
}

func VerifyVersion(version uint8) error {
	switch version {
	case GSERV_PROTOCOL_VERSION_UNKNOWN:
		return errors.New("Unknown protocol version")
	case GSERV_PROTOCOL_VERSION_FIRST:
		return nil
	default:
		return errors.New("Unknown protocol version")
	}
}

func VerifyProtocolType(protocol_type uint8) error {
	switch protocol_type {
	case GSERV_PROTOCOL_TYPE_UNKNOWN:
		return errors.New("Unknown protocol type")
	case GSERV_PROTOCOL_TYPE_HEADER_ERROR:
		return errors.New("Protocol error")
	case GSERV_PROTOCOL_TYPE_AUTH_PLAYER_LOGIN,
		GSERV_PROTOCOL_TYPE_AUTH_PLAYER_LOGOUT,
		GSERV_PROTOCOL_TYPE_JOIN_ROOM,
		GSERV_PROTOCOL_TYPE_ROOM_BROADCAST_DATA,
		GSERV_PROTOCOL_TYPE_ROOM_UNICAST_DATA,
		GSERV_PROTOCOL_TYPE_ROOM_MULTICAST_DATA:
		return nil
	default:
		return errors.New("Unknown protocol type")
	}
}

func VerifyProtocolHeader(data []byte) (GSERVProtocolHeader, error) {
	header := GSERVProtocolHeader{}

	buffer_decoder := custproto.NewBufferDecoder(data, binary.BigEndian)
	err := buffer_decoder.Decode(&header)
	if err != nil {
		return GSERVProtocolHeader{}, err
	}

	err = VerifyVersion(header.ProtocolVersion)
	if err != nil {
		return GSERVProtocolHeader{}, err
	}

	err = VerifyProtocolType(header.ProtocolType)
	if err != nil {
		return GSERVProtocolHeader{}, err
	}

	return header, nil
}
