import React from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { RootState } from '@/store';

interface PrivateRouteProps {
  children: React.ReactNode;
  requiredRoles?: string[];
  requiredPermissions?: string[];
}

const PrivateRoute: React.FC<PrivateRouteProps> = ({ 
  children, 
  requiredRoles = [],
  requiredPermissions = []
}) => {
  const location = useLocation();
  const { isAuthenticated, user } = useSelector((state: RootState) => state.auth);

  // 检查是否已登录
  if (!isAuthenticated || !user) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  // 检查角色权限
  if (requiredRoles.length > 0) {
    const userRoles = user.roles?.map(role => role.code) || [];
    const hasRequiredRole = requiredRoles.some(role => userRoles.includes(role));
    
    if (!hasRequiredRole) {
      return <Navigate to="/403" replace />;
    }
  }

  // 检查具体权限（这里简化处理，实际应该从后端获取用户权限）
  if (requiredPermissions.length > 0) {
    // TODO: 实现权限检查逻辑
  }

  return <>{children}</>;
};

export default PrivateRoute;
