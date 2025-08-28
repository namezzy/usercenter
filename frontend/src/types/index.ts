// 用户相关类型定义
export interface User {
  id: string;
  username: string;
  email?: string;
  phone?: string;
  nickname: string;
  avatar?: string;
  gender: number;
  birthday?: string;
  bio?: string;
  status: number;
  email_verified: boolean;
  phone_verified: boolean;
  last_login_at?: string;
  last_login_ip?: string;
  created_at: string;
  updated_at: string;
  roles: Role[];
}

// 角色类型
export interface Role {
  id: string;
  name: string;
  code: string;
  description: string;
  status: number;
  sort: number;
  created_at: string;
}

// 权限类型
export interface Permission {
  id: string;
  name: string;
  code: string;
  type: string;
  parent_id?: string;
  path?: string;
  method?: string;
  icon?: string;
  sort: number;
  status: number;
  children?: Permission[];
}

// 设备信息
export interface UserDevice {
  id: string;
  device_id: string;
  device_type: string;
  device_name: string;
  ip: string;
  user_agent: string;
  last_active: string;
  is_active: boolean;
  created_at: string;
}

// 用户日志
export interface UserLog {
  id: string;
  action: string;
  module: string;
  ip: string;
  user_agent: string;
  details: string;
  status: number;
  created_at: string;
}

// 系统通知
export interface SystemNotification {
  id: string;
  title: string;
  content: string;
  type: string;
  priority: number;
  status: number;
  send_to_all: boolean;
  send_at?: string;
  expire_at?: string;
  created_at: string;
}

// 用户通知
export interface UserNotification {
  id: string;
  notification_id: string;
  read_at?: string;
  is_read: boolean;
  notification: SystemNotification;
  created_at: string;
}

// 登录请求
export interface LoginRequest {
  username: string;
  password: string;
  captcha_id?: string;
  captcha_code?: string;
  device_info: {
    device_id?: string;
    device_type?: string;
    device_name?: string;
  };
}

// 登录响应
export interface LoginResponse {
  token: string;
  user: User;
  expires_at: number;
}

// 注册请求
export interface RegisterRequest {
  username: string;
  email?: string;
  phone?: string;
  password: string;
  nickname?: string;
  email_code?: string;
  sms_code?: string;
  captcha_id: string;
  captcha_code: string;
}

// 更新个人资料请求
export interface UpdateProfileRequest {
  nickname?: string;
  gender?: number;
  birthday?: string;
  bio?: string;
}

// 修改密码请求
export interface ChangePasswordRequest {
  old_password: string;
  new_password: string;
}

// 绑定邮箱请求
export interface BindEmailRequest {
  email: string;
  email_code: string;
}

// 绑定手机请求
export interface BindPhoneRequest {
  phone: string;
  sms_code: string;
}

// 验证码响应
export interface CaptchaResponse {
  captcha_id: string;
  captcha_img: string;
}

// 分页查询参数
export interface PageQuery {
  page: number;
  page_size: number;
  keyword?: string;
}

// 用户列表查询参数
export interface UserListQuery extends PageQuery {
  status?: number;
  role_code?: string;
}

// 分页响应
export interface PageResponse<T> {
  total: number;
  items: T[];
}

// API响应格式
export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data?: T;
}

// 性别枚举
export enum Gender {
  Unknown = 0,
  Male = 1,
  Female = 2,
}

// 用户状态枚举
export enum UserStatus {
  Normal = 1,
  Disabled = 2,
  Locked = 3,
}

// 通知类型枚举
export enum NotificationType {
  Info = 'info',
  Warning = 'warning',
  Error = 'error',
  Success = 'success',
}

// 权限类型枚举
export enum PermissionType {
  Menu = 'menu',
  Button = 'button',
  Data = 'data',
}

// 表单字段类型
export interface FormField {
  name: string;
  label: string;
  type: 'input' | 'password' | 'email' | 'select' | 'textarea' | 'date' | 'upload';
  rules?: any[];
  options?: { label: string; value: any }[];
  placeholder?: string;
}

// 菜单项类型
export interface MenuItem {
  key: string;
  label: string;
  icon?: React.ReactNode;
  children?: MenuItem[];
  path?: string;
  permission?: string;
}

// 路由配置类型
export interface RouteConfig {
  path: string;
  component: React.ComponentType<any>;
  exact?: boolean;
  meta?: {
    title?: string;
    requireAuth?: boolean;
    roles?: string[];
    permissions?: string[];
  };
  children?: RouteConfig[];
}
