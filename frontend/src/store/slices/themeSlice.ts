import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface ThemeState {
  mode: 'light' | 'dark';
  primaryColor: string;
  borderRadius: number;
  fontSize: 'small' | 'medium' | 'large';
  compactMode: boolean;
}

const initialState: ThemeState = {
  mode: 'light',
  primaryColor: '#1890ff',
  borderRadius: 6,
  fontSize: 'medium',
  compactMode: false,
};

const themeSlice = createSlice({
  name: 'theme',
  initialState,
  reducers: {
    // 切换主题模式
    toggleTheme: (state) => {
      state.mode = state.mode === 'light' ? 'dark' : 'light';
    },

    // 设置主题模式
    setThemeMode: (state, action: PayloadAction<'light' | 'dark'>) => {
      state.mode = action.payload;
    },

    // 设置主色调
    setPrimaryColor: (state, action: PayloadAction<string>) => {
      state.primaryColor = action.payload;
    },

    // 设置圆角大小
    setBorderRadius: (state, action: PayloadAction<number>) => {
      state.borderRadius = action.payload;
    },

    // 设置字体大小
    setFontSize: (state, action: PayloadAction<'small' | 'medium' | 'large'>) => {
      state.fontSize = action.payload;
    },

    // 切换紧凑模式
    toggleCompactMode: (state) => {
      state.compactMode = !state.compactMode;
    },

    // 设置紧凑模式
    setCompactMode: (state, action: PayloadAction<boolean>) => {
      state.compactMode = action.payload;
    },

    // 重置主题设置
    resetTheme: (state) => {
      return initialState;
    },
  },
});

export const {
  toggleTheme,
  setThemeMode,
  setPrimaryColor,
  setBorderRadius,
  setFontSize,
  toggleCompactMode,
  setCompactMode,
  resetTheme,
} = themeSlice.actions;

export default themeSlice.reducer;
