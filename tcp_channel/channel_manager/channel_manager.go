package channel_manager

import (
	//_interface "game_frame/interface"
	"game_frame/tcp_channel/channel"
	"sync"
)

type channelManager struct {
	channelMap map[int64]*channel.Channel
	channelMtx sync.RWMutex
}

func NewChannelManager() *channelManager {
	return &channelManager{
		channelMap: make(map[int64]*channel.Channel),
	}
}

//添加连接
func (c *channelManager) AddChannel(channel *channel.Channel,uid int64)  {
	c.channelMtx.RLock()
	c.channelMap[uid] = channel
	c.channelMtx.RUnlock()

}

//删除连接
func (c *channelManager) DelChannel(channel *channel.Channel,uid int64) error {
	c.channelMtx.RLock()
	channel,ok := c.channelMap[uid]
	if ok {
		err := channel.Conn.Close()
		if err != nil {
			return err
		}
		delete(c.channelMap,uid)
	}
	c.channelMtx.RUnlock()
	return nil
}

func (c *channelManager) GetchannelManager() map[int64]*channel.Channel {
	return c.channelMap
}