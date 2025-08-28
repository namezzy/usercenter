import React, { useState, useEffect } from 'react';
import { Row, Col, Card, Statistic, Table, List, Avatar, Progress, Space, Tag, Typography, DatePicker, Select } from 'antd';
import { 
  UserOutlined, 
  TeamOutlined, 
  SecurityScanOutlined, 
  BellOutlined,
  TrophyOutlined,
  ClockCircleOutlined,
  EyeOutlined,
  MessageOutlined,
  StarOutlined,
  RiseOutlined
} from '@ant-design/icons';
import { Line, Column, Pie } from '@ant-design/plots';
import { useAppSelector } from '@/store';
import { getDashboardStats, getRecentActivities, getUserGrowthData } from '@/services/user';
import './Dashboard.css';

const { Title, Text } = Typography;
const { RangePicker } = DatePicker;
const { Option } = Select;

interface DashboardStats {
  totalUsers: number;
  activeUsers: number;
  todayLogin: number;
  securityAlerts: number;
  notifications: number;
  userGrowth: number;
  systemLoad: number;
  responseTime: number;
}

interface Activity {
  id: string;
  type: string;
  description: string;
  time: string;
  user: {
    name: string;
    avatar: string;
  };
  status: 'success' | 'warning' | 'error' | 'info';
}

interface GrowthData {
  date: string;
  users: number;
  active: number;
}

const Dashboard: React.FC = () => {
  const { user } = useAppSelector(state => state.auth);
  const [stats, setStats] = useState<DashboardStats>({
    totalUsers: 0,
    activeUsers: 0,
    todayLogin: 0,
    securityAlerts: 0,
    notifications: 0,
    userGrowth: 0,
    systemLoad: 0,
    responseTime: 0
  });
  const [activities, setActivities] = useState<Activity[]>([]);
  const [growthData, setGrowthData] = useState<GrowthData[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadDashboardData();
  }, []);

  const loadDashboardData = async () => {
    try {
      setLoading(true);
      const [statsData, activitiesData, growthDataResult] = await Promise.all([
        getDashboardStats(),
        getRecentActivities(),
        getUserGrowthData({ days: 30 })
      ]);
      
      setStats(statsData);
      setActivities(activitiesData);
      setGrowthData(growthDataResult);
    } catch (error) {
      console.error('加载仪表盘数据失败:', error);
    } finally {
      setLoading(false);
    }
  };

  // 用户增长图表配置
  const lineConfig = {
    data: growthData,
    xField: 'date',
    yField: 'users',
    seriesField: 'type',
    smooth: true,
    color: ['#1890ff', '#52c41a'],
    point: {
      size: 3,
      shape: 'circle',
    },
    tooltip: {
      showMarkers: false,
    },
    legend: {
      position: 'top' as const,
    },
  };

  // 用户活跃度饼图配置
  const pieData = [
    { type: '活跃用户', value: stats.activeUsers },
    { type: '非活跃用户', value: stats.totalUsers - stats.activeUsers },
  ];

  const pieConfig = {
    data: pieData,
    angleField: 'value',
    colorField: 'type',
    radius: 0.8,
    label: {
      type: 'outer',
      content: '{name} {percentage}',
    },
    interactions: [{ type: 'element-active' }],
    color: ['#52c41a', '#f5f5f5'],
  };

  // 系统状态柱状图配置
  const columnData = [
    { name: 'CPU', value: 65 },
    { name: '内存', value: 78 },
    { name: '磁盘', value: 45 },
    { name: '网络', value: 32 },
  ];

  const columnConfig = {
    data: columnData,
    xField: 'name',
    yField: 'value',
    color: '#1890ff',
    columnWidthRatio: 0.6,
    meta: {
      value: { max: 100 },
    },
  };

  // 活动状态颜色映射
  const getActivityColor = (status: string) => {
    const colors = {
      success: '#52c41a',
      warning: '#faad14',
      error: '#ff4d4f',
      info: '#1890ff',
    };
    return colors[status as keyof typeof colors] || '#666';
  };

  // 活动类型图标映射
  const getActivityIcon = (type: string) => {
    const icons = {
      login: <UserOutlined />,
      register: <TeamOutlined />,
      security: <SecurityScanOutlined />,
      notification: <BellOutlined />,
      achievement: <TrophyOutlined />,
    };
    return icons[type as keyof typeof icons] || <ClockCircleOutlined />;
  };

  const quickActions = [
    { title: '用户管理', icon: <UserOutlined />, color: '#1890ff', path: '/admin/users' },
    { title: '角色权限', icon: <SecurityScanOutlined />, color: '#52c41a', path: '/admin/roles' },
    { title: '系统设置', icon: <BellOutlined />, color: '#faad14', path: '/admin/settings' },
    { title: '日志查看', icon: <EyeOutlined />, color: '#722ed1', path: '/admin/logs' },
  ];

  return (
    <div className="dashboard-container">
      <div className="dashboard-header">
        <Title level={2}>仪表盘</Title>
        <Text type="secondary">欢迎回来，{user?.username}！</Text>
      </div>

      {/* 统计卡片 */}
      <Row gutter={[16, 16]} className="stats-row">
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="总用户数"
              value={stats.totalUsers}
              prefix={<UserOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="活跃用户"
              value={stats.activeUsers}
              prefix={<TeamOutlined />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="今日登录"
              value={stats.todayLogin}
              prefix={<ClockCircleOutlined />}
              valueStyle={{ color: '#faad14' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="安全警报"
              value={stats.securityAlerts}
              prefix={<SecurityScanOutlined />}
              valueStyle={{ color: '#ff4d4f' }}
            />
          </Card>
        </Col>
      </Row>

      {/* 图表区域 */}
      <Row gutter={[16, 16]} className="charts-row">
        <Col xs={24} lg={16}>
          <Card title="用户增长趋势" loading={loading}>
            <Line {...lineConfig} height={300} />
          </Card>
        </Col>
        <Col xs={24} lg={8}>
          <Card title="用户活跃度" loading={loading}>
            <Pie {...pieConfig} height={300} />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]} className="content-row">
        {/* 系统状态 */}
        <Col xs={24} lg={8}>
          <Card title="系统状态" loading={loading}>
            <Column {...columnConfig} height={200} />
            <div className="system-info">
              <Space direction="vertical" style={{ width: '100%' }}>
                <div className="info-item">
                  <Text>系统负载：</Text>
                  <Progress percent={stats.systemLoad} size="small" />
                </div>
                <div className="info-item">
                  <Text>响应时间：</Text>
                  <Text strong>{stats.responseTime}ms</Text>
                </div>
              </Space>
            </div>
          </Card>
        </Col>

        {/* 最近活动 */}
        <Col xs={24} lg={8}>
          <Card title="最近活动" loading={loading}>
            <List
              itemLayout="horizontal"
              dataSource={activities.slice(0, 6)}
              renderItem={(item) => (
                <List.Item>
                  <List.Item.Meta
                    avatar={
                      <Avatar 
                        icon={getActivityIcon(item.type)} 
                        style={{ backgroundColor: getActivityColor(item.status) }}
                      />
                    }
                    title={
                      <Space>
                        <Text strong>{item.user.name}</Text>
                        <Tag color={getActivityColor(item.status)}>
                          {item.type}
                        </Tag>
                      </Space>
                    }
                    description={
                      <div>
                        <Text type="secondary">{item.description}</Text>
                        <br />
                        <Text type="secondary" style={{ fontSize: '12px' }}>
                          {item.time}
                        </Text>
                      </div>
                    }
                  />
                </List.Item>
              )}
            />
          </Card>
        </Col>

        {/* 快捷操作 */}
        <Col xs={24} lg={8}>
          <Card title="快捷操作">
            <Row gutter={[8, 8]}>
              {quickActions.map((action, index) => (
                <Col span={12} key={index}>
                  <Card 
                    size="small" 
                    hoverable 
                    className="quick-action-card"
                    onClick={() => window.location.href = action.path}
                  >
                    <div className="quick-action-content">
                      <div 
                        className="quick-action-icon" 
                        style={{ backgroundColor: action.color }}
                      >
                        {action.icon}
                      </div>
                      <Text strong>{action.title}</Text>
                    </div>
                  </Card>
                </Col>
              ))}
            </Row>
          </Card>
        </Col>
      </Row>

      {/* 性能指标 */}
      <Row gutter={[16, 16]} className="metrics-row">
        <Col span={24}>
          <Card title="性能指标">
            <Row gutter={[16, 16]}>
              <Col xs={24} sm={6}>
                <div className="metric-item">
                  <div className="metric-value">
                    <RiseOutlined style={{ color: '#52c41a' }} />
                    <span>{stats.userGrowth}%</span>
                  </div>
                  <div className="metric-label">用户增长率</div>
                </div>
              </Col>
              <Col xs={24} sm={6}>
                <div className="metric-item">
                  <div className="metric-value">
                    <StarOutlined style={{ color: '#faad14' }} />
                    <span>98.5%</span>
                  </div>
                  <div className="metric-label">系统可用性</div>
                </div>
              </Col>
              <Col xs={24} sm={6}>
                <div className="metric-item">
                  <div className="metric-value">
                    <MessageOutlined style={{ color: '#1890ff' }} />
                    <span>156</span>
                  </div>
                  <div className="metric-label">待处理消息</div>
                </div>
              </Col>
              <Col xs={24} sm={6}>
                <div className="metric-item">
                  <div className="metric-value">
                    <TrophyOutlined style={{ color: '#722ed1' }} />
                    <span>95</span>
                  </div>
                  <div className="metric-label">满意度评分</div>
                </div>
              </Col>
            </Row>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default Dashboard;
