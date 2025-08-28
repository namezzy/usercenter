import { createSlice, PayloadAction, createAsyncThunk } from '@reduxjs/toolkit';
import { User, LoginRequest, LoginResponse } from '@/types';
import { authApi } from '@/services/auth';

interface AuthState {
  isAuthenticated: boolean;
  token: string | null;
  user: User | null;
  loading: boolean;
  error: string | null;
}

const initialState: AuthState = {
  isAuthenticated: false,
  token: null,
  user: null,
  loading: false,
  error: null,
};

// 异步登录action
export const loginAsync = createAsyncThunk(
  'auth/login',
  async (loginData: LoginRequest, { rejectWithValue }) => {
    try {
      const response = await authApi.login(loginData);
      
      // 保存token到localStorage
      localStorage.setItem('token', response.token);
      localStorage.setItem('user', JSON.stringify(response.user));
      
      return response;
    } catch (error: any) {
      return rejectWithValue(error.message || '登录失败');
    }
  }
);

// 异步登出action
export const logoutAsync = createAsyncThunk(
  'auth/logout',
  async (_, { rejectWithValue }) => {
    try {
      await authApi.logout();
      
      // 清除本地存储
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      
      return;
    } catch (error: any) {
      // 即使API调用失败，也要清除本地存储
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      return rejectWithValue(error.message || '登出失败');
    }
  }
);

// 异步获取用户信息
export const getUserInfoAsync = createAsyncThunk(
  'auth/getUserInfo',
  async (_, { rejectWithValue }) => {
    try {
      const user = await authApi.getUserInfo();
      localStorage.setItem('user', JSON.stringify(user));
      return user;
    } catch (error: any) {
      return rejectWithValue(error.message || '获取用户信息失败');
    }
  }
);

// 刷新Token
export const refreshTokenAsync = createAsyncThunk(
  'auth/refreshToken',
  async (_, { rejectWithValue }) => {
    try {
      const response = await authApi.refreshToken();
      localStorage.setItem('token', response.token);
      return response.token;
    } catch (error: any) {
      // Token刷新失败，清除登录状态
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      return rejectWithValue(error.message || 'Token刷新失败');
    }
  }
);

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    // 清除错误
    clearError: (state) => {
      state.error = null;
    },
    
    // 从本地存储恢复登录状态
    restoreAuth: (state) => {
      const token = localStorage.getItem('token');
      const userStr = localStorage.getItem('user');
      
      if (token && userStr) {
        try {
          const user = JSON.parse(userStr);
          state.isAuthenticated = true;
          state.token = token;
          state.user = user;
        } catch (error) {
          // 解析失败，清除本地存储
          localStorage.removeItem('token');
          localStorage.removeItem('user');
        }
      }
    },
    
    // 更新用户信息
    updateUser: (state, action: PayloadAction<Partial<User>>) => {
      if (state.user) {
        state.user = { ...state.user, ...action.payload };
        localStorage.setItem('user', JSON.stringify(state.user));
      }
    },
  },
  extraReducers: (builder) => {
    // 登录
    builder
      .addCase(loginAsync.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(loginAsync.fulfilled, (state, action) => {
        state.loading = false;
        state.isAuthenticated = true;
        state.token = action.payload.token;
        state.user = action.payload.user;
        state.error = null;
      })
      .addCase(loginAsync.rejected, (state, action) => {
        state.loading = false;
        state.isAuthenticated = false;
        state.token = null;
        state.user = null;
        state.error = action.payload as string;
      });

    // 登出
    builder
      .addCase(logoutAsync.pending, (state) => {
        state.loading = true;
      })
      .addCase(logoutAsync.fulfilled, (state) => {
        state.loading = false;
        state.isAuthenticated = false;
        state.token = null;
        state.user = null;
        state.error = null;
      })
      .addCase(logoutAsync.rejected, (state, action) => {
        state.loading = false;
        // 即使登出失败，也要清除登录状态
        state.isAuthenticated = false;
        state.token = null;
        state.user = null;
        state.error = action.payload as string;
      });

    // 获取用户信息
    builder
      .addCase(getUserInfoAsync.pending, (state) => {
        state.loading = true;
      })
      .addCase(getUserInfoAsync.fulfilled, (state, action) => {
        state.loading = false;
        state.user = action.payload;
      })
      .addCase(getUserInfoAsync.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      });

    // 刷新Token
    builder
      .addCase(refreshTokenAsync.fulfilled, (state, action) => {
        state.token = action.payload;
      })
      .addCase(refreshTokenAsync.rejected, (state) => {
        state.isAuthenticated = false;
        state.token = null;
        state.user = null;
      });
  },
});

export const { clearError, restoreAuth, updateUser } = authSlice.actions;
export default authSlice.reducer;
