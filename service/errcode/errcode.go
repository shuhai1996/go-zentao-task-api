//go:generate stringer -type=ErrCode --linecomment

package errcode

type ErrCode int

const (
	// 微信小程序通用错误码
	ErrWechatNetworkBusy         ErrCode = 40000 // 网络繁忙，请稍后重试
	ErrWechatAuthTokenFailed     ErrCode = 40001 // 获取token失败，请重试
	ErrWechatCacheTokenFailed    ErrCode = 40002 // 获取token失败，请重试
	ErrWechatTokenExpired        ErrCode = 40003 // token已过期，请重试
	ErrWechatUserInfoFailed      ErrCode = 40004 // 获取用户信息失败，请重试
	ErrWechatCacheUserInfoFailed ErrCode = 40005 // 获取用户信息失败，请重试
	ErrWechatMobileFailed        ErrCode = 40006 // 获取用户手机号失败，请重试
	ErrWechatInvalidParam        ErrCode = 40100 // 无效的参数
	ErrWechatInvalidSign         ErrCode = 40101 // 无效的签名

	// 后台通用错误码
	ErrAdminNetworkBusy          ErrCode = 39000 // 网络繁忙，请稍后重试
	ErrAdminLoginExpired         ErrCode = 39001 // 登录已过期，请重新登录
	ErrAdminResourceForbidden    ErrCode = 39002 // 您没有权限操作此资源
	ErrAdminLoginTooFrequently   ErrCode = 39003 // 频繁登录，请5分钟后重试
	ErrAdminLoginAliyunAFSFailed ErrCode = 39004 // 人机验证失败
	ErrAdminInvalidParam         ErrCode = 39100 // 无效的参数
	ErrAdminInvalidSign          ErrCode = 39101 // 无效的签名
)
