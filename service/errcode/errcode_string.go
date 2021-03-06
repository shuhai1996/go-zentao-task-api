// Code generated by "stringer -type=ErrCode --linecomment"; DO NOT EDIT.

package errcode

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ErrWechatNetworkBusy-40000]
	_ = x[ErrWechatAuthTokenFailed-40001]
	_ = x[ErrWechatCacheTokenFailed-40002]
	_ = x[ErrWechatTokenExpired-40003]
	_ = x[ErrWechatUserInfoFailed-40004]
	_ = x[ErrWechatCacheUserInfoFailed-40005]
	_ = x[ErrWechatMobileFailed-40006]
	_ = x[ErrWechatInvalidParam-40100]
	_ = x[ErrWechatInvalidSign-40101]
	_ = x[ErrAdminNetworkBusy-39000]
	_ = x[ErrAdminLoginExpired-39001]
	_ = x[ErrAdminResourceForbidden-39002]
	_ = x[ErrAdminLoginTooFrequently-39003]
	_ = x[ErrAdminLoginAliyunAFSFailed-39004]
	_ = x[ErrAdminInvalidParam-39100]
	_ = x[ErrAdminInvalidSign-39101]
}

const (
	_ErrCode_name_0 = "网络繁忙，请稍后重试登录已过期，请重新登录您没有权限操作此资源频繁登录，请5分钟后重试人机验证失败"
	_ErrCode_name_1 = "无效的参数无效的签名"
	_ErrCode_name_2 = "网络繁忙，请稍后重试获取token失败，请重试获取token失败，请重试token已过期，请重试获取用户信息失败，请重试获取用户信息失败，请重试获取用户手机号失败，请重试"
	_ErrCode_name_3 = "无效的参数无效的签名"
)

var (
	_ErrCode_index_0 = [...]uint8{0, 30, 63, 93, 127, 145}
	_ErrCode_index_1 = [...]uint8{0, 15, 30}
	_ErrCode_index_2 = [...]uint8{0, 30, 59, 88, 114, 150, 186, 225}
	_ErrCode_index_3 = [...]uint8{0, 15, 30}
)

func (i ErrCode) String() string {
	switch {
	case 39000 <= i && i <= 39004:
		i -= 39000
		return _ErrCode_name_0[_ErrCode_index_0[i]:_ErrCode_index_0[i+1]]
	case 39100 <= i && i <= 39101:
		i -= 39100
		return _ErrCode_name_1[_ErrCode_index_1[i]:_ErrCode_index_1[i+1]]
	case 40000 <= i && i <= 40006:
		i -= 40000
		return _ErrCode_name_2[_ErrCode_index_2[i]:_ErrCode_index_2[i+1]]
	case 40100 <= i && i <= 40101:
		i -= 40100
		return _ErrCode_name_3[_ErrCode_index_3[i]:_ErrCode_index_3[i+1]]
	default:
		return "ErrCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
