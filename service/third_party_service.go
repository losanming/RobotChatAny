package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"new-robot-self/model"
	"new-robot-self/public"
	"strconv"
	"time"
)

// ThirdPartyService 第三方登录服务
type ThirdPartyService struct {
	UserService UserService
}

// GiteeUserInfo Gitee用户信息
type GiteeUserInfo struct {
	Id        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

// HandleGiteeLogin 处理Gitee登录
func (s *ThirdPartyService) HandleGiteeLogin(code string) (*model.User, error) {
	// 1. 用授权码换取访问令牌
	tokenData, err := s.getGiteeAccessToken(code)
	if err != nil {
		return nil, err
	}

	// 2. 获取Gitee用户信息
	giteeUser, err := s.getGiteeUserInfo(tokenData["access_token"].(string))
	if err != nil {
		return nil, err
	}

	// 3. 检查是否已有关联账号
	thirdParty, err := s.getUserByGiteeId(strconv.Itoa(giteeUser.Id))
	if err == nil {
		// 已存在关联，直接返回用户
		user, err := s.UserService.GetUserByCode(thirdParty.UserCode)
		if err != nil {
			return nil, errors.New("关联账号异常")
		}

		// 更新第三方账号信息
		thirdParty.AccessToken = tokenData["access_token"].(string)
		thirdParty.RefreshToken = tokenData["refresh_token"].(string)
		expiresIn, _ := strconv.ParseInt(tokenData["expires_in"].(string), 10, 64)
		thirdParty.ExpiresAt = time.Now().Unix() + expiresIn
		thirdParty.PlatformUsername = giteeUser.Name
		thirdParty.PlatformAvatar = giteeUser.AvatarUrl
		s.updateThirdPartyInfo(thirdParty)

		return user, nil
	}

	// 4. 如果有手机号，尝试关联现有账号
	var user *model.User
	if giteeUser.Phone != "" {
		user, err = s.UserService.GetUserByMobile(giteeUser.Phone)
		if err == nil {
			// 找到现有账号，创建关联
			s.createThirdPartyRelation(user.Code, "gitee", giteeUser, tokenData)
			return user, nil
		}
	}

	// 5. 没有找到关联账号，创建新账号
	// 注意：如果Gitee不提供手机号，这里需要额外的流程让用户提供手机号
	if giteeUser.Phone == "" {
		return nil, errors.New("需要绑定手机号才能完成注册")
	}

	// 创建新用户
	newUser := &model.User{
		Username:    "gitee_" + giteeUser.Login,
		Mobile:      giteeUser.Phone,
		Email:       giteeUser.Email,
		Nickname:    giteeUser.Name,
		Avatar:      giteeUser.AvatarUrl,
		Status:      1,
		CreatedAt:   time.Now().Unix(),
		LastLoginIp: "", // 这里需要从请求中获取IP
	}

	// 保存新用户
	if err := s.UserService.Create(newUser); err != nil {
		return nil, err
	}

	// 创建第三方关联
	s.createThirdPartyRelation(newUser.Code, "gitee", giteeUser, tokenData)

	return newUser, nil
}

// getGiteeAccessToken 获取Gitee访问令牌
func (s *ThirdPartyService) getGiteeAccessToken(code string) (map[string]interface{}, error) {
	giteeConfig := public.Config.Gitee

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("client_id", giteeConfig.ClientId)
	data.Set("client_secret", giteeConfig.ClientSecret)
	data.Set("redirect_uri", giteeConfig.RedirectUri)

	resp, err := http.PostForm("https://gitee.com/oauth/token", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if _, ok := result["access_token"]; !ok {
		return nil, errors.New("获取访问令牌失败")
	}

	return result, nil
}

// getGiteeUserInfo 获取Gitee用户信息
func (s *ThirdPartyService) getGiteeUserInfo(accessToken string) (*GiteeUserInfo, error) {
	req, err := http.NewRequest("GET", "https://gitee.com/api/v5/user", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("access_token", accessToken)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo GiteeUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// getUserByGiteeId 通过Gitee ID获取关联用户
func (s *ThirdPartyService) getUserByGiteeId(giteeId string) (*model.UserThirdParty, error) {
	var thirdParty model.UserThirdParty
	has, err := public.DB.Where("platform = ? AND platform_user_id = ?", "gitee", giteeId).Get(&thirdParty)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("未找到关联账号")
	}
	return &thirdParty, nil
}

// updateThirdPartyInfo 更新第三方账号信息
func (s *ThirdPartyService) updateThirdPartyInfo(thirdParty *model.UserThirdParty) error {
	_, err := public.DB.ID(thirdParty.Id).Update(thirdParty)
	return err
}

// createThirdPartyRelation 创建第三方账号关联
func (s *ThirdPartyService) createThirdPartyRelation(userCode string, platform string, giteeUser *GiteeUserInfo, tokenData map[string]interface{}) error {
	expiresIn, _ := strconv.ParseInt(tokenData["expires_in"].(string), 10, 64)

	thirdParty := &model.UserThirdParty{
		Code:             userCode,
		Platform:         platform,
		PlatformUserId:   strconv.Itoa(giteeUser.Id),
		PlatformUsername: giteeUser.Name,
		PlatformAvatar:   giteeUser.AvatarUrl,
		AccessToken:      tokenData["access_token"].(string),
		RefreshToken:     tokenData["refresh_token"].(string),
		ExpiresAt:        time.Now().Unix() + expiresIn,
	}

	_, err := public.DB.Insert(thirdParty)
	return err
}
