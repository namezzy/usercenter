import { 
  User,
  UserListQuery,
  UpdateProfileRequest,
  RegisterRequest,
  PageResponse
} from '@/types';
import { get, post, put, del, upload } from './request';

// 管理员相关API
export const adminApi = {
  // 获取用户列表
  getUsers: (params: any): Promise<{ data: User[]; total: number }> => {
    return get('/admin/users', params);
  },

  // 创建用户
  createUser: (data: any): Promise<void> => {
    return post('/admin/users', data);
  },

  // 获取用户详情
  getUser: (id: string): Promise<User> => {
    return get(`/admin/users/${id}`);
  },

  // 更新用户
  updateUser: (id: string, data: any): Promise<void> => {
    return put(`/admin/users/${id}`, data);
  },

  // 删除用户
  deleteUser: (id: string): Promise<void> => {
    return del(`/admin/users/${id}`);
  },

  // 导出用户
  exportUsers: (params: any): Promise<any> => {
    return get('/admin/users/export', params);
  },

  // 导入用户
  importUsers: (formData: FormData): Promise<void> => {
    return upload('/admin/users/import', formData);
  },

  // 锁定用户
  lockUser: (id: string): Promise<void> => {
    return post(`/admin/users/${id}/lock`);
  },

  // 解锁用户
  unlockUser: (id: string): Promise<void> => {
    return post(`/admin/users/${id}/unlock`);
  },

  // 重置密码
  resetPassword: (id: string): Promise<void> => {
    return post(`/admin/users/${id}/reset-password`);
  },

  // 获取角色列表
  getRoles: (params?: any): Promise<{ data: any[]; total: number }> => {
    return get('/admin/roles', params);
  },

  // 获取权限列表
  getPermissions: (params?: any): Promise<{ data: any[]; total: number }> => {
    return get('/admin/permissions', params);
  },

  // 获取统计信息
  getStatistics: (): Promise<{
    total_users: number;
    today_new: number;
    today_active: number;
    online_users: number;
  }> => {
    return get('/admin/statistics');
  },
};

// 便捷导出方法
export const {
  getUsers,
  createUser,
  updateUser,
  deleteUser,
  exportUsers,
  importUsers,
  lockUser,
  unlockUser,
  resetPassword,
  getRoles,
  getPermissions
} = adminApi;
