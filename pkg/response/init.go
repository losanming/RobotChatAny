package response

var msg = make(map[int]string)

type StatusCode int

const (
	// 成功状态码 (2xx)
	StatusOK        StatusCode = 200 // 请求成功
	StatusCreated   StatusCode = 201 // 创建成功
	StatusAccepted  StatusCode = 202 // 请求已接受
	StatusNoContent StatusCode = 204 // 请求成功但无返回内容

	// 客户端错误码 (4xx)
	ErrBadRequest       StatusCode = 400 // 请求参数错误
	ErrUnauthorized     StatusCode = 401 // 未授权
	ErrForbidden        StatusCode = 403 // 禁止访问
	ErrNotFound         StatusCode = 404 // 资源不存在
	ErrMethodNotAllowed StatusCode = 405 // 方法不允许
	ErrConflict         StatusCode = 409 // 资源冲突
	ErrTooManyRequests  StatusCode = 429 // 请求过于频繁

	// 服务器错误码 (5xx)
	ErrInternalServer     StatusCode = 500 // 服务器内部错误
	ErrServiceUnavailable StatusCode = 503 // 服务不可用
	ErrGatewayTimeout     StatusCode = 504 // 网关超时

	// 数据库操作错误码 (1xxx)
	ErrDBConnection  StatusCode = 1001 // 数据库连接失败
	ErrDBQuery       StatusCode = 1002 // 数据库查询失败
	ErrDBInsert      StatusCode = 1003 // 数据库插入失败
	ErrDBUpdate      StatusCode = 1004 // 数据库更新失败
	ErrDBDelete      StatusCode = 1005 // 数据库删除失败
	ErrDBDuplicate   StatusCode = 1006 // 数据重复
	ErrDBNotFound    StatusCode = 1007 // 数据不存在
	ErrDBTransaction StatusCode = 1008 // 事务操作失败

	// 用户相关错误码 (2xxx)
	ErrUserNotFound    StatusCode = 2001 // 用户不存在
	ErrUserExists      StatusCode = 2002 // 用户已存在
	ErrPasswordWrong   StatusCode = 2003 // 密码错误
	ErrTokenInvalid    StatusCode = 2004 // Token无效
	ErrTokenExpired    StatusCode = 2005 // Token过期
	ErrUserDisabled    StatusCode = 2006 // 用户已禁用
	ErrUserNotVerified StatusCode = 2007 // 用户未验证

	// 文件操作错误码 (3xxx)
	ErrFileNotFound StatusCode = 3001 // 文件不存在
	ErrFileUpload   StatusCode = 3002 // 文件上传失败
	ErrFileDownload StatusCode = 3003 // 文件下载失败
	ErrFileDelete   StatusCode = 3004 // 文件删除失败
	ErrFileSize     StatusCode = 3005 // 文件大小超限
	ErrFileFormat   StatusCode = 3006 // 文件格式不支持

	// 缓存操作错误码 (4xxx)
	ErrCacheSet     StatusCode = 4001 // 缓存设置失败
	ErrCacheGet     StatusCode = 4002 // 缓存获取失败
	ErrCacheDelete  StatusCode = 4003 // 缓存删除失败
	ErrCacheExpired StatusCode = 4004 // 缓存已过期

	// 第三方服务错误码 (5xxx)
	ErrThirdPartyService StatusCode = 5001 // 第三方服务调用失败
	ErrSMSService        StatusCode = 5002 // 短信服务错误
	ErrEmailService      StatusCode = 5003 // 邮件服务错误
	ErrPaymentService    StatusCode = 5004 // 支付服务错误
	ErrAIService         StatusCode = 5005 // AI服务错误

	// 业务逻辑错误码 (6xxx)
	ErrOrderCreate         StatusCode = 6001 // 订单创建失败
	ErrOrderStatus         StatusCode = 6002 // 订单状态错误
	ErrOrderPayment        StatusCode = 6003 // 订单支付失败
	ErrOrderRefund         StatusCode = 6004 // 订单退款失败
	ErrBalanceInsufficient StatusCode = 6005 // 余额不足

	//认证方面错误 (7xxx)
	ErrLoginFailed           StatusCode = 7001 // 登录失败
	ErrCaptchaGenerateFailed StatusCode = 7002 // 验证码生成失败
	//验证码验证失败
	ErrCaptchaVerifyFailed StatusCode = 7003
)

// StatusText 返回状态码对应的文本描述
func (c StatusCode) StatusText() string {
	switch c {
	// 成功状态码
	case StatusOK:
		return "请求成功"
	case StatusCreated:
		return "创建成功"
	case StatusAccepted:
		return "请求已接受"
	case StatusNoContent:
		return "请求成功但无返回内容"

	// 客户端错误码
	case ErrBadRequest:
		return "请求参数错误"
	case ErrUnauthorized:
		return "未授权"
	case ErrForbidden:
		return "禁止访问"
	case ErrNotFound:
		return "资源不存在"
	case ErrMethodNotAllowed:
		return "方法不允许"
	case ErrConflict:
		return "资源冲突"
	case ErrTooManyRequests:
		return "请求过于频繁"

	// 服务器错误码
	case ErrInternalServer:
		return "服务器内部错误"
	case ErrServiceUnavailable:
		return "服务不可用"
	case ErrGatewayTimeout:
		return "网关超时"

	// 数据库操作错误码
	case ErrDBConnection:
		return "数据库连接失败"
	case ErrDBQuery:
		return "数据库查询失败"
	case ErrDBInsert:
		return "数据库插入失败"
	case ErrDBUpdate:
		return "数据库更新失败"
	case ErrDBDelete:
		return "数据库删除失败"
	case ErrDBDuplicate:
		return "数据重复"
	case ErrDBNotFound:
		return "数据不存在"
	case ErrDBTransaction:
		return "事务操作失败"

	// 用户相关错误码
	case ErrUserNotFound:
		return "用户不存在"
	case ErrUserExists:
		return "用户已存在"
	case ErrPasswordWrong:
		return "密码错误"
	case ErrTokenInvalid:
		return "Token无效"
	case ErrTokenExpired:
		return "Token过期"
	case ErrUserDisabled:
		return "用户已禁用"
	case ErrUserNotVerified:
		return "用户未验证"

	// 文件操作错误码
	case ErrFileNotFound:
		return "文件不存在"
	case ErrFileUpload:
		return "文件上传失败"
	case ErrFileDownload:
		return "文件下载失败"
	case ErrFileDelete:
		return "文件删除失败"
	case ErrFileSize:
		return "文件大小超限"
	case ErrFileFormat:
		return "文件格式不支持"

	// 缓存操作错误码
	case ErrCacheSet:
		return "缓存设置失败"
	case ErrCacheGet:
		return "缓存获取失败"
	case ErrCacheDelete:
		return "缓存删除失败"
	case ErrCacheExpired:
		return "缓存已过期"

	// 第三方服务错误码
	case ErrThirdPartyService:
		return "第三方服务调用失败"
	case ErrSMSService:
		return "短信服务错误"
	case ErrEmailService:
		return "邮件服务错误"
	case ErrPaymentService:
		return "支付服务错误"
	case ErrAIService:
		return "AI服务错误"

	// 业务逻辑错误码
	case ErrOrderCreate:
		return "订单创建失败"
	case ErrOrderStatus:
		return "订单状态错误"
	case ErrOrderPayment:
		return "订单支付失败"
	case ErrOrderRefund:
		return "订单退款失败"
	case ErrBalanceInsufficient:
		return "余额不足"

	// 认证方面错误
	case ErrLoginFailed:
		return "登录失败"
	case ErrCaptchaGenerateFailed:
		return "验证码生成失败"
	case ErrCaptchaVerifyFailed:
		return "验证码验证失败"
	default:
		return "未知错误"
	}
}

// RegisterMsg 注册自定义消息
func RegisterMsg(code int, msg1 string) {
	msg[code] = msg1
}

// RegisterMsgMap 注册自定义消息, 多个
func RegisterMsgMap(msg1 map[int]string) {
	for code, msg2 := range msg1 {
		msg[code] = msg2
	}
}
