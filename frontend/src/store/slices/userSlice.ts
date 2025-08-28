import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface UserState {
  searchHistory: string[];
  preferences: {
    pageSize: number;
    theme: 'light' | 'dark';
    language: 'zh-CN' | 'en-US';
  };
}

const initialState: UserState = {
  searchHistory: [],
  preferences: {
    pageSize: 10,
    theme: 'light',
    language: 'zh-CN',
  },
};

const userSlice = createSlice({
  name: 'user',
  initialState,
  reducers: {
    // 添加搜索历史
    addSearchHistory: (state, action: PayloadAction<string>) => {
      const keyword = action.payload.trim();
      if (keyword && !state.searchHistory.includes(keyword)) {
        state.searchHistory.unshift(keyword);
        // 最多保存10条搜索历史
        if (state.searchHistory.length > 10) {
          state.searchHistory = state.searchHistory.slice(0, 10);
        }
      }
    },

    // 清除搜索历史
    clearSearchHistory: (state) => {
      state.searchHistory = [];
    },

    // 更新用户偏好设置
    updatePreferences: (state, action: PayloadAction<Partial<UserState['preferences']>>) => {
      state.preferences = { ...state.preferences, ...action.payload };
    },

    // 设置页面大小
    setPageSize: (state, action: PayloadAction<number>) => {
      state.preferences.pageSize = action.payload;
    },

    // 设置主题
    setTheme: (state, action: PayloadAction<'light' | 'dark'>) => {
      state.preferences.theme = action.payload;
    },

    // 设置语言
    setLanguage: (state, action: PayloadAction<'zh-CN' | 'en-US'>) => {
      state.preferences.language = action.payload;
    },
  },
});

export const {
  addSearchHistory,
  clearSearchHistory,
  updatePreferences,
  setPageSize,
  setTheme,
  setLanguage,
} = userSlice.actions;

export default userSlice.reducer;
