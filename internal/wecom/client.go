package wecom

import (
	"context"
	"encoding/json"

	"cnb.cool/zhiqiangwang/pkg/logx"
)

// AIBotClient 企业微信智能机器人客户端
// API文档: https://developer.work.weixin.qq.com/document/path/100719
type AIBotClient struct {
	ctx            context.Context
	Token          string
	EncodingAESKey string
}

// UserReq 用户请求消息
type UserReq struct {
	Msgid    string `json:"msgid"`
	Aibotid  string `json:"aibotid"`
	Chattype string `json:"chattype"`
	From     struct {
		Userid string `json:"userid"`
	} `json:"from"`
	Msgtype string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
	Stream struct {
		Id string `json:"id"`
	} `json:"stream"`
}

// UserResp 回复消息
type UserResp struct {
	Msgtype string `json:"msgtype"`
	Stream  Stream `json:"stream"`
}

// Stream 流式消息
type Stream struct {
	Id      string `json:"id"`
	Finish  bool   `json:"finish"`
	Content string `json:"content"`
	MsgItem []struct {
		Msgtype string `json:"msgtype"`
		Image   struct {
			Base64 string `json:"base64"`
			Md5    string `json:"md5"`
		} `json:"image"`
	} `json:"msg_item"`
}

// NewAIBotClient 创建企业微信AI机器人客户端
func NewAIBotClient(ctx context.Context, token, encodingAESKey string) (*AIBotClient, error) {
	return &AIBotClient{
		ctx:            ctx,
		Token:          token,
		EncodingAESKey: encodingAESKey,
	}, nil
}

// VerifyURL 验证URL(用于配置回调地址)
func (c *AIBotClient) VerifyURL(signature, timestamp, nonce, echoStr string) (string, error) {
	wx, _, err := NewWXBizJsonMsgCrypt(c.Token, c.EncodingAESKey, "")
	if err != nil {
		return "", err
	}

	code, sReplyEchoStr := wx.VerifyURL(signature, timestamp, nonce, echoStr)
	if code != WXBizMsgCrypt_OK {
		logx.Error("VerifyURL failed: code %d", code)
		return "", getErrorMessage(code)
	}

	return sReplyEchoStr, nil
}

// DecryptUserReq 解密用户请求消息
func (c *AIBotClient) DecryptUserReq(signature, timestamp, nonce, msg string) (*UserReq, error) {
	wx, _, err := NewWXBizJsonMsgCrypt(c.Token, c.EncodingAESKey, "")
	if err != nil {
		return nil, err
	}

	code, reqMsg := wx.DecryptMsg(msg, signature, timestamp, nonce)
	if code != WXBizMsgCrypt_OK {
		return nil, getErrorMessage(code)
	}

	var data UserReq
	logx.Debug("Decrypted user request: %s", reqMsg)
	err = json.Unmarshal([]byte(reqMsg), &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// MakeStreamResp 生成流式响应消息
func (c *AIBotClient) MakeStreamResp(nonce, id, content string, isFinish bool) (string, error) {
	logx.Debug("Making stream response: id %s, content %s, finish %v", id, content, isFinish)

	wx, _, err := NewWXBizJsonMsgCrypt(c.Token, c.EncodingAESKey, "")
	if err != nil {
		return "", err
	}

	resp := UserResp{
		Msgtype: "stream",
		Stream: Stream{
			Id:      id,
			Finish:  isFinish,
			Content: content,
			MsgItem: nil,
		},
	}

	b, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}

	code, msg := wx.EncryptMsg(string(b), nonce)
	if code != WXBizMsgCrypt_OK {
		logx.Error("MakeStreamResp failed: code %d", code)
		return "", getErrorMessage(code)
	}

	return msg, nil
}
