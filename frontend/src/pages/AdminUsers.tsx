import React, { useState, useEffect } from 'react';
import {
  Card,
  Table,
  Button,
  Form,
  Input,
  Modal,
  Select,
  Tag,
  Space,
  Popconfirm,
  message,
  Row,
  Col,
  Avatar,
  Badge,
  Tooltip,
  Upload,
  Switch,
  DatePicker,
  Typography,
  Divider,
  Drawer,
  Tabs
} from 'antd';
import {
  UserOutlined,
  EditOutlined,
  DeleteOutlined,
  PlusOutlined,
  SearchOutlined,
  ExportOutlined,
  ImportOutlined,
  ReloadOutlined,
  EyeOutlined,
  LockOutlined,
  UnlockOutlined,
  DownOutlined,
  SettingOutlined,
  SecurityScanOutlined,
  TeamOutlined
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { useAppSelector } from '@/store';
import { 
  getUsers, 
  createUser, 
  updateUser, 
  deleteUser, 
  exportUsers, 
  importUsers,
  lockUser,
  unlockUser,
  resetPassword,
  getRoles,
  getPermissions
} from '@/services/admin';
import './AdminUsers.css';

const { Option } = Select;
const { Title, Text } = Typography;
const { TabPane } = Tabs;
const { RangePicker } = DatePicker;

interface User {
  id: string;
  username: string;
  email: string;
  phone: string;
  realName: string;
  avatar: string;
  status: 'active' | 'inactive' | 'locked';
  roles: string[];
  lastLoginTime: string;
  createdAt: string;
  loginCount: number;
  verified: boolean;
}

interface Role {
  id: string;
  name: string;
  description: string;
  permissions: string[];
}

interface QueryParams {
  page: number;
  pageSize: number;
  keyword?: string;
  status?: string;
  role?: string;
  dateRange?: [string, string];
}

const AdminUsers: React.FC = () => {
  const [users, setUsers] = useState<User[]>([]);
  const [roles, setRoles] = useState<Role[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);
  const [detailVisible, setDetailVisible] = useState(false);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
  
  const [queryParams, setQueryParams] = useState<QueryParams>({
    page: 1,
    pageSize: 10
  });

  const [form] = Form.useForm();
  const [searchForm] = Form.useForm();

  useEffect(() => {
    loadUsers();
    loadRoles();
  }, [queryParams]);

  const loadUsers = async () => {
    try {
      setLoading(true);
      const response = await getUsers(queryParams);
      setUsers(response.data);
      setTotal(response.total);
    } catch (error) {
      message.error('加载用户列表失败');
    } finally {
      setLoading(false);
    }
  };

  const loadRoles = async () => {
    try {
      const response = await getRoles();
      setRoles(response.data);
    } catch (error) {
      console.error('加载角色列表失败:', error);
    }
  };

  const handleSearch = (values: any) => {
    setQueryParams({
      ...queryParams,
      page: 1,
      keyword: values.keyword,
      status: values.status,
      role: values.role,
      dateRange: values.dateRange ? values.dateRange.map((d: any) => d.format('YYYY-MM-DD')) : undefined
    });
  };

  const handleReset = () => {
    searchForm.resetFields();
    setQueryParams({
      page: 1,
      pageSize: 10
    });
  };

  const handleCreateUser = () => {
    setEditingUser(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEditUser = (user: User) => {
    setEditingUser(user);
    form.setFieldsValue({
      ...user,
      roleIds: user.roles
    });
    setModalVisible(true);
  };

  const handleSaveUser = async (values: any) => {
    try {
      if (editingUser) {
        await updateUser(editingUser.id, values);
        message.success('用户更新成功');
      } else {
        await createUser(values);
        message.success('用户创建成功');
      }
      setModalVisible(false);
      loadUsers();
    } catch (error: any) {
      message.error(error.message || '操作失败');
    }
  };

  const handleDeleteUser = async (userId: string) => {
    try {
      await deleteUser(userId);
      message.success('用户删除成功');
      loadUsers();
    } catch (error: any) {
      message.error(error.message || '删除失败');
    }
  };

  const handleBatchDelete = async () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请选择要删除的用户');
      return;
    }
    
    Modal.confirm({
      title: '批量删除用户',
      content: `确定要删除选中的 ${selectedRowKeys.length} 个用户吗？`,
      onOk: async () => {
        try {
          for (const userId of selectedRowKeys) {
            await deleteUser(userId as string);
          }
          message.success('批量删除成功');
          setSelectedRowKeys([]);
          loadUsers();
        } catch (error: any) {
          message.error(error.message || '批量删除失败');
        }
      }
    });
  };

  const handleLockUser = async (userId: string, lock: boolean) => {
    try {
      if (lock) {
        await lockUser(userId);
        message.success('用户已锁定');
      } else {
        await unlockUser(userId);
        message.success('用户已解锁');
      }
      loadUsers();
    } catch (error: any) {
      message.error(error.message || '操作失败');
    }
  };

  const handleResetPassword = async (userId: string) => {
    Modal.confirm({
      title: '重置密码',
      content: '确定要重置该用户的密码吗？新密码将通过邮件发送给用户。',
      onOk: async () => {
        try {
          await resetPassword(userId);
          message.success('密码重置成功');
        } catch (error: any) {
          message.error(error.message || '密码重置失败');
        }
      }
    });
  };

  const handleExport = async () => {
    try {
      await exportUsers(queryParams);
      message.success('导出成功');
    } catch (error: any) {
      message.error(error.message || '导出失败');
    }
  };

  const handleImport = async (file: File) => {
    try {
      const formData = new FormData();
      formData.append('file', file);
      await importUsers(formData);
      message.success('导入成功');
      loadUsers();
    } catch (error: any) {
      message.error(error.message || '导入失败');
    }
  };

  const handleViewDetail = (user: User) => {
    setSelectedUser(user);
    setDetailVisible(true);
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active':
        return 'green';
      case 'inactive':
        return 'orange';
      case 'locked':
        return 'red';
      default:
        return 'default';
    }
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case 'active':
        return '正常';
      case 'inactive':
        return '未激活';
      case 'locked':
        return '已锁定';
      default:
        return '未知';
    }
  };

  const columns: ColumnsType<User> = [
    {
      title: '用户',
      key: 'user',
      width: 200,
      render: (_, record) => (
        <div className="user-info">
          <Avatar src={record.avatar} icon={<UserOutlined />} />
          <div className="user-details">
            <div className="username">{record.username}</div>
            <div className="real-name">{record.realName}</div>
          </div>
        </div>
      )
    },
    {
      title: '联系方式',
      key: 'contact',
      width: 200,
      render: (_, record) => (
        <div>
          <div>{record.email}</div>
          <div className="phone">{record.phone}</div>
        </div>
      )
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status) => (
        <Tag color={getStatusColor(status)}>
          {getStatusText(status)}
        </Tag>
      )
    },
    {
      title: '角色',
      dataIndex: 'roles',
      key: 'roles',
      width: 150,
      render: (roles) => (
        <div>
          {roles.map((role: string) => (
            <Tag key={role} color="blue">{role}</Tag>
          ))}
        </div>
      )
    },
    {
      title: '认证状态',
      dataIndex: 'verified',
      key: 'verified',
      width: 100,
      render: (verified) => (
        <Badge
          status={verified ? 'success' : 'warning'}
          text={verified ? '已认证' : '未认证'}
        />
      )
    },
    {
      title: '最后登录',
      dataIndex: 'lastLoginTime',
      key: 'lastLoginTime',
      width: 150,
      render: (time) => time || '从未登录'
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      width: 150
    },
    {
      title: '操作',
      key: 'actions',
      width: 200,
      fixed: 'right',
      render: (_, record) => (
        <Space>
          <Tooltip title="查看详情">
            <Button
              type="text"
              icon={<EyeOutlined />}
              onClick={() => handleViewDetail(record)}
            />
          </Tooltip>
          
          <Tooltip title="编辑">
            <Button
              type="text"
              icon={<EditOutlined />}
              onClick={() => handleEditUser(record)}
            />
          </Tooltip>
          
          <Tooltip title={record.status === 'locked' ? '解锁' : '锁定'}>
            <Button
              type="text"
              icon={record.status === 'locked' ? <UnlockOutlined /> : <LockOutlined />}
              onClick={() => handleLockUser(record.id, record.status !== 'locked')}
            />
          </Tooltip>
          
          <Tooltip title="重置密码">
            <Button
              type="text"
              icon={<SettingOutlined />}
              onClick={() => handleResetPassword(record.id)}
            />
          </Tooltip>
          
          <Popconfirm
            title="确定要删除该用户吗？"
            onConfirm={() => handleDeleteUser(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Tooltip title="删除">
              <Button
                type="text"
                danger
                icon={<DeleteOutlined />}
              />
            </Tooltip>
          </Popconfirm>
        </Space>
      )
    }
  ];

  const rowSelection = {
    selectedRowKeys,
    onChange: setSelectedRowKeys,
  };

  return (
    <div className="admin-users-container">
      <Card>
        <div className="page-header">
          <Title level={3}>用户管理</Title>
          <Space>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={handleCreateUser}
            >
              新增用户
            </Button>
            
            <Upload
              accept=".xlsx,.xls,.csv"
              showUploadList={false}
              beforeUpload={(file) => {
                handleImport(file);
                return false;
              }}
            >
              <Button icon={<ImportOutlined />}>
                导入用户
              </Button>
            </Upload>
            
            <Button
              icon={<ExportOutlined />}
              onClick={handleExport}
            >
              导出用户
            </Button>
          </Space>
        </div>

        <div className="search-section">
          <Form
            form={searchForm}
            layout="inline"
            onFinish={handleSearch}
            style={{ marginBottom: 16 }}
          >
            <Form.Item name="keyword">
              <Input
                placeholder="搜索用户名、邮箱、手机号"
                prefix={<SearchOutlined />}
                style={{ width: 250 }}
              />
            </Form.Item>
            
            <Form.Item name="status">
              <Select placeholder="状态" style={{ width: 120 }} allowClear>
                <Option value="active">正常</Option>
                <Option value="inactive">未激活</Option>
                <Option value="locked">已锁定</Option>
              </Select>
            </Form.Item>
            
            <Form.Item name="role">
              <Select placeholder="角色" style={{ width: 120 }} allowClear>
                {roles.map(role => (
                  <Option key={role.id} value={role.name}>
                    {role.name}
                  </Option>
                ))}
              </Select>
            </Form.Item>
            
            <Form.Item name="dateRange">
              <RangePicker placeholder={['开始日期', '结束日期']} />
            </Form.Item>
            
            <Form.Item>
              <Space>
                <Button type="primary" htmlType="submit" icon={<SearchOutlined />}>
                  搜索
                </Button>
                <Button onClick={handleReset} icon={<ReloadOutlined />}>
                  重置
                </Button>
              </Space>
            </Form.Item>
          </Form>
        </div>

        {selectedRowKeys.length > 0 && (
          <div className="batch-actions">
            <Space>
              <Text>已选择 {selectedRowKeys.length} 项</Text>
              <Button
                danger
                icon={<DeleteOutlined />}
                onClick={handleBatchDelete}
              >
                批量删除
              </Button>
              <Button onClick={() => setSelectedRowKeys([])}>
                取消选择
              </Button>
            </Space>
          </div>
        )}

        <Table
          columns={columns}
          dataSource={users}
          rowKey="id"
          loading={loading}
          rowSelection={rowSelection}
          pagination={{
            current: queryParams.page,
            pageSize: queryParams.pageSize,
            total,
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total, range) => 
              `第 ${range[0]}-${range[1]} 条/共 ${total} 条`,
            onChange: (page, pageSize) => {
              setQueryParams({ ...queryParams, page, pageSize: pageSize! });
            }
          }}
          scroll={{ x: 1200 }}
        />
      </Card>

      {/* 新增/编辑用户模态框 */}
      <Modal
        title={editingUser ? '编辑用户' : '新增用户'}
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        footer={null}
        width={600}
        destroyOnClose
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSaveUser}
        >
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                label="用户名"
                name="username"
                rules={[
                  { required: true, message: '请输入用户名' },
                  { min: 3, max: 20, message: '用户名长度为3-20位' }
                ]}
              >
                <Input />
              </Form.Item>
            </Col>
            
            <Col span={12}>
              <Form.Item
                label="真实姓名"
                name="realName"
                rules={[{ required: true, message: '请输入真实姓名' }]}
              >
                <Input />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                label="邮箱"
                name="email"
                rules={[
                  { required: true, message: '请输入邮箱' },
                  { type: 'email', message: '请输入有效的邮箱地址' }
                ]}
              >
                <Input />
              </Form.Item>
            </Col>
            
            <Col span={12}>
              <Form.Item
                label="手机号"
                name="phone"
                rules={[
                  { required: true, message: '请输入手机号' },
                  { pattern: /^1[3-9]\d{9}$/, message: '请输入有效的手机号' }
                ]}
              >
                <Input />
              </Form.Item>
            </Col>
          </Row>

          {!editingUser && (
            <Row gutter={16}>
              <Col span={12}>
                <Form.Item
                  label="密码"
                  name="password"
                  rules={[
                    { required: true, message: '请输入密码' },
                    { min: 8, message: '密码至少8位' }
                  ]}
                >
                  <Input.Password />
                </Form.Item>
              </Col>
              
              <Col span={12}>
                <Form.Item
                  label="确认密码"
                  name="confirmPassword"
                  dependencies={['password']}
                  rules={[
                    { required: true, message: '请确认密码' },
                    ({ getFieldValue }) => ({
                      validator(_, value) {
                        if (!value || getFieldValue('password') === value) {
                          return Promise.resolve();
                        }
                        return Promise.reject(new Error('两次输入的密码不一致'));
                      },
                    }),
                  ]}
                >
                  <Input.Password />
                </Form.Item>
              </Col>
            </Row>
          )}

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                label="角色"
                name="roleIds"
                rules={[{ required: true, message: '请选择角色' }]}
              >
                <Select mode="multiple" placeholder="选择角色">
                  {roles.map(role => (
                    <Option key={role.id} value={role.id}>
                      {role.name}
                    </Option>
                  ))}
                </Select>
              </Form.Item>
            </Col>
            
            <Col span={12}>
              <Form.Item
                label="状态"
                name="status"
                rules={[{ required: true, message: '请选择状态' }]}
              >
                <Select>
                  <Option value="active">正常</Option>
                  <Option value="inactive">未激活</Option>
                  <Option value="locked">已锁定</Option>
                </Select>
              </Form.Item>
            </Col>
          </Row>

          <Form.Item>
            <Space style={{ width: '100%', justifyContent: 'flex-end' }}>
              <Button onClick={() => setModalVisible(false)}>
                取消
              </Button>
              <Button type="primary" htmlType="submit">
                {editingUser ? '更新' : '创建'}
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>

      {/* 用户详情抽屉 */}
      <Drawer
        title="用户详情"
        placement="right"
        width={600}
        open={detailVisible}
        onClose={() => setDetailVisible(false)}
      >
        {selectedUser && (
          <div className="user-detail">
            <div className="detail-header">
              <Avatar size={80} src={selectedUser.avatar} icon={<UserOutlined />} />
              <div className="detail-info">
                <Title level={4}>{selectedUser.realName}</Title>
                <Text type="secondary">@{selectedUser.username}</Text>
                <div style={{ marginTop: 8 }}>
                  <Tag color={getStatusColor(selectedUser.status)}>
                    {getStatusText(selectedUser.status)}
                  </Tag>
                  {selectedUser.verified && (
                    <Tag color="green">已认证</Tag>
                  )}
                </div>
              </div>
            </div>

            <Divider />

            <Tabs defaultActiveKey="basic">
              <TabPane tab="基本信息" key="basic">
                <div className="detail-section">
                  <div className="detail-item">
                    <span className="label">邮箱：</span>
                    <span>{selectedUser.email}</span>
                  </div>
                  <div className="detail-item">
                    <span className="label">手机号：</span>
                    <span>{selectedUser.phone}</span>
                  </div>
                  <div className="detail-item">
                    <span className="label">角色：</span>
                    <span>
                      {selectedUser.roles.map(role => (
                        <Tag key={role} color="blue">{role}</Tag>
                      ))}
                    </span>
                  </div>
                  <div className="detail-item">
                    <span className="label">登录次数：</span>
                    <span>{selectedUser.loginCount}</span>
                  </div>
                  <div className="detail-item">
                    <span className="label">最后登录：</span>
                    <span>{selectedUser.lastLoginTime || '从未登录'}</span>
                  </div>
                  <div className="detail-item">
                    <span className="label">创建时间：</span>
                    <span>{selectedUser.createdAt}</span>
                  </div>
                </div>
              </TabPane>
              
              <TabPane tab="操作记录" key="logs">
                <div className="detail-section">
                  <Text type="secondary">暂无操作记录</Text>
                </div>
              </TabPane>
            </Tabs>
          </div>
        )}
      </Drawer>
    </div>
  );
};

export default AdminUsers;
