// Code generated by mockery v2.46.3. DO NOT EDIT.

package mock

import (
	context "context"

	dto "github.com/mandarine-io/Backend/internal/domain/dto"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"

	mock "github.com/stretchr/testify/mock"

	oauth "github.com/mandarine-io/Backend/pkg/oauth"
)

// AuthServiceMock is an autogenerated mock type for the AuthService type
type AuthServiceMock struct {
	mock.Mock
}

type AuthServiceMock_Expecter struct {
	mock *mock.Mock
}

func (_m *AuthServiceMock) EXPECT() *AuthServiceMock_Expecter {
	return &AuthServiceMock_Expecter{mock: &_m.Mock}
}

// FetchUserInfo provides a mock function with given fields: ctx, provider, input
func (_m *AuthServiceMock) FetchUserInfo(ctx context.Context, provider string, input dto.FetchUserInfoInput) (oauth.UserInfo, error) {
	ret := _m.Called(ctx, provider, input)

	if len(ret) == 0 {
		panic("no return value specified for FetchUserInfo")
	}

	var r0 oauth.UserInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, dto.FetchUserInfoInput) (oauth.UserInfo, error)); ok {
		return rf(ctx, provider, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, dto.FetchUserInfoInput) oauth.UserInfo); ok {
		r0 = rf(ctx, provider, input)
	} else {
		r0 = ret.Get(0).(oauth.UserInfo)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, dto.FetchUserInfoInput) error); ok {
		r1 = rf(ctx, provider, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthServiceMock_FetchUserInfo_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FetchUserInfo'
type AuthServiceMock_FetchUserInfo_Call struct {
	*mock.Call
}

// FetchUserInfo is a helper method to define mock.On call
//   - ctx context.Context
//   - provider string
//   - input dto.FetchUserInfoInput
func (_e *AuthServiceMock_Expecter) FetchUserInfo(ctx interface{}, provider interface{}, input interface{}) *AuthServiceMock_FetchUserInfo_Call {
	return &AuthServiceMock_FetchUserInfo_Call{Call: _e.mock.On("FetchUserInfo", ctx, provider, input)}
}

func (_c *AuthServiceMock_FetchUserInfo_Call) Run(run func(ctx context.Context, provider string, input dto.FetchUserInfoInput)) *AuthServiceMock_FetchUserInfo_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(dto.FetchUserInfoInput))
	})
	return _c
}

func (_c *AuthServiceMock_FetchUserInfo_Call) Return(_a0 oauth.UserInfo, _a1 error) *AuthServiceMock_FetchUserInfo_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AuthServiceMock_FetchUserInfo_Call) RunAndReturn(run func(context.Context, string, dto.FetchUserInfoInput) (oauth.UserInfo, error)) *AuthServiceMock_FetchUserInfo_Call {
	_c.Call.Return(run)
	return _c
}

// GetConsentPageUrl provides a mock function with given fields: _a0, provider, redirectUrl
func (_m *AuthServiceMock) GetConsentPageUrl(_a0 context.Context, provider string, redirectUrl string) (dto.GetConsentPageUrlOutput, error) {
	ret := _m.Called(_a0, provider, redirectUrl)

	if len(ret) == 0 {
		panic("no return value specified for GetConsentPageUrl")
	}

	var r0 dto.GetConsentPageUrlOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (dto.GetConsentPageUrlOutput, error)); ok {
		return rf(_a0, provider, redirectUrl)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) dto.GetConsentPageUrlOutput); ok {
		r0 = rf(_a0, provider, redirectUrl)
	} else {
		r0 = ret.Get(0).(dto.GetConsentPageUrlOutput)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(_a0, provider, redirectUrl)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthServiceMock_GetConsentPageUrl_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetConsentPageUrl'
type AuthServiceMock_GetConsentPageUrl_Call struct {
	*mock.Call
}

// GetConsentPageUrl is a helper method to define mock.On call
//   - _a0 context.Context
//   - provider string
//   - redirectUrl string
func (_e *AuthServiceMock_Expecter) GetConsentPageUrl(_a0 interface{}, provider interface{}, redirectUrl interface{}) *AuthServiceMock_GetConsentPageUrl_Call {
	return &AuthServiceMock_GetConsentPageUrl_Call{Call: _e.mock.On("GetConsentPageUrl", _a0, provider, redirectUrl)}
}

func (_c *AuthServiceMock_GetConsentPageUrl_Call) Run(run func(_a0 context.Context, provider string, redirectUrl string)) *AuthServiceMock_GetConsentPageUrl_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *AuthServiceMock_GetConsentPageUrl_Call) Return(_a0 dto.GetConsentPageUrlOutput, _a1 error) *AuthServiceMock_GetConsentPageUrl_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AuthServiceMock_GetConsentPageUrl_Call) RunAndReturn(run func(context.Context, string, string) (dto.GetConsentPageUrlOutput, error)) *AuthServiceMock_GetConsentPageUrl_Call {
	_c.Call.Return(run)
	return _c
}

// Login provides a mock function with given fields: ctx, input
func (_m *AuthServiceMock) Login(ctx context.Context, input dto.LoginInput) (dto.JwtTokensOutput, error) {
	ret := _m.Called(ctx, input)

	if len(ret) == 0 {
		panic("no return value specified for Login")
	}

	var r0 dto.JwtTokensOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, dto.LoginInput) (dto.JwtTokensOutput, error)); ok {
		return rf(ctx, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, dto.LoginInput) dto.JwtTokensOutput); ok {
		r0 = rf(ctx, input)
	} else {
		r0 = ret.Get(0).(dto.JwtTokensOutput)
	}

	if rf, ok := ret.Get(1).(func(context.Context, dto.LoginInput) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthServiceMock_Login_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Login'
type AuthServiceMock_Login_Call struct {
	*mock.Call
}

// Login is a helper method to define mock.On call
//   - ctx context.Context
//   - input dto.LoginInput
func (_e *AuthServiceMock_Expecter) Login(ctx interface{}, input interface{}) *AuthServiceMock_Login_Call {
	return &AuthServiceMock_Login_Call{Call: _e.mock.On("Login", ctx, input)}
}

func (_c *AuthServiceMock_Login_Call) Run(run func(ctx context.Context, input dto.LoginInput)) *AuthServiceMock_Login_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(dto.LoginInput))
	})
	return _c
}

func (_c *AuthServiceMock_Login_Call) Return(_a0 dto.JwtTokensOutput, _a1 error) *AuthServiceMock_Login_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AuthServiceMock_Login_Call) RunAndReturn(run func(context.Context, dto.LoginInput) (dto.JwtTokensOutput, error)) *AuthServiceMock_Login_Call {
	_c.Call.Return(run)
	return _c
}

// Logout provides a mock function with given fields: ctx, jti
func (_m *AuthServiceMock) Logout(ctx context.Context, jti string) error {
	ret := _m.Called(ctx, jti)

	if len(ret) == 0 {
		panic("no return value specified for Logout")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, jti)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AuthServiceMock_Logout_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Logout'
type AuthServiceMock_Logout_Call struct {
	*mock.Call
}

// Logout is a helper method to define mock.On call
//   - ctx context.Context
//   - jti string
func (_e *AuthServiceMock_Expecter) Logout(ctx interface{}, jti interface{}) *AuthServiceMock_Logout_Call {
	return &AuthServiceMock_Logout_Call{Call: _e.mock.On("Logout", ctx, jti)}
}

func (_c *AuthServiceMock_Logout_Call) Run(run func(ctx context.Context, jti string)) *AuthServiceMock_Logout_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *AuthServiceMock_Logout_Call) Return(_a0 error) *AuthServiceMock_Logout_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AuthServiceMock_Logout_Call) RunAndReturn(run func(context.Context, string) error) *AuthServiceMock_Logout_Call {
	_c.Call.Return(run)
	return _c
}

// RecoveryPassword provides a mock function with given fields: ctx, input, localizer
func (_m *AuthServiceMock) RecoveryPassword(ctx context.Context, input dto.RecoveryPasswordInput, localizer *i18n.Localizer) error {
	ret := _m.Called(ctx, input, localizer)

	if len(ret) == 0 {
		panic("no return value specified for RecoveryPassword")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, dto.RecoveryPasswordInput, *i18n.Localizer) error); ok {
		r0 = rf(ctx, input, localizer)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AuthServiceMock_RecoveryPassword_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RecoveryPassword'
type AuthServiceMock_RecoveryPassword_Call struct {
	*mock.Call
}

// RecoveryPassword is a helper method to define mock.On call
//   - ctx context.Context
//   - input dto.RecoveryPasswordInput
//   - localizer *i18n.Localizer
func (_e *AuthServiceMock_Expecter) RecoveryPassword(ctx interface{}, input interface{}, localizer interface{}) *AuthServiceMock_RecoveryPassword_Call {
	return &AuthServiceMock_RecoveryPassword_Call{Call: _e.mock.On("RecoveryPassword", ctx, input, localizer)}
}

func (_c *AuthServiceMock_RecoveryPassword_Call) Run(run func(ctx context.Context, input dto.RecoveryPasswordInput, localizer *i18n.Localizer)) *AuthServiceMock_RecoveryPassword_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(dto.RecoveryPasswordInput), args[2].(*i18n.Localizer))
	})
	return _c
}

func (_c *AuthServiceMock_RecoveryPassword_Call) Return(_a0 error) *AuthServiceMock_RecoveryPassword_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AuthServiceMock_RecoveryPassword_Call) RunAndReturn(run func(context.Context, dto.RecoveryPasswordInput, *i18n.Localizer) error) *AuthServiceMock_RecoveryPassword_Call {
	_c.Call.Return(run)
	return _c
}

// RefreshTokens provides a mock function with given fields: ctx, refreshToken
func (_m *AuthServiceMock) RefreshTokens(ctx context.Context, refreshToken string) (dto.JwtTokensOutput, error) {
	ret := _m.Called(ctx, refreshToken)

	if len(ret) == 0 {
		panic("no return value specified for RefreshTokens")
	}

	var r0 dto.JwtTokensOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (dto.JwtTokensOutput, error)); ok {
		return rf(ctx, refreshToken)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) dto.JwtTokensOutput); ok {
		r0 = rf(ctx, refreshToken)
	} else {
		r0 = ret.Get(0).(dto.JwtTokensOutput)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, refreshToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthServiceMock_RefreshTokens_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RefreshTokens'
type AuthServiceMock_RefreshTokens_Call struct {
	*mock.Call
}

// RefreshTokens is a helper method to define mock.On call
//   - ctx context.Context
//   - refreshToken string
func (_e *AuthServiceMock_Expecter) RefreshTokens(ctx interface{}, refreshToken interface{}) *AuthServiceMock_RefreshTokens_Call {
	return &AuthServiceMock_RefreshTokens_Call{Call: _e.mock.On("RefreshTokens", ctx, refreshToken)}
}

func (_c *AuthServiceMock_RefreshTokens_Call) Run(run func(ctx context.Context, refreshToken string)) *AuthServiceMock_RefreshTokens_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *AuthServiceMock_RefreshTokens_Call) Return(_a0 dto.JwtTokensOutput, _a1 error) *AuthServiceMock_RefreshTokens_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AuthServiceMock_RefreshTokens_Call) RunAndReturn(run func(context.Context, string) (dto.JwtTokensOutput, error)) *AuthServiceMock_RefreshTokens_Call {
	_c.Call.Return(run)
	return _c
}

// Register provides a mock function with given fields: ctx, input, localizer
func (_m *AuthServiceMock) Register(ctx context.Context, input dto.RegisterInput, localizer *i18n.Localizer) error {
	ret := _m.Called(ctx, input, localizer)

	if len(ret) == 0 {
		panic("no return value specified for Register")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, dto.RegisterInput, *i18n.Localizer) error); ok {
		r0 = rf(ctx, input, localizer)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AuthServiceMock_Register_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Register'
type AuthServiceMock_Register_Call struct {
	*mock.Call
}

// Register is a helper method to define mock.On call
//   - ctx context.Context
//   - input dto.RegisterInput
//   - localizer *i18n.Localizer
func (_e *AuthServiceMock_Expecter) Register(ctx interface{}, input interface{}, localizer interface{}) *AuthServiceMock_Register_Call {
	return &AuthServiceMock_Register_Call{Call: _e.mock.On("Register", ctx, input, localizer)}
}

func (_c *AuthServiceMock_Register_Call) Run(run func(ctx context.Context, input dto.RegisterInput, localizer *i18n.Localizer)) *AuthServiceMock_Register_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(dto.RegisterInput), args[2].(*i18n.Localizer))
	})
	return _c
}

func (_c *AuthServiceMock_Register_Call) Return(_a0 error) *AuthServiceMock_Register_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AuthServiceMock_Register_Call) RunAndReturn(run func(context.Context, dto.RegisterInput, *i18n.Localizer) error) *AuthServiceMock_Register_Call {
	_c.Call.Return(run)
	return _c
}

// RegisterConfirm provides a mock function with given fields: ctx, input
func (_m *AuthServiceMock) RegisterConfirm(ctx context.Context, input dto.RegisterConfirmInput) error {
	ret := _m.Called(ctx, input)

	if len(ret) == 0 {
		panic("no return value specified for RegisterConfirm")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, dto.RegisterConfirmInput) error); ok {
		r0 = rf(ctx, input)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AuthServiceMock_RegisterConfirm_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RegisterConfirm'
type AuthServiceMock_RegisterConfirm_Call struct {
	*mock.Call
}

// RegisterConfirm is a helper method to define mock.On call
//   - ctx context.Context
//   - input dto.RegisterConfirmInput
func (_e *AuthServiceMock_Expecter) RegisterConfirm(ctx interface{}, input interface{}) *AuthServiceMock_RegisterConfirm_Call {
	return &AuthServiceMock_RegisterConfirm_Call{Call: _e.mock.On("RegisterConfirm", ctx, input)}
}

func (_c *AuthServiceMock_RegisterConfirm_Call) Run(run func(ctx context.Context, input dto.RegisterConfirmInput)) *AuthServiceMock_RegisterConfirm_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(dto.RegisterConfirmInput))
	})
	return _c
}

func (_c *AuthServiceMock_RegisterConfirm_Call) Return(_a0 error) *AuthServiceMock_RegisterConfirm_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AuthServiceMock_RegisterConfirm_Call) RunAndReturn(run func(context.Context, dto.RegisterConfirmInput) error) *AuthServiceMock_RegisterConfirm_Call {
	_c.Call.Return(run)
	return _c
}

// RegisterOrLogin provides a mock function with given fields: ctx, userInfo
func (_m *AuthServiceMock) RegisterOrLogin(ctx context.Context, userInfo oauth.UserInfo) (dto.JwtTokensOutput, error) {
	ret := _m.Called(ctx, userInfo)

	if len(ret) == 0 {
		panic("no return value specified for RegisterOrLogin")
	}

	var r0 dto.JwtTokensOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, oauth.UserInfo) (dto.JwtTokensOutput, error)); ok {
		return rf(ctx, userInfo)
	}
	if rf, ok := ret.Get(0).(func(context.Context, oauth.UserInfo) dto.JwtTokensOutput); ok {
		r0 = rf(ctx, userInfo)
	} else {
		r0 = ret.Get(0).(dto.JwtTokensOutput)
	}

	if rf, ok := ret.Get(1).(func(context.Context, oauth.UserInfo) error); ok {
		r1 = rf(ctx, userInfo)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthServiceMock_RegisterOrLogin_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RegisterOrLogin'
type AuthServiceMock_RegisterOrLogin_Call struct {
	*mock.Call
}

// RegisterOrLogin is a helper method to define mock.On call
//   - ctx context.Context
//   - userInfo oauth.UserInfo
func (_e *AuthServiceMock_Expecter) RegisterOrLogin(ctx interface{}, userInfo interface{}) *AuthServiceMock_RegisterOrLogin_Call {
	return &AuthServiceMock_RegisterOrLogin_Call{Call: _e.mock.On("RegisterOrLogin", ctx, userInfo)}
}

func (_c *AuthServiceMock_RegisterOrLogin_Call) Run(run func(ctx context.Context, userInfo oauth.UserInfo)) *AuthServiceMock_RegisterOrLogin_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(oauth.UserInfo))
	})
	return _c
}

func (_c *AuthServiceMock_RegisterOrLogin_Call) Return(_a0 dto.JwtTokensOutput, _a1 error) *AuthServiceMock_RegisterOrLogin_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AuthServiceMock_RegisterOrLogin_Call) RunAndReturn(run func(context.Context, oauth.UserInfo) (dto.JwtTokensOutput, error)) *AuthServiceMock_RegisterOrLogin_Call {
	_c.Call.Return(run)
	return _c
}

// ResetPassword provides a mock function with given fields: ctx, input
func (_m *AuthServiceMock) ResetPassword(ctx context.Context, input dto.ResetPasswordInput) error {
	ret := _m.Called(ctx, input)

	if len(ret) == 0 {
		panic("no return value specified for ResetPassword")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, dto.ResetPasswordInput) error); ok {
		r0 = rf(ctx, input)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AuthServiceMock_ResetPassword_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ResetPassword'
type AuthServiceMock_ResetPassword_Call struct {
	*mock.Call
}

// ResetPassword is a helper method to define mock.On call
//   - ctx context.Context
//   - input dto.ResetPasswordInput
func (_e *AuthServiceMock_Expecter) ResetPassword(ctx interface{}, input interface{}) *AuthServiceMock_ResetPassword_Call {
	return &AuthServiceMock_ResetPassword_Call{Call: _e.mock.On("ResetPassword", ctx, input)}
}

func (_c *AuthServiceMock_ResetPassword_Call) Run(run func(ctx context.Context, input dto.ResetPasswordInput)) *AuthServiceMock_ResetPassword_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(dto.ResetPasswordInput))
	})
	return _c
}

func (_c *AuthServiceMock_ResetPassword_Call) Return(_a0 error) *AuthServiceMock_ResetPassword_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AuthServiceMock_ResetPassword_Call) RunAndReturn(run func(context.Context, dto.ResetPasswordInput) error) *AuthServiceMock_ResetPassword_Call {
	_c.Call.Return(run)
	return _c
}

// VerifyRecoveryCode provides a mock function with given fields: ctx, input
func (_m *AuthServiceMock) VerifyRecoveryCode(ctx context.Context, input dto.VerifyRecoveryCodeInput) error {
	ret := _m.Called(ctx, input)

	if len(ret) == 0 {
		panic("no return value specified for VerifyRecoveryCode")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, dto.VerifyRecoveryCodeInput) error); ok {
		r0 = rf(ctx, input)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AuthServiceMock_VerifyRecoveryCode_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'VerifyRecoveryCode'
type AuthServiceMock_VerifyRecoveryCode_Call struct {
	*mock.Call
}

// VerifyRecoveryCode is a helper method to define mock.On call
//   - ctx context.Context
//   - input dto.VerifyRecoveryCodeInput
func (_e *AuthServiceMock_Expecter) VerifyRecoveryCode(ctx interface{}, input interface{}) *AuthServiceMock_VerifyRecoveryCode_Call {
	return &AuthServiceMock_VerifyRecoveryCode_Call{Call: _e.mock.On("VerifyRecoveryCode", ctx, input)}
}

func (_c *AuthServiceMock_VerifyRecoveryCode_Call) Run(run func(ctx context.Context, input dto.VerifyRecoveryCodeInput)) *AuthServiceMock_VerifyRecoveryCode_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(dto.VerifyRecoveryCodeInput))
	})
	return _c
}

func (_c *AuthServiceMock_VerifyRecoveryCode_Call) Return(_a0 error) *AuthServiceMock_VerifyRecoveryCode_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AuthServiceMock_VerifyRecoveryCode_Call) RunAndReturn(run func(context.Context, dto.VerifyRecoveryCodeInput) error) *AuthServiceMock_VerifyRecoveryCode_Call {
	_c.Call.Return(run)
	return _c
}

// NewAuthServiceMock creates a new instance of AuthServiceMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAuthServiceMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *AuthServiceMock {
	mock := &AuthServiceMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
