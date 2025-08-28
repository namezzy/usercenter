import { 
  LoginRequest, 
  LoginResponse, 
  RegisterRequest, 
  CaptchaResponse,
  User 
} from '@/types';
import { get, post } from './request';

// 认证相关API
export const authApi = {
  // 获取验证码
  getCaptcha: (): Promise<CaptchaResponse> => {
    return get('/auth/captcha');
  },

  // 发送邮箱验证码
  sendEmailCode: (email: string, purpose: string): Promise<void> => {
    return post('/auth/send-email-code', { email, purpose });
  },

  // 发送短信验证码
  sendSMSCode: (phone: string, purpose: string): Promise<void> => {
    return post('/auth/send-sms-code', { phone, purpose });
  },

  // 注册
  register: (data: {
    username: string;
    email?: string;
    phone?: string;
    password: string;
    verificationCode: string;
    registerType: 'email' | 'phone';
  }): Promise<any> => {
    return post('/auth/register', data);
  },

  // 发送验证码
  sendVerificationCode: (data: {
    type: 'email' | 'phone';
    target: string;
    purpose: 'register' | 'login' | 'reset_password';
  }): Promise<any> => {
    return post('/auth/send-verification-code', data);
  },

  // 用户登录
  login: (data: LoginRequest): Promise<LoginResponse> => {
    return post('/auth/login', data);
  },

  // 用户登出
  logout: (): Promise<void> => {
    return post('/auth/logout');
  },

  // 刷新Token
  refreshToken: (): Promise<{ token: string }> => {
    return post('/auth/refresh');
  },

  // 获取当前用户信息
  getUserInfo: (): Promise<User> => {
    return get('/auth/user');
  },
};

// 便捷导出方法
export const { 
  getCaptcha, 
  sendEmailCode, 
  sendSMSCode, 
  register, 
  sendVerificationCode,
  login, 
  logout, 
  refreshToken, 
  getUserInfo 
} = authApi;
