import React, { useState, useEffect } from 'react';
import { 
  Card, 
  Row, 
  Col, 
  Avatar, 
  Button, 
  Form, 
  Input, 
  Upload, 
  message, 
  Tabs, 
  Table, 
  Tag, 
  Switch,
  Modal,
  List,
  Typography,
  Space,
  Divider,
  Progress,
  Badge,
  Timeline
} from 'antd';
import { 
  UserOutlined, 
  EditOutlined, 
  CameraOutlined, 
  PhoneOutlined, 
  MailOutlined,
  LockOutlined,
  SecurityScanOutlined,
  MobileOutlined,
  DesktopOutlined,
  ClockCircleOutlined,
  EnvironmentOutlined,
  BellOutlined,
  EyeOutlined,
  DownloadOutlined,
  DeleteOutlined,
  CheckCircleOutlined,
  ExclamationCircleOutlined
} from '@ant-design/icons';
import { useAppSelector, useAppDispatch } from '@/store';
import { updateProfile, changePassword, uploadAvatar, getLoginLogs, getSecuritySettings, updateSecuritySettings } from '@/services/user';
import './Profile.css';

const { TabPane } = Tabs;
const { Title, Text, Paragraph } = Typography;
const { TextArea } = Input;

interface LoginLog {
  id: string;
  ip: string;
  location: string;
  device: string;
  browser: string;
  loginTime: string;
  status: 'success' | 'failed';
}

interface SecuritySettings {
  twoFactorEnabled: boolean;
  emailNotification: boolean;
  smsNotification: boolean;
  loginNotification: boolean;
  passwordExpireDays: number;
  sessionTimeout: number;
}

const Profile: React.FC = () => {
  const { user } = useAppSelector(state => state.auth);
  const dispatch = useAppDispatch();
  
  const [activeTab, setActiveTab] = useState('1');
  const [editMode, setEditMode] = useState(false);
  const [loading, setLoading] = useState(false);
  const [avatarLoading, setAvatarLoading] = useState(false);
  const [passwordModalVisible, setPasswordModalVisible] = useState(false);
  const [loginLogs, setLoginLogs] = useState<LoginLog[]>([]);
  const [securitySettings, setSecuritySettings] = useState<SecuritySettings>({
    twoFactorEnabled: false,
    emailNotification: true,
    smsNotification: true,
    loginNotification: true,
    passwordExpireDays: 90,
    sessionTimeout: 30
  });
  
  const [profileForm] = Form.useForm();
  const [passwordForm] = Form.useForm();

  useEffect(() => {
    if (user) {
      profileForm.setFieldsValue({
        username: user.username,
        email: user.email,
        phone: user.phone,
        realName: user.realName,
        bio: user.bio,
        location: user.location,
        website: user.website
      });
    }
    loadLoginLogs();
    loadSecuritySettings();
  }, [user, profileForm]);

  const loadLoginLogs = async () => {
    try {
      const logs = await getLoginLogs({ page: 1, pageSize: 10 });
      setLoginLogs(logs.data);
    } catch (error) {
      console.error('加载登录日志失败:', error);
    }
  };

  const loadSecuritySettings = async () => {
    try {
      const settings = await getSecuritySettings();
      setSecuritySettings(settings);
    } catch (error) {
      console.error('加载安全设置失败:', error);
    }
  };

  // 更新个人信息
  const handleUpdateProfile = async (values: any) => {
    try {
      setLoading(true);
      await updateProfile(values);
      message.success('个人信息更新成功');
      setEditMode(false);
    } catch (error: any) {
      message.error(error.message || '更新失败');
    } finally {
      setLoading(false);
    }
  };

  // 修改密码
  const handleChangePassword = async (values: any) => {
    try {
      setLoading(true);
      await changePassword({
        oldPassword: values.oldPassword,
        newPassword: values.newPassword
      });
      message.success('密码修改成功');
      setPasswordModalVisible(false);
      passwordForm.resetFields();
    } catch (error: any) {
      message.error(error.message || '密码修改失败');
    } finally {
      setLoading(false);
    }
  };

  // 上传头像
  const handleAvatarUpload = async (file: File) => {
    try {
      setAvatarLoading(true);
      const formData = new FormData();
      formData.append('avatar', file);
      await uploadAvatar(formData);
      message.success('头像上传成功');
    } catch (error: any) {
      message.error(error.message || '头像上传失败');
    } finally {
      setAvatarLoading(false);
    }
  };

  // 更新安全设置
  const handleSecuritySettingChange = async (key: keyof SecuritySettings, value: boolean | number) => {
    try {
      const newSettings = { ...securitySettings, [key]: value };
      await updateSecuritySettings(newSettings);
      setSecuritySettings(newSettings);
      message.success('设置已更新');
    } catch (error: any) {
      message.error(error.message || '设置更新失败');
    }
  };

  // 导出数据
  const handleExportData = () => {
    Modal.confirm({
      title: '导出个人数据',
      content: '确定要导出您的个人数据吗？这将包括您的基本信息、登录记录等。',
      onOk: async () => {
        try {
          // 实际实现中应该调用导出API
          message.success('数据导出请求已提交，我们将通过邮件发送给您');
        } catch (error: any) {
          message.error(error.message || '导出失败');
        }
      }
    });
  };

  // 删除账户
  const handleDeleteAccount = () => {
    Modal.confirm({
      title: '删除账户',
      content: '确定要删除您的账户吗？此操作不可撤销，所有数据将被永久删除。',
      okText: '确认删除',
      okType: 'danger',
      onOk: async () => {
        try {
          // 实际实现中应该调用删除API
          message.success('账户删除请求已提交');
        } catch (error: any) {
          message.error(error.message || '删除失败');
        }
      }
    });
  };

  const uploadProps = {
    showUploadList: false,
    beforeUpload: (file: File) => {
      const isImage = file.type.startsWith('image/');
      if (!isImage) {
        message.error('只能上传图片文件');
        return false;
      }
      const isLt2M = file.size / 1024 / 1024 < 2;
      if (!isLt2M) {
        message.error('图片大小不能超过 2MB');
        return false;
      }
      handleAvatarUpload(file);
      return false;
    }
  };

  const loginLogColumns = [
    {
      title: 'IP地址',
      dataIndex: 'ip',
      key: 'ip',
    },
    {
      title: '位置',
      dataIndex: 'location',
      key: 'location',
    },
    {
      title: '设备',
      dataIndex: 'device',
      key: 'device',
      render: (device: string) => (
        <Space>
          {device.includes('Mobile') ? <MobileOutlined /> : <DesktopOutlined />}
          {device}
        </Space>
      )
    },
    {
      title: '浏览器',
      dataIndex: 'browser',
      key: 'browser',
    },
    {
      title: '登录时间',
      dataIndex: 'loginTime',
      key: 'loginTime',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'success' ? 'green' : 'red'}>
          {status === 'success' ? '成功' : '失败'}
        </Tag>
      )
    }
  ];

  return (
    <div className="profile-container">
      <Row gutter={[24, 24]}>
        <Col xs={24} lg={8}>
          <Card className="profile-card">
            <div className="profile-header">
              <div className="avatar-container">
                <Avatar 
                  size={120} 
                  src={user?.avatar} 
                  icon={<UserOutlined />}
                />
                <Upload {...uploadProps}>
                  <Button 
                    className="avatar-upload-btn"
                    icon={<CameraOutlined />}
                    loading={avatarLoading}
                    type="primary"
                    size="small"
                  >
                    更换头像
                  </Button>
                </Upload>
              </div>
              
              <div className="profile-info">
                <Title level={3}>{user?.realName || user?.username}</Title>
                <Text type="secondary">@{user?.username}</Text>
                <div className="profile-stats">
                  <div className="stat-item">
                    <span className="stat-value">{user?.loginCount || 0}</span>
                    <span className="stat-label">登录次数</span>
                  </div>
                  <div className="stat-item">
                    <span className="stat-value">{user?.points || 0}</span>
                    <span className="stat-label">积分</span>
                  </div>
                </div>
                
                <div className="profile-badges">
                  <Badge count="VIP" style={{ backgroundColor: '#87d068' }} />
                  <Badge count="认证" style={{ backgroundColor: '#108ee9' }} />
                </div>
              </div>
            </div>
          </Card>

          <Card title="账户概览" style={{ marginTop: 16 }}>
            <List
              itemLayout="horizontal"
              dataSource={[
                {
                  title: '账户状态',
                  description: '正常',
                  icon: <CheckCircleOutlined style={{ color: '#52c41a' }} />
                },
                {
                  title: '实名认证',
                  description: user?.verified ? '已认证' : '未认证',
                  icon: user?.verified ? 
                    <CheckCircleOutlined style={{ color: '#52c41a' }} /> :
                    <ExclamationCircleOutlined style={{ color: '#faad14' }} />
                },
                {
                  title: '双因子认证',
                  description: securitySettings.twoFactorEnabled ? '已启用' : '未启用',
                  icon: securitySettings.twoFactorEnabled ?
                    <CheckCircleOutlined style={{ color: '#52c41a' }} /> :
                    <ExclamationCircleOutlined style={{ color: '#ff4d4f' }} />
                }
              ]}
              renderItem={item => (
                <List.Item>
                  <List.Item.Meta
                    avatar={item.icon}
                    title={item.title}
                    description={item.description}
                  />
                </List.Item>
              )}
            />
          </Card>
        </Col>

        <Col xs={24} lg={16}>
          <Card>
            <Tabs activeKey={activeTab} onChange={setActiveTab}>
              <TabPane tab="基本信息" key="1">
                <Form
                  form={profileForm}
                  layout="vertical"
                  onFinish={handleUpdateProfile}
                  disabled={!editMode}
                >
                  <Row gutter={16}>
                    <Col span={12}>
                      <Form.Item label="用户名" name="username">
                        <Input disabled />
                      </Form.Item>
                    </Col>
                    <Col span={12}>
                      <Form.Item label="真实姓名" name="realName">
                        <Input />
                      </Form.Item>
                    </Col>
                  </Row>
                  
                  <Row gutter={16}>
                    <Col span={12}>
                      <Form.Item label="邮箱" name="email">
                        <Input prefix={<MailOutlined />} />
                      </Form.Item>
                    </Col>
                    <Col span={12}>
                      <Form.Item label="手机号" name="phone">
                        <Input prefix={<PhoneOutlined />} />
                      </Form.Item>
                    </Col>
                  </Row>

                  <Row gutter={16}>
                    <Col span={12}>
                      <Form.Item label="所在地" name="location">
                        <Input prefix={<EnvironmentOutlined />} />
                      </Form.Item>
                    </Col>
                    <Col span={12}>
                      <Form.Item label="个人网站" name="website">
                        <Input />
                      </Form.Item>
                    </Col>
                  </Row>

                  <Form.Item label="个人简介" name="bio">
                    <TextArea rows={4} placeholder="介绍一下自己..." />
                  </Form.Item>

                  <Form.Item>
                    <Space>
                      {editMode ? (
                        <>
                          <Button type="primary" htmlType="submit" loading={loading}>
                            保存
                          </Button>
                          <Button onClick={() => setEditMode(false)}>
                            取消
                          </Button>
                        </>
                      ) : (
                        <Button 
                          type="primary" 
                          icon={<EditOutlined />} 
                          onClick={() => setEditMode(true)}
                        >
                          编辑信息
                        </Button>
                      )}
                    </Space>
                  </Form.Item>
                </Form>
              </TabPane>

              <TabPane tab="安全设置" key="2">
                <div className="security-settings">
                  <div className="setting-section">
                    <Title level={4}>密码设置</Title>
                    <div className="setting-item">
                      <div className="setting-info">
                        <Text strong>登录密码</Text>
                        <Text type="secondary">用于登录账户的密码</Text>
                      </div>
                      <Button 
                        icon={<LockOutlined />} 
                        onClick={() => setPasswordModalVisible(true)}
                      >
                        修改密码
                      </Button>
                    </div>
                  </div>

                  <Divider />

                  <div className="setting-section">
                    <Title level={4}>双因子认证</Title>
                    <div className="setting-item">
                      <div className="setting-info">
                        <Text strong>双因子认证</Text>
                        <Text type="secondary">为您的账户添加额外的安全保护</Text>
                      </div>
                      <Switch 
                        checked={securitySettings.twoFactorEnabled}
                        onChange={(checked) => handleSecuritySettingChange('twoFactorEnabled', checked)}
                      />
                    </div>
                  </div>

                  <Divider />

                  <div className="setting-section">
                    <Title level={4}>通知设置</Title>
                    <div className="setting-item">
                      <div className="setting-info">
                        <Text strong>邮件通知</Text>
                        <Text type="secondary">接收重要通知邮件</Text>
                      </div>
                      <Switch 
                        checked={securitySettings.emailNotification}
                        onChange={(checked) => handleSecuritySettingChange('emailNotification', checked)}
                      />
                    </div>
                    
                    <div className="setting-item">
                      <div className="setting-info">
                        <Text strong>短信通知</Text>
                        <Text type="secondary">接收重要通知短信</Text>
                      </div>
                      <Switch 
                        checked={securitySettings.smsNotification}
                        onChange={(checked) => handleSecuritySettingChange('smsNotification', checked)}
                      />
                    </div>
                    
                    <div className="setting-item">
                      <div className="setting-info">
                        <Text strong>登录通知</Text>
                        <Text type="secondary">新设备登录时发送通知</Text>
                      </div>
                      <Switch 
                        checked={securitySettings.loginNotification}
                        onChange={(checked) => handleSecuritySettingChange('loginNotification', checked)}
                      />
                    </div>
                  </div>
                </div>
              </TabPane>

              <TabPane tab="登录日志" key="3">
                <div className="login-logs">
                  <div className="logs-header">
                    <Title level={4}>登录记录</Title>
                    <Text type="secondary">最近10次登录记录</Text>
                  </div>
                  
                  <Table
                    columns={loginLogColumns}
                    dataSource={loginLogs}
                    rowKey="id"
                    pagination={false}
                    scroll={{ x: 800 }}
                  />
                </div>
              </TabPane>

              <TabPane tab="数据管理" key="4">
                <div className="data-management">
                  <div className="section">
                    <Title level={4}>数据导出</Title>
                    <Paragraph type="secondary">
                      您可以导出您的个人数据，包括基本信息、登录记录等。
                    </Paragraph>
                    <Button 
                      icon={<DownloadOutlined />} 
                      type="primary"
                      onClick={handleExportData}
                    >
                      导出数据
                    </Button>
                  </div>

                  <Divider />

                  <div className="section">
                    <Title level={4}>账户删除</Title>
                    <Paragraph type="secondary">
                      删除账户后，所有数据将被永久删除且无法恢复。
                    </Paragraph>
                    <Button 
                      icon={<DeleteOutlined />} 
                      danger
                      onClick={handleDeleteAccount}
                    >
                      删除账户
                    </Button>
                  </div>
                </div>
              </TabPane>
            </Tabs>
          </Card>
        </Col>
      </Row>

      {/* 修改密码模态框 */}
      <Modal
        title="修改密码"
        open={passwordModalVisible}
        onCancel={() => setPasswordModalVisible(false)}
        footer={null}
        destroyOnClose
      >
        <Form
          form={passwordForm}
          layout="vertical"
          onFinish={handleChangePassword}
        >
          <Form.Item
            label="当前密码"
            name="oldPassword"
            rules={[{ required: true, message: '请输入当前密码' }]}
          >
            <Input.Password />
          </Form.Item>
          
          <Form.Item
            label="新密码"
            name="newPassword"
            rules={[
              { required: true, message: '请输入新密码' },
              { min: 8, message: '密码至少8位' }
            ]}
          >
            <Input.Password />
          </Form.Item>
          
          <Form.Item
            label="确认新密码"
            name="confirmPassword"
            dependencies={['newPassword']}
            rules={[
              { required: true, message: '请确认新密码' },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue('newPassword') === value) {
                    return Promise.resolve();
                  }
                  return Promise.reject(new Error('两次输入的密码不一致'));
                },
              }),
            ]}
          >
            <Input.Password />
          </Form.Item>
          
          <Form.Item style={{ marginBottom: 0 }}>
            <Space style={{ width: '100%', justifyContent: 'flex-end' }}>
              <Button onClick={() => setPasswordModalVisible(false)}>
                取消
              </Button>
              <Button type="primary" htmlType="submit" loading={loading}>
                确认修改
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default Profile;
