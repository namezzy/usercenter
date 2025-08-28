import React, { useState, useEffect } from 'react';
import { Layout, Menu, Avatar, Dropdown, Space, Badge, Button, Drawer } from 'antd';
import { 
  MenuFoldOutlined, 
  MenuUnfoldOutlined, 
  UserOutlined, 
  SettingOutlined, 
  LogoutOutlined,
  BellOutlined,
  MenuOutlined
} from '@ant-design/icons';
import { useNavigate, useLocation } from 'react-router-dom';
import { useSelector, useDispatch } from 'react-redux';
import { RootState, AppDispatch } from '@/store';
import { logoutAsync } from '@/store/slices/authSlice';
import { useMediaQuery } from '@/hooks/useMediaQuery';

const { Header, Sider, Content } = Layout;

interface MainLayoutProps {
  children: React.ReactNode;
}

const MainLayout: React.FC<MainLayoutProps> = ({ children }) => {
  const [collapsed, setCollapsed] = useState(false);
  const [mobileMenuVisible, setMobileMenuVisible] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();
  const dispatch = useDispatch<AppDispatch>();
  const isMobile = useMediaQuery('(max-width: 768px)');
  
  const { user } = useSelector((state: RootState) => state.auth);
  const { mode } = useSelector((state: RootState) => state.theme);

  // 菜单项配置
  const menuItems = [
    {
      key: '/dashboard',
      icon: <UserOutlined />,
      label: '仪表盘',
    },
    {
      key: '/profile',
      icon: <UserOutlined />,
      label: '个人中心',
      children: [
        {
          key: '/profile/info',
          label: '基本信息',
        },
        {
          key: '/profile/security',
          label: '安全设置',
        },
        {
          key: '/profile/devices',
          label: '设备管理',
        },
        {
          key: '/profile/logs',
          label: '操作日志',
        },
      ],
    },
  ];

  // 管理员菜单
  const adminMenuItems = [
    {
      key: '/admin',
      icon: <SettingOutlined />,
      label: '系统管理',
      children: [
        {
          key: '/admin/users',
          label: '用户管理',
        },
        {
          key: '/admin/roles',
          label: '角色管理',
        },
        {
          key: '/admin/permissions',
          label: '权限管理',
        },
        {
          key: '/admin/logs',
          label: '系统日志',
        },
      ],
    },
  ];

  // 根据用户角色显示不同菜单
  const getMenuItems = () => {
    let items = [...menuItems];
    
    if (user?.roles?.some(role => ['admin', 'super_admin'].includes(role.code))) {
      items = [...items, ...adminMenuItems];
    }
    
    return items;
  };

  // 用户下拉菜单
  const userMenuItems = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人中心',
      onClick: () => navigate('/profile'),
    },
    {
      key: 'settings',
      icon: <SettingOutlined />,
      label: '设置',
      onClick: () => navigate('/settings'),
    },
    {
      type: 'divider',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      onClick: handleLogout,
    },
  ];

  // 处理登出
  async function handleLogout() {
    try {
      await dispatch(logoutAsync()).unwrap();
      navigate('/login');
    } catch (error) {
      console.error('Logout failed:', error);
    }
  }

  // 处理菜单点击
  const handleMenuClick = ({ key }: { key: string }) => {
    navigate(key);
    if (isMobile) {
      setMobileMenuVisible(false);
    }
  };

  // 移动端菜单
  const MobileMenu = () => (
    <Drawer
      title="菜单"
      placement="left"
      onClose={() => setMobileMenuVisible(false)}
      open={mobileMenuVisible}
      bodyStyle={{ padding: 0 }}
    >
      <Menu
        theme={mode}
        mode="inline"
        selectedKeys={[location.pathname]}
        items={getMenuItems()}
        onClick={handleMenuClick}
      />
    </Drawer>
  );

  useEffect(() => {
    if (isMobile) {
      setCollapsed(true);
    }
  }, [isMobile]);

  return (
    <Layout style={{ minHeight: '100vh' }}>
      {/* 移动端菜单抽屉 */}
      {isMobile && <MobileMenu />}
      
      {/* 侧边栏 */}
      {!isMobile && (
        <Sider 
          trigger={null} 
          collapsible 
          collapsed={collapsed}
          theme={mode}
          width={256}
          collapsedWidth={80}
        >
          <div className="logo" style={{ 
            height: 64, 
            display: 'flex', 
            alignItems: 'center', 
            justifyContent: 'center',
            fontSize: collapsed ? 16 : 20,
            fontWeight: 'bold',
            color: mode === 'dark' ? '#fff' : '#1890ff'
          }}>
            {collapsed ? 'UC' : '用户中心'}
          </div>
          
          <Menu
            theme={mode}
            mode="inline"
            selectedKeys={[location.pathname]}
            items={getMenuItems()}
            onClick={handleMenuClick}
          />
        </Sider>
      )}

      <Layout>
        {/* 顶部导航 */}
        <Header style={{ 
          padding: '0 16px', 
          background: mode === 'dark' ? '#001529' : '#fff',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          borderBottom: `1px solid ${mode === 'dark' ? '#303030' : '#f0f0f0'}`
        }}>
          <div style={{ display: 'flex', alignItems: 'center' }}>
            {isMobile ? (
              <Button
                type="text"
                icon={<MenuOutlined />}
                onClick={() => setMobileMenuVisible(true)}
                style={{ fontSize: 16 }}
              />
            ) : (
              <Button
                type="text"
                icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
                onClick={() => setCollapsed(!collapsed)}
                style={{ fontSize: 16 }}
              />
            )}
          </div>

          <Space size="middle">
            {/* 通知铃铛 */}
            <Badge count={0} size="small">
              <Button 
                type="text" 
                icon={<BellOutlined />} 
                onClick={() => navigate('/notifications')}
              />
            </Badge>

            {/* 用户头像和信息 */}
            <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
              <Space style={{ cursor: 'pointer' }}>
                <Avatar 
                  src={user?.avatar} 
                  icon={<UserOutlined />}
                  size="small"
                />
                <span style={{ color: mode === 'dark' ? '#fff' : undefined }}>
                  {user?.nickname || user?.username}
                </span>
              </Space>
            </Dropdown>
          </Space>
        </Header>

        {/* 主内容区 */}
        <Content style={{ 
          margin: '16px',
          padding: '16px',
          background: mode === 'dark' ? '#141414' : '#fff',
          borderRadius: 8,
          minHeight: 'calc(100vh - 112px)',
          overflow: 'auto'
        }}>
          {children}
        </Content>
      </Layout>
    </Layout>
  );
};

export default MainLayout;
