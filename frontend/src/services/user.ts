import { 
  User,
  UpdateProfileRequest,
  ChangePasswordRequest,
  BindEmailRequest,
  BindPhoneRequest,
  UserDevice,
  UserLog,
  PageResponse
} from '@/types';
import { get, post, put, del, upload } from './request';

// 用户相关API
export const userApi = {
  // 获取个人资料
  getProfile: (): Promise<User> => {
    return get('/profile');
  },

  // 更新个人资料
  updateProfile: (data: UpdateProfileRequest): Promise<void> => {
    return put('/profile', data);
  },

  // 修改密码
  changePassword: (data: ChangePasswordRequest): Promise<void> => {
    return put('/profile/password', data);
  },

  // 上传头像
  uploadAvatar: (file: File): Promise<{ avatar_url: string }> => {
    const formData = new FormData();
    formData.append('avatar', file);
    return upload('/profile/avatar', formData);
  },

  // 绑定邮箱
  bindEmail: (data: BindEmailRequest): Promise<void> => {
    return post('/profile/bind-email', data);
  },

  // 绑定手机号
  bindPhone: (data: BindPhoneRequest): Promise<void> => {
    return post('/profile/bind-phone', data);
  },

  // 获取设备列表
  getDevices: (): Promise<UserDevice[]> => {
    return get('/profile/devices');
  },

  // 移除设备
  removeDevice: (deviceId: string): Promise<void> => {
    return del(`/profile/devices/${deviceId}`);
  },

  // 获取操作日志
  getLogs: (page: number = 1, pageSize: number = 10): Promise<PageResponse<UserLog>> => {
    return get('/profile/logs', { page, page_size: pageSize });
  },

  // 获取登录日志
  getLoginLogs: (params: { page: number; pageSize: number }): Promise<{ data: any[]; total: number }> => {
    return get('/profile/login-logs', params);
  },

  // 获取安全设置
  getSecuritySettings: (): Promise<any> => {
    return get('/profile/security-settings');
  },

  // 更新安全设置
  updateSecuritySettings: (data: any): Promise<void> => {
    return put('/profile/security-settings', data);
  },

  // 获取仪表盘统计数据
  getDashboardStats: (): Promise<any> => {
    return get('/dashboard/stats');
  },

  // 获取最近活动
  getRecentActivities: (): Promise<any[]> => {
    return get('/dashboard/activities');
  },

  // 获取用户增长数据
  getUserGrowthData: (params: { days: number }): Promise<any[]> => {
    return get('/dashboard/growth', params);
  },
};

// 便捷导出方法
export const {
  getProfile,
  updateProfile,
  changePassword,
  uploadAvatar,
  bindEmail,
  bindPhone,
  getDevices,
  removeDevice,
  getLogs,
  getLoginLogs,
  getSecuritySettings,
  updateSecuritySettings,
  getDashboardStats,
  getRecentActivities,
  getUserGrowthData
} = userApi;
