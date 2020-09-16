package handler

import (
	"fmt"
	"strings"

	"github.com/qiniu/x/xlog"
	"github.com/qrtc/qlive/protocol"
)

// MockAccount 模拟的账号服务。
type MockAccount struct {
	accounts []*protocol.Account
}

func (m *MockAccount) GetAccountByPhoneNumber(xl *xlog.Logger, phoneNumber string) (*protocol.Account, error) {
	for _, account := range m.accounts {
		if account.PhoneNumber == phoneNumber {
			return account, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (m *MockAccount) GetAccountByID(xl *xlog.Logger, id string) (*protocol.Account, error) {
	for _, account := range m.accounts {
		if account.ID == id {
			return account, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (m *MockAccount) CreateAccount(xl *xlog.Logger, account *protocol.Account) error {
	if account.ID == "" || account.PhoneNumber == "" {
		return fmt.Errorf("bad request")
	}
	for _, a := range m.accounts {
		if a.ID == account.ID || a.PhoneNumber == account.PhoneNumber {
			return fmt.Errorf("conflict")
		}
	}
	m.accounts = append(m.accounts, account)
	return nil
}

func (m *MockAccount) UpdateAccount(xl *xlog.Logger, id string, account *protocol.Account) (*protocol.Account, error) {
	if account.ID != "" && account.ID != id {
		return nil, fmt.Errorf("bad request")
	}
	var oldAccount *protocol.Account
	for _, a := range m.accounts {
		if a.ID == id {
			oldAccount = a
			break
		}
	}
	if oldAccount == nil {
		return nil, fmt.Errorf("not found")
	}
	if account.PhoneNumber != "" && account.PhoneNumber != oldAccount.PhoneNumber {
		return nil, fmt.Errorf("bad request")
	}
	oldAccount.Nickname = account.Nickname
	oldAccount.Gender = account.Gender
	return oldAccount, nil
}

// AccountLogin 模拟账号登录，返回token。
func (m *MockAccount) AccountLogin(xl *xlog.Logger, id string) (token string, err error) {
	return id + "#" + "login-token", nil
}

// AccountLogout 模拟账号退出登录。
func (m *MockAccount) AccountLogout(xl *xlog.Logger, id string) error {
	return nil
}

// MockSMSCode 模拟的短信服务。
type MockSMSCode struct {
	// 模拟发送出错的情况。
	NumberToError map[string]error
}

// Send 模拟发送验证码
func (m *MockSMSCode) Send(xl *xlog.Logger, phoneNumber string) error {
	if m.NumberToError != nil {
		return m.NumberToError[phoneNumber]
	}
	return nil
}

// Validate 模拟检查输入的验证码。
func (m *MockSMSCode) Validate(xl *xlog.Logger, phoneNumber string, smsCode string) error {
	if smsCode == "123456" {
		return nil
	}
	return fmt.Errorf("wrong sms code")
}

// MockAuth 模拟的认证服务。
type MockAuth struct{}

// GetIDByToken 从token 中获取用户ID。
func (m *MockAuth) GetIDByToken(xl *xlog.Logger, token string) (string, error) {
	parts := strings.SplitN(token, "#", 2)
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid token")
	}
	return parts[0], nil
}

// MockRoom 模拟的房间管理服务。
type MockRoom struct{}
