package logic

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"ai-roleplay/common/response"
	"ai-roleplay/services/speech/api/internal/svc"
	"ai-roleplay/services/speech/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SpeechToTextLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSpeechToTextLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SpeechToTextLogic {
	return &SpeechToTextLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SpeechToTextLogic) SpeechToText(req *types.SttRequest) (resp *types.SttResponse, err error) {
	// 1. 参数验证
	if err := l.validateSTTRequest(req); err != nil {
		return &types.SttResponse{
			Code: response.INVALID_PARAMS,
			Msg:  err.Error(),
		}, nil
	}

	// 2. 检查用户权限和频率限制
	userId := l.ctx.Value("userId").(int64)
	if err := l.checkRateLimit(userId); err != nil {
		return &types.SttResponse{
			Code: response.RATE_LIMIT_EXCEEDED,
			Msg:  "语音识别请求过于频繁",
		}, nil
	}

	// 3. 处理音频数据
	audioData, err := l.processAudioData(req.AudioData, req.AudioFormat)
	if err != nil {
		logx.Errorf("音频数据处理失败: %v", err)
		return &types.SttResponse{
			Code: response.AUDIO_PROCESS_ERROR,
			Msg:  "音频数据格式错误",
		}, nil
	}

	// 4. 音频质量检查
	if err := l.validateAudioQuality(audioData); err != nil {
		return &types.SttResponse{
			Code: response.AUDIO_QUALITY_ERROR,
			Msg:  err.Error(),
		}, nil
	}

	// 5. 调用ASR服务
	startTime := time.Now()
	recognitionResult, err := l.callASRService(audioData, req.Language)
	if err != nil {
		logx.Errorf("ASR服务调用失败: %v", err)
		return &types.SttResponse{
			Code: response.ASR_SERVICE_ERROR,
			Msg:  "语音识别服务暂时不可用",
		}, nil
	}

	// 6. 后处理识别结果
	processedText := l.postProcessRecognitionResult(recognitionResult.Text)

	// 7. 内容安全检查
	if err := l.contentSafetyCheck(processedText); err != nil {
		logx.Warnf("语音识别结果触发安全检查: %s", processedText)
		return &types.SttResponse{
			Code: response.CONTENT_UNSAFE,
			Msg:  "识别内容包含不当信息",
		}, nil
	}

	// 8. 记录使用统计
	duration := time.Since(startTime)
	go l.recordSTTMetrics(userId, len(audioData), len(processedText), duration, recognitionResult.Confidence)

	// 9. 缓存结果（可选）
	if l.svcCtx.Config.Cache.EnableSTTCache {
		go l.cacheSTTResult(audioData, processedText, recognitionResult.Confidence)
	}

	return &types.SttResponse{
		Code: response.SUCCESS,
		Msg:  "识别成功",
		Data: &types.SttData{
			Text:       processedText,
			Confidence: recognitionResult.Confidence,
			Language:   recognitionResult.Language,
			Duration:   duration.Milliseconds(),
			Words:      recognitionResult.Words,
		},
	}, nil
}

// ASR识别结果结构
type ASRResult struct {
	Text       string            `json:"text"`
	Confidence float64           `json:"confidence"`
	Language   string            `json:"language"`
	Words      []*types.WordInfo `json:"words,omitempty"`
}

// 验证STT请求
func (l *SpeechToTextLogic) validateSTTRequest(req *types.SttRequest) error {
	if req.AudioData == "" {
		return fmt.Errorf("音频数据不能为空")
	}

	// 检查Base64格式
	if _, err := base64.StdEncoding.DecodeString(req.AudioData); err != nil {
		return fmt.Errorf("音频数据格式错误，需要Base64编码")
	}

	// 检查语言支持
	supportedLanguages := []string{"zh-CN", "en-US", "ja-JP", "ko-KR"}
	if req.Language == "" {
		req.Language = "zh-CN" // 默认中文
	}

	supported := false
	for _, lang := range supportedLanguages {
		if req.Language == lang {
			supported = true
			break
		}
	}
	if !supported {
		return fmt.Errorf("不支持的语言: %s", req.Language)
	}

	return nil
}

// 检查频率限制
func (l *SpeechToTextLogic) checkRateLimit(userId int64) error {
	key := fmt.Sprintf("stt_rate_limit:user:%d", userId)

	count, err := l.svcCtx.Redis.Incr(l.ctx, key).Result()
	if err != nil {
		logx.Errorf("Redis操作失败: %v", err)
		return nil // Redis失败时不限制
	}

	if count == 1 {
		l.svcCtx.Redis.Expire(l.ctx, key, time.Minute)
	}

	// 每分钟最多30次STT请求
	if count > 30 {
		return fmt.Errorf("语音识别请求频率超限")
	}

	return nil
}

// 处理音频数据
func (l *SpeechToTextLogic) processAudioData(audioDataBase64, format string) ([]byte, error) {
	// 解码Base64
	audioData, err := base64.StdEncoding.DecodeString(audioDataBase64)
	if err != nil {
		return nil, fmt.Errorf("Base64解码失败: %w", err)
	}

	// 检查文件大小（最大10MB）
	maxSize := 10 * 1024 * 1024
	if len(audioData) > maxSize {
		return nil, fmt.Errorf("音频文件过大，最大支持10MB")
	}

	// 根据格式进行转换（如果需要）
	if format != "" && format != "wav" && format != "mp3" && format != "m4a" {
		return nil, fmt.Errorf("不支持的音频格式: %s", format)
	}

	return audioData, nil
}

// 验证音频质量
func (l *SpeechToTextLogic) validateAudioQuality(audioData []byte) error {
	// 检查音频时长（最少0.5秒，最多60秒）
	// 这里简化处理，实际需要解析音频文件头
	minSize := 1000   // 大约0.5秒的音频
	maxSize := 600000 // 大约60秒的音频

	if len(audioData) < minSize {
		return fmt.Errorf("音频时长过短，至少需要0.5秒")
	}

	if len(audioData) > maxSize {
		return fmt.Errorf("音频时长过长，最多支持60秒")
	}

	return nil
}

// 调用ASR服务
func (l *SpeechToTextLogic) callASRService(audioData []byte, language string) (*ASRResult, error) {
	// 根据配置选择ASR服务提供商
	provider := l.svcCtx.Config.ASR.Provider

	switch provider {
	case "azure":
		return l.callAzureASR(audioData, language)
	case "aliyun":
		return l.callAliyunASR(audioData, language)
	case "tencent":
		return l.callTencentASR(audioData, language)
	case "baidu":
		return l.callBaiduASR(audioData, language)
	default:
		return l.callDefaultASR(audioData, language)
	}
}

// 调用Azure ASR服务
func (l *SpeechToTextLogic) callAzureASR(audioData []byte, language string) (*ASRResult, error) {
	// Azure Speech Service API调用
	endpoint := l.svcCtx.Config.ASR.Azure.Endpoint
	apiKey := l.svcCtx.Config.ASR.Azure.ApiKey

	// 构建请求URL
	url := fmt.Sprintf("%s/speech/recognition/conversation/cognitiveservices/v1?language=%s", endpoint, language)

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(l.ctx, "POST", url, bytes.NewReader(audioData))
	if err != nil {
		return nil, err
	}

	// 设置请求头
	req.Header.Set("Ocp-Apim-Subscription-Key", apiKey)
	req.Header.Set("Content-Type", "audio/wav")
	req.Header.Set("Accept", "application/json")

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Azure ASR请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Azure ASR返回错误: %d, %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var azureResp struct {
		RecognitionStatus string  `json:"RecognitionStatus"`
		DisplayText       string  `json:"DisplayText"`
		Confidence        float64 `json:"Confidence"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&azureResp); err != nil {
		return nil, fmt.Errorf("解析Azure ASR响应失败: %w", err)
	}

	if azureResp.RecognitionStatus != "Success" {
		return nil, fmt.Errorf("Azure ASR识别失败: %s", azureResp.RecognitionStatus)
	}

	return &ASRResult{
		Text:       azureResp.DisplayText,
		Confidence: azureResp.Confidence,
		Language:   language,
	}, nil
}

// 调用默认ASR服务（Web Speech API模拟）
func (l *SpeechToTextLogic) callDefaultASR(audioData []byte, language string) (*ASRResult, error) {
	// 模拟ASR识别过程
	time.Sleep(500 * time.Millisecond)

	// 根据音频数据长度模拟不同的识别结果
	var text string
	var confidence float64

	switch {
	case len(audioData) < 5000:
		text = "你好"
		confidence = 0.95
	case len(audioData) < 15000:
		text = "今天天气真不错"
		confidence = 0.90
	case len(audioData) < 25000:
		text = "我想和你聊聊关于人工智能的话题"
		confidence = 0.88
	default:
		text = "这是一段比较长的语音内容，包含了多个句子和丰富的信息"
		confidence = 0.85
	}

	return &ASRResult{
		Text:       text,
		Confidence: confidence,
		Language:   language,
	}, nil
}

// 后处理识别结果
func (l *SpeechToTextLogic) postProcessRecognitionResult(text string) string {
	// 1. 去除首尾空白
	text = strings.TrimSpace(text)

	// 2. 标点符号规范化
	text = strings.ReplaceAll(text, "。。", "。")
	text = strings.ReplaceAll(text, "，，", "，")
	text = strings.ReplaceAll(text, "？？", "？")
	text = strings.ReplaceAll(text, "！！", "！")

	// 3. 过滤特殊字符
	// 保留中文、英文、数字、常用标点
	// 这里简化处理，实际可能需要更复杂的正则表达式

	// 4. 长度限制
	maxLength := 500
	if len([]rune(text)) > maxLength {
		runes := []rune(text)
		text = string(runes[:maxLength])
	}

	return text
}

// 内容安全检查
func (l *SpeechToTextLogic) contentSafetyCheck(text string) error {
	// 基础敏感词检查
	sensitiveWords := []string{"暴力", "色情", "政治敏感"}
	lowerText := strings.ToLower(text)

	for _, word := range sensitiveWords {
		if strings.Contains(lowerText, strings.ToLower(word)) {
			return fmt.Errorf("内容包含敏感词汇")
		}
	}

	return nil
}

// 记录STT使用统计
func (l *SpeechToTextLogic) recordSTTMetrics(userId int64, audioSize, textLength int, duration time.Duration, confidence float64) {
	metrics := &model.STTMetrics{
		UserId:     userId,
		AudioSize:  audioSize,
		TextLength: textLength,
		Duration:   duration.Milliseconds(),
		Confidence: confidence,
		CreatedAt:  time.Now(),
	}

	if err := l.svcCtx.STTMetricsModel.Insert(context.Background(), metrics); err != nil {
		logx.Errorf("记录STT统计失败: %v", err)
	}
}

// 缓存STT结果
func (l *SpeechToTextLogic) cacheSTTResult(audioData []byte, text string, confidence float64) {
	// 计算音频数据哈希作为缓存key
	hash := l.calculateAudioHash(audioData)
	cacheKey := fmt.Sprintf("stt:result:%s", hash)

	result := map[string]interface{}{
		"text":       text,
		"confidence": confidence,
		"timestamp":  time.Now().Unix(),
	}

	// 缓存1小时
	if err := l.svcCtx.CacheManager.SetSTTResult(context.Background(), cacheKey, result, time.Hour); err != nil {
		logx.Errorf("缓存STT结果失败: %v", err)
	}
}

// 计算音频数据哈希
func (l *SpeechToTextLogic) calculateAudioHash(audioData []byte) string {
	// 使用MD5计算哈希（简化处理）
	return fmt.Sprintf("%x", md5.Sum(audioData))
}
