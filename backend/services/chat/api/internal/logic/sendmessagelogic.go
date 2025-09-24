package logic

import (
	"context"
	"fmt"
	"strings"
	"time"

	"ai-roleplay/common/response"
	"ai-roleplay/services/chat/api/internal/svc"
	"ai-roleplay/services/chat/api/internal/types"

	"ai-roleplay/services/ai/rpc/airpc"
	"ai-roleplay/services/character/rpc/characterrpc"
	"ai-roleplay/services/conversation/model"
	"ai-roleplay/services/speech/rpc/speechrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendMessageLogic) SendMessage(req *types.SendMessageRequest) (resp *types.SendMessageResponse, err error) {
	// 1. 获取当前用户ID
	userId := l.ctx.Value("userId").(int64)

	// 2. 参数验证
	if err := l.validateSendMessageRequest(req); err != nil {
		return &types.SendMessageResponse{
			Code: response.INVALID_PARAMS,
			Msg:  err.Error(),
		}, nil
	}

	// 3. 验证对话是否存在且属于当前用户
	conversation, err := l.svcCtx.ConversationModel.FindOne(l.ctx, req.ConversationId)
	if err != nil {
		return &types.SendMessageResponse{
			Code: response.CONVERSATION_NOT_FOUND,
			Msg:  "对话不存在",
		}, nil
	}

	if conversation.UserId != userId {
		return &types.SendMessageResponse{
			Code: response.PERMISSION_DENIED,
			Msg:  "无权限访问此对话",
		}, nil
	}

	// 4. 处理语音消息
	var audioUrl string
	var duration int64
	if req.Type == "audio" && req.AudioData != "" {
		// 调用语音识别服务
		sttResp, err := l.svcCtx.SpeechRpc.SpeechToText(l.ctx, &speechrpc.SttRequest{
			AudioData: req.AudioData,
			Language:  "zh-CN",
		})

		if err != nil {
			logx.Errorf("语音识别失败: %v", err)
			return &types.SendMessageResponse{
				Code: response.SPEECH_RECOGNITION_ERROR,
				Msg:  "语音识别失败",
			}, nil
		}

		req.Content = sttResp.Text
		audioUrl = req.AudioUrl
		duration = req.AudioDuration
	}

	// 5. 内容安全检查
	if err := l.contentSafetyCheck(req.Content); err != nil {
		return &types.SendMessageResponse{
			Code: response.CONTENT_UNSAFE,
			Msg:  "消息内容包含不当信息",
		}, nil
	}

	// 6. 保存用户消息
	now := time.Now()
	userMessage := &model.Message{
		ConversationId: req.ConversationId,
		UserId:         userId,
		Type:           "user",
		Content:        req.Content,
		AudioUrl:       audioUrl,
		AudioDuration:  duration,
		CreatedAt:      now,
	}

	err = l.svcCtx.MessageModel.Insert(l.ctx, userMessage)
	if err != nil {
		logx.Errorf("保存用户消息失败: %v", err)
		return &types.SendMessageResponse{
			Code: response.DATABASE_ERROR,
			Msg:  "保存消息失败",
		}, nil
	}

	// 7. 获取角色信息
	character, err := l.svcCtx.CharacterRpc.GetCharacter(l.ctx, &characterrpc.GetCharacterRequest{
		Id: conversation.CharacterId,
	})
	if err != nil {
		logx.Errorf("获取角色信息失败: %v", err)
		return &types.SendMessageResponse{
			Code: response.CHARACTER_NOT_FOUND,
			Msg:  "角色不存在",
		}, nil
	}

	// 8. 调用AI服务生成回复
	aiResp, err := l.callAIService(req.ConversationId, character, req.Content)
	if err != nil {
		logx.Errorf("AI生成回复失败: %v", err)
		return &types.SendMessageResponse{
			Code: response.AI_SERVICE_ERROR,
			Msg:  "AI服务暂时不可用",
		}, nil
	}

	// 9. 保存AI回复消息
	aiMessage := &model.Message{
		ConversationId: req.ConversationId,
		UserId:         userId,
		Type:           "ai",
		Content:        aiResp.Content,
		CreatedAt:      time.Now(),
	}

	err = l.svcCtx.MessageModel.Insert(l.ctx, aiMessage)
	if err != nil {
		logx.Errorf("保存AI消息失败: %v", err)
		// AI消息保存失败不影响用户消息
	}

	// 10. 更新对话最后更新时间和消息数量
	err = l.svcCtx.ConversationModel.UpdateLastMessage(l.ctx, req.ConversationId, time.Now())
	if err != nil {
		logx.Errorf("更新对话信息失败: %v", err)
	}

	// 11. 异步生成语音（如果需要）
	var aiAudioUrl string
	if req.NeedTTS {
		go l.generateTTSAsync(aiResp.Content, character.VoiceSettings)
	}

	// 12. 构造响应
	return &types.SendMessageResponse{
		Code: response.SUCCESS,
		Msg:  "发送成功",
		Data: &types.SendMessageData{
			UserMessage: &types.Message{
				Id:            userMessage.Id,
				Type:          userMessage.Type,
				Content:       userMessage.Content,
				AudioUrl:      userMessage.AudioUrl,
				AudioDuration: userMessage.AudioDuration,
				Timestamp:     userMessage.CreatedAt.Format("2006-01-02 15:04:05"),
			},
			AiMessage: &types.Message{
				Id:            aiMessage.Id,
				Type:          aiMessage.Type,
				Content:       aiMessage.Content,
				AudioUrl:      aiAudioUrl,
				AudioDuration: 0,
				Timestamp:     aiMessage.CreatedAt.Format("2006-01-02 15:04:05"),
			},
		},
	}, nil
}

// 验证发送消息请求
func (l *SendMessageLogic) validateSendMessageRequest(req *types.SendMessageRequest) error {
	if req.ConversationId <= 0 {
		return fmt.Errorf("对话ID无效")
	}

	if req.Type != "text" && req.Type != "audio" {
		return fmt.Errorf("消息类型无效")
	}

	if req.Type == "text" && len(req.Content) == 0 {
		return fmt.Errorf("文字消息内容不能为空")
	}

	if req.Type == "audio" && req.AudioData == "" {
		return fmt.Errorf("语音消息数据不能为空")
	}

	if len(req.Content) > 2000 {
		return fmt.Errorf("消息内容过长，最多2000个字符")
	}

	return nil
}

// 内容安全检查
func (l *SendMessageLogic) contentSafetyCheck(content string) error {
	// 调用内容安全检查服务
	// 这里可以集成第三方内容审核服务

	// 简单的敏感词检查
	sensitiveWords := []string{"政治敏感词", "色情", "暴力"}
	for _, word := range sensitiveWords {
		if strings.Contains(content, word) {
			return fmt.Errorf("内容包含敏感信息")
		}
	}

	return nil
}

// 调用AI服务
func (l *SendMessageLogic) callAIService(conversationId int64, character *characterrpc.Character, userInput string) (*airpc.ChatResponse, error) {
	// 获取对话历史上下文
	context, err := l.getConversationContext(conversationId)
	if err != nil {
		return nil, err
	}

	// 构建AI请求
	aiReq := &airpc.ChatRequest{
		ConversationId:  conversationId,
		CharacterId:     character.Id,
		CharacterPrompt: character.Prompt,
		UserInput:       userInput,
		Context:         context,
		MaxTokens:       1000,
		Temperature:     0.7,
	}

	// 调用AI服务
	return l.svcCtx.AiRpc.Chat(l.ctx, aiReq)
}

// 获取对话上下文
func (l *SendMessageLogic) getConversationContext(conversationId int64) ([]*airpc.Message, error) {
	// 从缓存获取最近的对话历史
	cachedContext, err := l.svcCtx.CacheManager.GetConversationContext(l.ctx, conversationId)
	if err == nil && len(cachedContext) > 0 {
		return l.convertToAIMessages(cachedContext), nil
	}

	// 从数据库获取最近10条消息
	messages, err := l.svcCtx.MessageModel.FindRecentByConversationId(l.ctx, conversationId, 10)
	if err != nil {
		return nil, err
	}

	// 缓存上下文
	go func() {
		if err := l.svcCtx.CacheManager.SetConversationContext(context.Background(), conversationId, messages); err != nil {
			logx.Errorf("缓存对话上下文失败: %v", err)
		}
	}()

	return l.convertToAIMessages(messages), nil
}

// 转换消息格式
func (l *SendMessageLogic) convertToAIMessages(messages []*model.Message) []*airpc.Message {
	aiMessages := make([]*airpc.Message, 0, len(messages))

	for _, msg := range messages {
		aiMessages = append(aiMessages, &airpc.Message{
			Role:    msg.Type,
			Content: msg.Content,
		})
	}

	return aiMessages
}

// 异步生成TTS
func (l *SendMessageLogic) generateTTSAsync(text string, voiceSettings map[string]interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ttsReq := &speechrpc.TtsRequest{
		Text:   text,
		Voice:  "zh-CN-XiaoxiaoNeural", // 默认音色
		Speed:  1.0,
		Pitch:  1.0,
		Volume: 0.8,
	}

	// 应用角色语音设置
	if rate, ok := voiceSettings["rate"].(float64); ok {
		ttsReq.Speed = rate
	}
	if pitch, ok := voiceSettings["pitch"].(float64); ok {
		ttsReq.Pitch = pitch
	}
	if volume, ok := voiceSettings["volume"].(float64); ok {
		ttsReq.Volume = volume
	}

	resp, err := l.svcCtx.SpeechRpc.TextToSpeech(ctx, ttsReq)
	if err != nil {
		logx.Errorf("TTS生成失败: %v", err)
		return
	}

	logx.Infof("TTS生成成功，音频URL: %s", resp.AudioUrl)
}
