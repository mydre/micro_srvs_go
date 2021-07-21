package forms

type SendSmsForm struct{
	//手机号码格式有规范可寻，自定义validator
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"`
	//希望type在指定的范围内进行选择
	Type uint `form:"type" json:"type" binding:"required,oneof=1 2"`//1代表注册，2代表动态验证码登录
	//1. 注册发送短息验证码和动态验证码登录发送验证码
}