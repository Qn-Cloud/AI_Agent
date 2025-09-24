package logic

import (
	"context"
	"errors"
	"time"

	"ai-roleplay/common/response"
	"ai-roleplay/common/utils"
	"ai-roleplay/services/user/api/internal/svc"
	"ai-roleplay/services/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	// 1. 参数验证
	if err := l.validateLoginRequest(req); err != nil {
		return &types.LoginResponse{
			Code: response.INVALID_PARAMS,
			Msg:  err.Error(),
		}, nil
	}

	// 2. 查找用户
	user, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)
	if err != nil {
		// 用户不存在或数据库错误
		return &types.LoginResponse{
			Code: response.USER_NOT_FOUND,
			Msg:  "用户名或密码错误",
		}, nil
	}

	// 3. 检查用户状态
	if user.Status != 1 {
		return &types.LoginResponse{
			Code: response.USER_DISABLED,
			Msg:  "账户已被禁用",
		}, nil
	}

	// 4. 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		// 记录登录失败
		l.logLoginAttempt(user.Id, false, "密码错误")

		return &types.LoginResponse{
			Code: response.INVALID_PASSWORD,
			Msg:  "用户名或密码错误",
		}, nil
	}

	// 5. 生成JWT Token
	now := time.Now()
	accessTokenExpire := l.svcCtx.Config.Auth.AccessExpire
	refreshTokenExpire := accessTokenExpire * 7 // 刷新token有效期为7倍访问token

	accessToken, err := utils.GenerateJWTToken(utils.JWTPayload{
		UserId:   user.Id,
		Username: user.Username,
		ExpireAt: now.Add(time.Second * time.Duration(accessTokenExpire)).Unix(),
	}, l.svcCtx.Config.Auth.AccessSecret)

	if err != nil {
		logx.Errorf("生成access token失败: %v", err)
		return &types.LoginResponse{
			Code: response.INTERNAL_ERROR,
			Msg:  "登录失败，请重试",
		}, nil
	}

	// 6. 生成刷新Token（如果记住登录）
	var refreshToken string
	if req.Remember {
		refreshToken, err = utils.GenerateJWTToken(utils.JWTPayload{
			UserId:   user.Id,
			Username: user.Username,
			ExpireAt: now.Add(time.Second * time.Duration(refreshTokenExpire)).Unix(),
		}, l.svcCtx.Config.Auth.RefreshSecret)

		if err != nil {
			logx.Errorf("生成refresh token失败: %v", err)
			// 刷新token生成失败不影响登录
		}
	}

	// 7. 更新最后登录时间
	err = l.svcCtx.UserModel.UpdateLastLoginTime(l.ctx, user.Id, now)
	if err != nil {
		logx.Errorf("更新用户最后登录时间失败: %v", err)
		// 不影响登录流程
	}

	// 8. 缓存用户会话信息
	sessionInfo := &utils.UserSession{
		UserId:    user.Id,
		Username:  user.Username,
		LoginTime: now,
		ExpireAt:  now.Add(time.Second * time.Duration(accessTokenExpire)),
	}

	if err := l.svcCtx.CacheManager.SetUserSession(l.ctx, accessToken, sessionInfo); err != nil {
		logx.Errorf("缓存用户会话失败: %v", err)
		// 不影响登录流程
	}

	// 9. 记录登录成功
	l.logLoginAttempt(user.Id, true, "登录成功")

	// 10. 构造响应
	userInfo := &types.UserInfo{
		Id:        user.Id,
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Bio:       user.Bio,
		Status:    int(user.Status),
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return &types.LoginResponse{
		Code:         response.SUCCESS,
		Msg:          "登录成功",
		Data:         userInfo,
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// 验证登录请求参数
func (l *LoginLogic) validateLoginRequest(req *types.LoginRequest) error {
	if len(req.Username) < 3 || len(req.Username) > 20 {
		return errors.New("用户名长度必须在3-20个字符之间")
	}

	if len(req.Password) < 6 || len(req.Password) > 20 {
		return errors.New("密码长度必须在6-20个字符之间")
	}

	return nil
}

// 记录登录尝试
func (l *LoginLogic) logLoginAttempt(userId int64, success bool, message string) {
	attempt := &model.LoginAttempt{
		UserId:    userId,
		Success:   success,
		Message:   message,
		IpAddress: l.getClientIP(),
		UserAgent: l.getUserAgent(),
		CreatedAt: time.Now(),
	}

	// 异步记录，不影响主流程
	go func() {
		if err := l.svcCtx.LoginAttemptModel.Insert(context.Background(), attempt); err != nil {
			logx.Errorf("记录登录尝试失败: %v", err)
		}
	}()
}

// 获取客户端IP
func (l *LoginLogic) getClientIP() string {
	// 从context中获取HTTP请求信息
	// 这里需要根据实际的中间件实现来获取
	return "127.0.0.1" // 占位符
}

// 获取用户代理
func (l *LoginLogic) getUserAgent() string {
	// 从context中获取HTTP请求信息
	return "Unknown" // 占位符
}
