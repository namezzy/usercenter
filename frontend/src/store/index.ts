import { configureStore } from '@reduxjs/toolkit';
import { persistStore, persistReducer } from 'redux-persist';
import storage from 'redux-persist/lib/storage';
import { combineReducers } from '@reduxjs/toolkit';
import { useDispatch, useSelector, TypedUseSelectorHook } from 'react-redux';

import authSlice from './slices/authSlice';
import userSlice from './slices/userSlice';
import themeSlice from './slices/themeSlice';

// 持久化配置
const persistConfig = {
  key: 'root',
  storage,
  whitelist: ['auth', 'theme'], // 只持久化这些reducer
};

// 根reducer
const rootReducer = combineReducers({
  auth: authSlice,
  user: userSlice,
  theme: themeSlice,
});

// 持久化reducer
const persistedReducer = persistReducer(persistConfig, rootReducer);

// 创建store
export const store = configureStore({
  reducer: persistedReducer,
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        ignoredActions: ['persist/PERSIST', 'persist/REHYDRATE'],
      },
    }),
});

// 创建persistor
export const persistor = persistStore(store);

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

// 类型化的hooks
export const useAppDispatch = () => useDispatch<AppDispatch>();
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector;
