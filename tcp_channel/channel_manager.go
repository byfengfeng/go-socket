package tcp_channel

import (
	_interface "game_frame/interface"
)

type channelManager struct {
	channelMap map[int64]*channel
}

func NewChannelManager() _interface.IChannelManager {
	return &channelManager{
		channelMap: map[int64]*channel{},
	}
}
