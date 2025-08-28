import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse, AxiosError } from 'axios';
import { message } from 'antd';
import { ApiResponse } from '@/types';

// 创建axios实例
const request: AxiosInstance = axios.create({
  baseURL: process.env.REACT_APP_API_BASE_URL || '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器
request.interceptors.request.use(
  (config: AxiosRequestConfig) => {
    // 添加认证token
    const token = localStorage.getItem('token');
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    
    // 添加设备ID
    const deviceId = getDeviceId();
    if (deviceId && config.headers) {
      config.headers['X-Device-ID'] = deviceId;
    }
    
    return config;
  },
  (error: AxiosError) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    const { code, message: msg, data } = response.data;
    
    if (code === 200) {
      return { ...response, data };
    } else {
      message.error(msg || '请求失败');
      return Promise.reject(new Error(msg || '请求失败'));
    }
  },
  (error: AxiosError<ApiResponse>) => {
    const { response } = error;
    
    if (response) {
      const { status, data } = response;
      
      switch (status) {
        case 401:
          message.error('登录已过期，请重新登录');
          // 清除本地存储的用户信息
          localStorage.removeItem('token');
          localStorage.removeItem('user');
          // 跳转到登录页
          window.location.href = '/login';
          break;
        case 403:
          message.error('权限不足');
          break;
        case 404:
          message.error('请求的资源不存在');
          break;
        case 429:
          message.error('请求过于频繁，请稍后再试');
          break;
        case 500:
          message.error('服务器内部错误');
          break;
        default:
          message.error(data?.message || '请求失败');
      }
    } else if (error.request) {
      message.error('网络请求失败，请检查网络连接');
    } else {
      message.error('请求配置错误');
    }
    
    return Promise.reject(error);
  }
);

// 获取或生成设备ID
function getDeviceId(): string {
  let deviceId = localStorage.getItem('device_id');
  if (!deviceId) {
    deviceId = generateDeviceId();
    localStorage.setItem('device_id', deviceId);
  }
  return deviceId;
}

// 生成设备ID
function generateDeviceId(): string {
  return 'web_' + Math.random().toString(36).substr(2, 9) + '_' + Date.now();
}

// 通用GET请求
export const get = <T = any>(url: string, params?: any): Promise<T> => {
  return request.get(url, { params }).then(res => res.data);
};

// 通用POST请求
export const post = <T = any>(url: string, data?: any): Promise<T> => {
  return request.post(url, data).then(res => res.data);
};

// 通用PUT请求
export const put = <T = any>(url: string, data?: any): Promise<T> => {
  return request.put(url, data).then(res => res.data);
};

// 通用DELETE请求
export const del = <T = any>(url: string, params?: any): Promise<T> => {
  return request.delete(url, { params }).then(res => res.data);
};

// 上传文件
export const upload = <T = any>(url: string, formData: FormData): Promise<T> => {
  return request.post(url, formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  }).then(res => res.data);
};

export default request;
