import React, { useEffect } from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { useDispatch } from 'react-redux';
import { AppDispatch } from './store';
import { restoreAuth } from './store/slices/authSlice';

// 组件导入
import Login from './pages/Login';
import Register from './pages/Register';
import Dashboard from './pages/Dashboard';
import Profile from './pages/Profile';
import AdminUsers from './pages/admin/Users';
import PrivateRoute from './components/PrivateRoute';
import MainLayout from './components/MainLayout';

const App: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>();

  useEffect(() => {
    // 应用启动时恢复登录状态
    dispatch(restoreAuth());
  }, [dispatch]);

  return (
    <Routes>
      {/* 公开路由 */}
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />
      
      {/* 受保护的路由 */}
      <Route path="/*" element={
        <PrivateRoute>
          <MainLayout>
            <Routes>
              <Route path="/" element={<Navigate to="/dashboard" replace />} />
              <Route path="/dashboard" element={<Dashboard />} />
              <Route path="/profile/*" element={<Profile />} />
              
              {/* 管理员路由 */}
              <Route path="/admin/users" element={
                <PrivateRoute requiredRoles={['admin', 'super_admin']}>
                  <AdminUsers />
                </PrivateRoute>
              } />
              
              {/* 404页面 */}
              <Route path="*" element={<div>页面不存在</div>} />
            </Routes>
          </MainLayout>
        </PrivateRoute>
      } />
    </Routes>
  );
};

export default App;
