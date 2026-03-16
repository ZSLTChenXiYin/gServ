package tcpserv

import (
	"bytes"
	"encoding/binary"
	"gServ/core/gameserv"
	"gServ/core/log"
	"gServ/core/repository"
	"gServ/pkg/jwt"
	"io"
	"net"

	"github.com/ZSLTChenXiYin/custproto"
)

func handleProtocol(conn net.Conn, v any) error {
	encoder := custproto.NewEncoder(bytes.NewBuffer(nil), binary.BigEndian)
	err := encoder.Encode(v)
	if err != nil {
		return err
	}
	_, err = conn.Write(encoder.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	// 处理TCP连接，确认玩家登录
	login_bytes_slice := make([]byte, 2)
	_, err := io.ReadFull(conn, login_bytes_slice)
	if err != nil {
		log.StdErrorf("TCP服务处理玩家%s首次登录请求读取失败: %v", conn.RemoteAddr().String(), err)
		return
	}
	login_header, err := VerifyProtocolHeader(login_bytes_slice)
	if err != nil {
		log.StdErrorf("TCP服务处理玩家%s首次登录请求验证请求头失败: %v", conn.RemoteAddr().String(), err)
		err = handleProtocol(conn, NewGSERVProtocolHeaderError(err.Error()))
		if err != nil {
			log.StdErrorf("TCP服务处理玩家%s首次登录请求验证请求头响应失败: %v", conn.RemoteAddr().String(), err)
			return
		}
		return
	}
	if login_header.ProtocolType != GSERV_PROTOCOL_TYPE_AUTH_PLAYER_LOGIN {
		log.StdErrorf("TCP服务处理玩家%s首次登录请求验证协议类型失败: 登录协议错误，当前登录协议类型%d", conn.RemoteAddr().String(), login_header.ProtocolType)
		err = handleProtocol(conn, NewGSERVProtocolHeaderError("协议类型错误"))
		if err != nil {
			log.StdErrorf("TCP服务处理玩家%s首次登录请求响应验证协议类型失败: %v", conn.RemoteAddr().String(), err)
		}
		return
	}
	// 读取玩家Token信息
	conn_stream_decoder := custproto.NewStreamDecoder(conn, binary.BigEndian)
	login_req := &GSERVProtocolAuthPlayerLoginRequest{}
	err = conn_stream_decoder.Decode(login_req)
	if err != nil {
		log.StdErrorf("TCP服务处理玩家%s首次登录请求结构化失败: %v", conn.RemoteAddr().String(), err)
		err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
			GSERV_PROTOCOL_AUTH_PLAYER_LOGIN_RESPONSE_FAILURE,
			err.Error(),
		))
		if err != nil {
			log.StdErrorf("TCP服务处理玩家%s首次登录请求结构化响应失败: %v", conn.RemoteAddr().String(), err)
		}
		return
	}
	login_auth_player, err := jwt.ParseAuthPlayerToken(string(login_req.Token))
	if err != nil {
		log.StdErrorf("TCP服务处理玩家%s首次登录请求解析Token失败: %v", conn.RemoteAddr().String(), err)
		err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
			GSERV_PROTOCOL_AUTH_PLAYER_LOGIN_RESPONSE_FAILURE,
			err.Error(),
		))
		if err != nil {
			log.StdErrorf("TCP服务处理玩家%s首次登录请求解析Token响应失败: %v", conn.RemoteAddr().String(), err)
		}
		return
	}
	// 验证玩家Token信息
	login_player_expired_at, err := repository.ScanPlayerExpiredAt(login_auth_player.ID)
	if err != nil {
		log.StdErrorf("TCP服务处理玩家%s首次登录请求验证Token失败: %v", conn.RemoteAddr().String(), err)
		err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
			GSERV_PROTOCOL_AUTH_PLAYER_LOGIN_RESPONSE_FAILURE,
			err.Error(),
		))
		if err != nil {
			log.StdErrorf("TCP服务处理玩家%s首次登录请求解析Token响应失败: %v", conn.RemoteAddr().String(), err)
		}
		return
	}
	// 如果玩家Token中的过期时间小于当前数据库中的过期时间，则表示该Token已失效
	if login_auth_player.ExpiredAt.Before(*login_player_expired_at) {
		err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
			GSERV_PROTOCOL_AUTH_PLAYER_LOGIN_RESPONSE_FAILURE,
			"Token已失效",
		))
		if err != nil {
			log.StdErrorf("TCP服务处理玩家%s首次登录请求验证过期时间响应失败: %v", conn.RemoteAddr().String(), err)
		}
		return
	}

	// 玩家登录后在服务器上线
	err = gameserv.PlayerOnline(uint(login_req.GameID), login_auth_player.ID)
	if err != nil {
		log.StdErrorf("TCP服务处理玩家%s首次登录请求玩家上线失败: %v", conn.RemoteAddr().String(), err)
		err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
			GSERV_PROTOCOL_AUTH_PLAYER_LOGIN_RESPONSE_FAILURE,
			err.Error(),
		))
		if err != nil {
			log.StdErrorf("TCP服务处理玩家%s首次登录请求玩家上线响应失败: %v", conn.RemoteAddr().String(), err)
		}
		return
	}

	err = player_manager.Create(uint(login_req.GameID), login_auth_player.ID, conn)
	if err != nil {
		log.StdErrorf("TCP服务处理玩家%s首次登录请求创建玩家连接失败: %v", conn.RemoteAddr().String(), err)
		err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
			GSERV_PROTOCOL_AUTH_PLAYER_LOGIN_RESPONSE_FAILURE,
			err.Error(),
		))
		if err != nil {
			log.StdErrorf("TCP服务处理玩家%s首次登录请求创建玩家连接响应失败: %v", conn.RemoteAddr().String(), err)
		}
		return
	}
	defer player_manager.Delete(uint(login_req.GameID), login_auth_player.ID)

	err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
		GSERV_PROTOCOL_AUTH_PLAYER_LOGIN_RESPONSE_SUCCESS,
		"玩家登录成功",
	))
	if err != nil {
		log.StdErrorf("TCP服务处理玩家%s首次登录请求响应失败: %v", conn.RemoteAddr().String(), err)
		return
	}

	log.StdInfof("TCP服务处理玩家%s首次登录请求成功", conn.RemoteAddr().String())

	for {
		bytes_slice := make([]byte, 2)
		_, err := io.ReadFull(conn, bytes_slice)
		if err != nil {
			log.StdErrorf("TCP服务处理玩家%s请求失败: %v", conn.RemoteAddr().String(), err)
			return
		}
		header, err := VerifyProtocolHeader(bytes_slice)
		if err != nil {
			log.StdErrorf("TCP服务处理玩家%s请求验证请求头失败: %v", conn.RemoteAddr().String(), err)
			err = handleProtocol(conn, NewGSERVProtocolHeaderError(err.Error()))
			if err != nil {
				log.StdErrorf("TCP服务处理玩家%s请求验证请求头响应失败: %v", conn.RemoteAddr().String(), err)
			}
			return
		}

		switch header.ProtocolType {
		case GSERV_PROTOCOL_TYPE_AUTH_PLAYER_LOGIN:
			// 重复登录，清理登录请求后续的协议数据
			req := &GSERVProtocolAuthPlayerLoginRequest{}
			err = conn_stream_decoder.Decode(req)
			if err != nil {
				log.StdErrorf("TCP服务处理玩家%s登录请求结构化失败: %v", conn.RemoteAddr().String(), err)
				err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
					GSERV_PROTOCOL_AUTH_PLAYER_LOGIN_RESPONSE_FAILURE,
					err.Error(),
				))
				if err != nil {
					log.StdErrorf("TCP服务处理玩家%s登录请求结构化响应失败: %v", conn.RemoteAddr().String(), err)
				}
				return
			}
			err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
				GSERV_PROTOCOL_AUTH_PLAYER_LOGIN_RESPONSE_FAILURE,
				"重复登录协议错误",
			))
			if err != nil {
				log.StdErrorf("TCP服务处理玩家%s登录请求响应失败: %v", conn.RemoteAddr().String(), err)
			}

			log.StdInfof("TCP服务处理玩家%s重复登录请求成功", conn.RemoteAddr().String())

			continue
		case GSERV_PROTOCOL_TYPE_AUTH_PLAYER_LOGOUT:
			// 玩家登出
			req := &GSERVProtocolAuthPlayerLogoutRequest{}
			err = conn_stream_decoder.Decode(req)
			if err != nil {
				log.StdErrorf("TCP服务处理玩家%s登出请求结构化失败: %v", conn.RemoteAddr().String(), err)
				err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
					GSERV_PROTOCOL_AUTH_PLAYER_LOGOUT_RESPONSE_FAILURE,
					err.Error(),
				))
				if err != nil {
					log.StdErrorf("TCP服务处理玩家%s登出请求结构化响应失败: %v", conn.RemoteAddr().String(), err)
				}
				return
			}

			// 读取玩家Token信息
			auth_player, err := jwt.ParseAuthPlayerToken(string(req.Token))
			if err != nil {
				log.StdErrorf("TCP服务处理玩家%s登出请求解析Token失败: %v", conn.RemoteAddr().String(), err)
				err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
					GSERV_PROTOCOL_AUTH_PLAYER_LOGOUT_RESPONSE_FAILURE,
					err.Error(),
				))
				if err != nil {
					log.StdErrorf("TCP服务处理玩家%s登出请求解析Token响应失败: %v", conn.RemoteAddr().String(), err)
				}
				return
			}

			// 玩家登出后在服务器下线
			gameserv.PlayerOffline(uint(req.GameID), auth_player.ID)

			// 玩家登出返回成功
			err = handleProtocol(conn, NewGSERVProtocolAuthPlayerLogoutResponse(
				GSERV_PROTOCOL_AUTH_PLAYER_LOGOUT_RESPONSE_SUCCESS,
				GSERV_PROTOCOL_MESSAGE_TYPE_TEXT,
				"玩家登出成功",
			))
			if err != nil {
				log.StdErrorf("TCP服务处理玩家%s登出请求响应失败: %v", conn.RemoteAddr().String(), err)
			}

			log.StdInfof("TCP服务处理玩家%s登出请求成功", conn.RemoteAddr().String())

			return
		case GSERV_PROTOCOL_TYPE_JOIN_ROOM:
			// 玩家加入房间请求，客户端调用创建房间接口后，再调用该接口加入房间
			req := &GSERVProtocolJoinRoomRequest{}
			err = conn_stream_decoder.Decode(req)
			if err != nil {
				err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
					GSERV_PROTOCOL_JOIN_ROOM_RESPONSE_FAILURE,
					err.Error(),
				))
				if err != nil {
				}
				return
			}

			// 读取玩家Token信息
			auth_player, err := jwt.ParseAuthPlayerToken(string(req.Token))
			if err != nil {
				err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
					GSERV_PROTOCOL_AUTH_PLAYER_LOGOUT_RESPONSE_FAILURE,
					err.Error(),
				))
				if err != nil {
					return
				}
				continue
			}

			// 读取玩家数据
			player := gameserv.GetOnlinePlayer(uint(req.GameID), auth_player.ID)
			if player == nil {
				err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
					GSERV_PROTOCOL_JOIN_ROOM_RESPONSE_FAILURE,
					"玩家未登录",
				))
				if err != nil {
					return
				}
				continue
			}

			// 查询玩家是否在房间中，在房间中就调用离开房间
			if player.GetCurrentRoomID() != 0 {
				err = gameserv.LeaveRoom(uint(req.GameID), player.GetCurrentRoomID(), auth_player.ID)
				if err != nil {
					err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
						GSERV_PROTOCOL_JOIN_ROOM_RESPONSE_FAILURE,
						err.Error(),
					))
					if err != nil {
						return
					}
					continue
				}
			}

			err = gameserv.JoinRoom(uint(req.GameID), req.RoomID, auth_player.ID)
			if err != nil {
				err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
					GSERV_PROTOCOL_JOIN_ROOM_RESPONSE_FAILURE,
					err.Error(),
				))
				if err != nil {
					return
				}
				continue
			}

			err = handleProtocol(conn, NewGSERVProtocolJoinRoomResponse(
				GSERV_PROTOCOL_JOIN_ROOM_RESPONSE_SUCCESS,
				GSERV_PROTOCOL_MESSAGE_TYPE_TEXT,
				"玩家加入房间成功",
			))
			if err != nil {
			}

			continue
		case GSERV_PROTOCOL_TYPE_ROOM_BROADCAST_DATA:
			// 向房间所有玩家（需要包括自己）广播数据请求
			req := &GSERVProtocolRoomBroadcastDataRequest{}
			err = conn_stream_decoder.Decode(req)
			if err != nil {
				err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
					GSERV_PROTOCOL_ROOM_BROADCAST_DATA_RESPONSE_FAILURE,
					err.Error(),
				))
				if err != nil {
				}
				return
			}

			// 验证token
			auth_player, err := jwt.ParseAuthPlayerToken(string(req.Token))
			if err != nil {
				err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
					GSERV_PROTOCOL_ROOM_BROADCAST_DATA_RESPONSE_FAILURE,
					"Token验证失败",
				))
				if err != nil {
					return
				}
				continue
			}

			// 获取玩家信息
			player := gameserv.GetOnlinePlayer(uint(req.GameID), auth_player.ID)
			if player == nil {
				err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
					GSERV_PROTOCOL_ROOM_BROADCAST_DATA_RESPONSE_FAILURE,
					"玩家未登录",
				))
				if err != nil {
					return
				}
				return
			}

			// 检查玩家是否在房间中
			if player.GetCurrentRoomID() != req.RoomID {
				err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
					GSERV_PROTOCOL_ROOM_BROADCAST_DATA_RESPONSE_FAILURE,
					"玩家不在房间内",
				))
				if err != nil {
					return
				}
				continue
			}

			// 广播
			err = broadcast(uint(req.GameID), req.RoomID, auth_player.ID, req.Data)
			if err != nil {
				err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
					GSERV_PROTOCOL_ROOM_BROADCAST_DATA_RESPONSE_FAILURE,
					err.Error(),
				))
				if err != nil {
					return
				}
				continue
			}

			continue
		case GSERV_PROTOCOL_TYPE_ROOM_UNICAST_DATA:
			// 向房间指定玩家（可以指定自己）广播数据请求
			req := &GSERVProtocolRoomUnicastDataRequest{}
			err = conn_stream_decoder.Decode(req)
			if err != nil {
				err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
					GSERV_PROTOCOL_ROOM_BROADCAST_DATA_RESPONSE_FAILURE,
					err.Error(),
				))
				if err != nil {
				}
				return
			}

			// 验证token
			auth_player, err := jwt.ParseAuthPlayerToken(string(req.Token))
			if err != nil {
				err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
					GSERV_PROTOCOL_ROOM_BROADCAST_DATA_RESPONSE_FAILURE,
					"Token验证失败",
				))
				if err != nil {
					return
				}
				continue
			}

			// 获取发送者信息
			sender := gameserv.GetOnlinePlayer(uint(req.GameID), auth_player.ID)
			if sender == nil {
				err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
					GSERV_PROTOCOL_ROOM_BROADCAST_DATA_RESPONSE_FAILURE,
					"玩家未登录",
				))
				if err != nil {
					return
				}
				return
			}

			// 检查发送者是否在房间中
			if sender.GetCurrentRoomID() != req.RoomID {
				err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
					GSERV_PROTOCOL_ROOM_BROADCAST_DATA_RESPONSE_FAILURE,
					"玩家不在房间内",
				))
				if err != nil {
					return
				}
				continue
			}

			// 检查目标玩家是否在房间中
			target_player := gameserv.GetOnlinePlayer(uint(req.GameID), uint(req.PlayerID))
			if target_player == nil || target_player.GetCurrentRoomID() != req.RoomID {
				err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
					GSERV_PROTOCOL_ROOM_BROADCAST_DATA_RESPONSE_FAILURE,
					"目标玩家不在房间内",
				))
				if err != nil {
					return
				}
				continue
			}

			// 单播
			err = unicast(uint(req.GameID), req.RoomID, auth_player.ID, uint(req.PlayerID), req.Data)
			if err != nil {
				err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
					GSERV_PROTOCOL_ROOM_BROADCAST_DATA_RESPONSE_FAILURE,
					err.Error(),
				))
				if err != nil {
					return
				}
				continue
			}

			continue
		case GSERV_PROTOCOL_TYPE_ROOM_MULTICAST_DATA:
			// 向房间多个指定玩家广播数据请求
			req := &GSERVProtocolRoomMulticastDataRequest{}
			err = conn_stream_decoder.Decode(req)
			if err != nil {
				err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
					GSERV_PROTOCOL_ROOM_BROADCAST_DATA_RESPONSE_FAILURE,
					err.Error(),
				))
				if err != nil {
				}
				return
			}

			// 验证token
			auth_player, err := jwt.ParseAuthPlayerToken(string(req.Token))
			if err != nil {
				err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
					GSERV_PROTOCOL_ROOM_BROADCAST_DATA_RESPONSE_FAILURE,
					"Token验证失败",
				))
				if err != nil {
					return
				}
				continue
			}

			// 获取发送者信息
			sender := gameserv.GetOnlinePlayer(uint(req.GameID), auth_player.ID)
			if sender == nil {
				err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
					GSERV_PROTOCOL_ROOM_BROADCAST_DATA_RESPONSE_FAILURE,
					"玩家未登录",
				))
				if err != nil {
					return
				}
				return
			}

			// 检查发送者是否在房间中
			if sender.GetCurrentRoomID() != req.RoomID {
				err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
					GSERV_PROTOCOL_ROOM_BROADCAST_DATA_RESPONSE_FAILURE,
					"玩家不在房间内",
				))
				if err != nil {
					return
				}
				continue
			}

			// 组播
			target_player_ids := make([]uint, req.PlayerIDCount)
			for index := 0; index < int(req.PlayerIDCount); index++ {
				target_player_ids[index] = uint(req.PlayerIDs[index])
			}
			err = multicast(uint(req.GameID), req.RoomID, auth_player.ID, target_player_ids, req.Data)
			if err != nil {
				err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
					GSERV_PROTOCOL_ROOM_BROADCAST_DATA_RESPONSE_FAILURE,
					err.Error(),
				))
				if err != nil {
					return
				}
				continue
			}

			continue
			/*
				case GSERV_PROTOCOL_TYPE_DATA_STORE_REQUEST:
					// 数据存储请求
					req := &GSERVProtocolDataStoreRequest{}
					err = conn_stream_decoder.Decode(req)
					if err != nil {
						err = handleProtocol(conn, NewGSERVProtocolDataStoreResponse(
							GSERV_PROTOCOL_DATA_STORE_RESPONSE_FAILURE,
							GSERV_PROTOCOL_MESSAGE_TYPE_ERROR,
							"协议解析失败: "+err.Error(),
						))
						if err != nil {
							return
						}
						continue
					}

					// 验证token
					auth_player, err := jwt.ParseAuthPlayerToken(string(req.Token))
					if err != nil {
						err = handleProtocol(conn, NewGSERVProtocolDataStoreResponse(
							GSERV_PROTOCOL_DATA_STORE_RESPONSE_FAILURE,
							GSERV_PROTOCOL_MESSAGE_TYPE_ERROR,
							"Token验证失败",
						))
						if err != nil {
							return
						}
						continue
					}

					// 获取数据存储服务
					storage := gameserv.GetDataStorage(uint(req.GameID), auth_player.ID)
					if storage == nil {
						err = handleProtocol(conn, NewGSERVProtocolDataStoreResponse(
							GSERV_PROTOCOL_DATA_STORE_RESPONSE_FAILURE,
							GSERV_PROTOCOL_MESSAGE_TYPE_ERROR,
							"数据存储服务初始化失败",
						))
						if err != nil {
							return
						}
						continue
					}

					// 保存数据
					key := string(req.Key)
					var data any
					err = json.Unmarshal(req.Data, &data)
					if err != nil {
						err = handleProtocol(conn, NewGSERVProtocolDataStoreResponse(
							GSERV_PROTOCOL_DATA_STORE_RESPONSE_FAILURE,
							GSERV_PROTOCOL_MESSAGE_TYPE_ERROR,
							"数据解析失败: "+err.Error(),
						))
						if err != nil {
							return
						}
						continue
					}

					err = storage.SaveData(key, data)
					if err != nil {
						err = handleProtocol(conn, NewGSERVProtocolDataStoreResponse(
							GSERV_PROTOCOL_DATA_STORE_RESPONSE_FAILURE,
							GSERV_PROTOCOL_MESSAGE_TYPE_ERROR,
							"数据保存失败: "+err.Error(),
						))
						if err != nil {
							return
						}
						continue
					}

					// 发送成功响应
					err = handleProtocol(conn, NewGSERVProtocolDataStoreResponse(
						GSERV_PROTOCOL_DATA_STORE_RESPONSE_SUCCESS,
						GSERV_PROTOCOL_MESSAGE_TYPE_TEXT,
						"数据保存成功",
					))
					if err != nil {
						return
					}

					continue
				case GSERV_PROTOCOL_TYPE_DATA_LOAD_REQUEST:
					// 数据加载请求
					req := &GSERVProtocolDataLoadRequest{}
					err = conn_stream_decoder.Decode(req)
					if err != nil {
						err = handleProtocol(conn, NewGSERVProtocolDataLoadResponse(
							GSERV_PROTOCOL_DATA_LOAD_RESPONSE_FAILURE,
							[]byte(""),
						))
						if err != nil {
							return
						}
						continue
					}

					// 验证token
					auth_player, err := jwt.ParseAuthPlayerToken(string(req.Token))
					if err != nil {
						err = handleProtocol(conn, NewGSERVProtocolDataLoadResponse(
							GSERV_PROTOCOL_DATA_LOAD_RESPONSE_FAILURE,
							[]byte(""),
						))
						if err != nil {
							return
						}
						continue
					}

					// 获取数据存储服务
					storage := gameserv.GetDataStorage(uint(req.GameID), auth_player.ID)
					if storage == nil {
						err = handleProtocol(conn, NewGSERVProtocolDataLoadResponse(
							GSERV_PROTOCOL_DATA_LOAD_RESPONSE_FAILURE,
							[]byte(""),
						))
						if err != nil {
							return
						}
						continue
					}

					// 加载数据
					key := string(req.Key)
					data, err := storage.LoadData(key)
					if err != nil {
						err = handleProtocol(conn, NewGSERVProtocolDataLoadResponse(
							GSERV_PROTOCOL_DATA_LOAD_RESPONSE_FAILURE,
							[]byte(""),
						))
						if err != nil {
							return
						}
						continue
					}

					// 序列化数据
					jsonData, err := json.Marshal(data)
					if err != nil {
						err = handleProtocol(conn, NewGSERVProtocolDataLoadResponse(
							GSERV_PROTOCOL_DATA_LOAD_RESPONSE_FAILURE,
							[]byte(""),
						))
						if err != nil {
							return
						}
						continue
					}

					// 发送成功响应
					err = handleProtocol(conn, NewGSERVProtocolDataLoadResponse(
						GSERV_PROTOCOL_DATA_LOAD_RESPONSE_SUCCESS,
						jsonData,
					))
					if err != nil {
						return
					}

					continue
				case GSERV_PROTOCOL_TYPE_DATA_DELETE_REQUEST:
					// 数据删除请求
					req := &GSERVProtocolDataDeleteRequest{}
					err = conn_stream_decoder.Decode(req)
					if err != nil {
						err = handleProtocol(conn, NewGSERVProtocolDataDeleteResponse(
							GSERV_PROTOCOL_DATA_DELETE_RESPONSE_FAILURE,
							GSERV_PROTOCOL_MESSAGE_TYPE_ERROR,
							"协议解析失败: "+err.Error(),
						))
						if err != nil {
							return
						}
						continue
					}

					// 验证token
					auth_player, err := jwt.ParseAuthPlayerToken(string(req.Token))
					if err != nil {
						err = handleProtocol(conn, NewGSERVProtocolDataDeleteResponse(
							GSERV_PROTOCOL_DATA_DELETE_RESPONSE_FAILURE,
							GSERV_PROTOCOL_MESSAGE_TYPE_ERROR,
							"Token验证失败",
						))
						if err != nil {
							return
						}
						continue
					}

					// 获取数据存储服务
					storage := gameserv.GetDataStorage(uint(req.GameID), auth_player.ID)
					if storage == nil {
						err = handleProtocol(conn, NewGSERVProtocolDataDeleteResponse(
							GSERV_PROTOCOL_DATA_DELETE_RESPONSE_FAILURE,
							GSERV_PROTOCOL_MESSAGE_TYPE_ERROR,
							"数据存储服务初始化失败",
						))
						if err != nil {
							return
						}
						continue
					}

					// 删除数据
					key := string(req.Key)
					err = storage.DeleteData(key)
					if err != nil {
						err = handleProtocol(conn, NewGSERVProtocolDataDeleteResponse(
							GSERV_PROTOCOL_DATA_DELETE_RESPONSE_FAILURE,
							GSERV_PROTOCOL_MESSAGE_TYPE_ERROR,
							"数据删除失败: "+err.Error(),
						))
						if err != nil {
							return
						}
						continue
					}

					// 发送成功响应
					err = handleProtocol(conn, NewGSERVProtocolDataDeleteResponse(
						GSERV_PROTOCOL_DATA_DELETE_RESPONSE_SUCCESS,
						GSERV_PROTOCOL_MESSAGE_TYPE_TEXT,
						"数据删除成功",
					))
					if err != nil {
						return
					}

					continue
			*/
		default:
			err = handleProtocol(conn, NewGSERVProtocolUniversalResponseError(
				GSERV_PROTOCOL_MESSAGE_TYPE_ERROR,
				"未知协议",
			))
			if err != nil {
				return
			}

			return
		}
	}
}
