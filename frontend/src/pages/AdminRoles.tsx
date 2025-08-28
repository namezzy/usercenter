import React, { useState, useEffect } from 'react';
import {
  Card,
  Table,
  Button,
  Form,
  Input,
  Modal,
  Select,
  Tree,
  Tag,
  Space,
  Popconfirm,
  message,
  Row,
  Col,
  Typography,
  Tabs,
  Transfer,
  Checkbox,
  Divider
} from 'antd';
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  SecurityScanOutlined,
  TeamOutlined,
  UserOutlined,
  SettingOutlined,
  CheckOutlined,
  CloseOutlined
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import type { DataNode } from 'antd/es/tree';
import { useAppSelector } from '@/store';
import { getRoles, createRole, updateRole, deleteRole, getPermissions } from '@/services/admin';
import './AdminRoles.css';

const { Title, Text } = Typography;
const { TabPane } = Tabs;
const { Option } = Select;

interface Role {
  id: string;
  name: string;
  description: string;
  permissions: string[];
  userCount: number;
  status: 'active' | 'inactive';
  createdAt: string;
  updatedAt: string;
}

interface Permission {
  id: string;
  name: string;
  code: string;
  description: string;
  resource: string;
  action: string;
  children?: Permission[];
}

const AdminRoles: React.FC = () => {
  const [roles, setRoles] = useState<Role[]>([]);
  const [permissions, setPermissions] = useState<Permission[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [permissionModalVisible, setPermissionModalVisible] = useState(false);
  const [editingRole, setEditingRole] = useState<Role | null>(null);
  const [selectedRole, setSelectedRole] = useState<Role | null>(null);
  const [selectedPermissions, setSelectedPermissions] = useState<string[]>([]);
  
  const [form] = Form.useForm();

  useEffect(() => {
    loadRoles();
    loadPermissions();
  }, []);

  const loadRoles = async () => {
    try {
      setLoading(true);
      const response = await getRoles();
      setRoles(response.data);
    } catch (error) {
      message.error('加载角色列表失败');
    } finally {
      setLoading(false);
    }
  };

  const loadPermissions = async () => {
    try {
      const response = await getPermissions();
      setPermissions(response.data);
    } catch (error) {
      console.error('加载权限列表失败:', error);
    }
  };

  const handleCreateRole = () => {
    setEditingRole(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEditRole = (role: Role) => {
    setEditingRole(role);
    form.setFieldsValue(role);
    setModalVisible(true);
  };

  const handleSaveRole = async (values: any) => {
    try {
      if (editingRole) {
        await updateRole(editingRole.id, values);
        message.success('角色更新成功');
      } else {
        await createRole(values);
        message.success('角色创建成功');
      }
      setModalVisible(false);
      loadRoles();
    } catch (error: any) {
      message.error(error.message || '操作失败');
    }
  };

  const handleDeleteRole = async (roleId: string) => {
    try {
      await deleteRole(roleId);
      message.success('角色删除成功');
      loadRoles();
    } catch (error: any) {
      message.error(error.message || '删除失败');
    }
  };

  const handleManagePermissions = (role: Role) => {
    setSelectedRole(role);
    setSelectedPermissions(role.permissions);
    setPermissionModalVisible(true);
  };

  const handleSavePermissions = async () => {
    if (!selectedRole) return;
    
    try {
      await updateRole(selectedRole.id, {
        ...selectedRole,
        permissions: selectedPermissions
      });
      message.success('权限分配成功');
      setPermissionModalVisible(false);
      loadRoles();
    } catch (error: any) {
      message.error(error.message || '权限分配失败');
    }
  };

  const getStatusColor = (status: string) => {
    return status === 'active' ? 'green' : 'red';
  };

  const getStatusText = (status: string) => {
    return status === 'active' ? '启用' : '禁用';
  };

  // 构建权限树
  const buildPermissionTree = (permissions: Permission[]): DataNode[] => {
    const groupedPermissions: { [key: string]: Permission[] } = {};
    
    permissions.forEach(permission => {
      if (!groupedPermissions[permission.resource]) {
        groupedPermissions[permission.resource] = [];
      }
      groupedPermissions[permission.resource].push(permission);
    });

    return Object.keys(groupedPermissions).map(resource => ({
      title: resource,
      key: resource,
      children: groupedPermissions[resource].map(permission => ({
        title: `${permission.name} (${permission.action})`,
        key: permission.id,
        isLeaf: true,
      }))
    }));
  };

  const roleColumns: ColumnsType<Role> = [
    {
      title: '角色名称',
      dataIndex: 'name',
      key: 'name',
      render: (name, record) => (
        <div>
          <div className="role-name">{name}</div>
          <div className="role-description">{record.description}</div>
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
      title: '用户数量',
      dataIndex: 'userCount',
      key: 'userCount',
      width: 100,
      render: (count) => (
        <Tag color="blue">
          <UserOutlined /> {count}
        </Tag>
      )
    },
    {
      title: '权限数量',
      dataIndex: 'permissions',
      key: 'permissions',
      width: 120,
      render: (permissions) => (
        <Tag color="purple">
          <SecurityScanOutlined /> {permissions.length}
        </Tag>
      )
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
          <Button
            type="text"
            icon={<EditOutlined />}
            onClick={() => handleEditRole(record)}
          >
            编辑
          </Button>
          
          <Button
            type="text"
            icon={<SecurityScanOutlined />}
            onClick={() => handleManagePermissions(record)}
          >
            权限
          </Button>
          
          <Popconfirm
            title="确定要删除该角色吗？"
            onConfirm={() => handleDeleteRole(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button
              type="text"
              danger
              icon={<DeleteOutlined />}
            >
              删除
            </Button>
          </Popconfirm>
        </Space>
      )
    }
  ];

  return (
    <div className="admin-roles-container">
      <Card>
        <div className="page-header">
          <Title level={3}>
            <TeamOutlined /> 角色权限管理
          </Title>
          <Button
            type="primary"
            icon={<PlusOutlined />}
            onClick={handleCreateRole}
          >
            新增角色
          </Button>
        </div>

        <Table
          columns={roleColumns}
          dataSource={roles}
          rowKey="id"
          loading={loading}
          pagination={{
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total, range) => 
              `第 ${range[0]}-${range[1]} 条/共 ${total} 条`,
          }}
          scroll={{ x: 1000 }}
        />
      </Card>

      {/* 新增/编辑角色模态框 */}
      <Modal
        title={editingRole ? '编辑角色' : '新增角色'}
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        footer={null}
        width={600}
        destroyOnClose
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSaveRole}
        >
          <Form.Item
            label="角色名称"
            name="name"
            rules={[
              { required: true, message: '请输入角色名称' },
              { min: 2, max: 20, message: '角色名称长度为2-20位' }
            ]}
          >
            <Input placeholder="请输入角色名称" />
          </Form.Item>
          
          <Form.Item
            label="角色描述"
            name="description"
            rules={[{ required: true, message: '请输入角色描述' }]}
          >
            <Input.TextArea 
              rows={3} 
              placeholder="请输入角色描述"
              maxLength={200}
              showCount
            />
          </Form.Item>
          
          <Form.Item
            label="状态"
            name="status"
            rules={[{ required: true, message: '请选择状态' }]}
          >
            <Select placeholder="请选择状态">
              <Option value="active">启用</Option>
              <Option value="inactive">禁用</Option>
            </Select>
          </Form.Item>

          <Form.Item>
            <Space style={{ width: '100%', justifyContent: 'flex-end' }}>
              <Button onClick={() => setModalVisible(false)}>
                取消
              </Button>
              <Button type="primary" htmlType="submit">
                {editingRole ? '更新' : '创建'}
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>

      {/* 权限分配模态框 */}
      <Modal
        title={`为角色"${selectedRole?.name}"分配权限`}
        open={permissionModalVisible}
        onCancel={() => setPermissionModalVisible(false)}
        width={800}
        footer={
          <Space>
            <Button onClick={() => setPermissionModalVisible(false)}>
              取消
            </Button>
            <Button type="primary" onClick={handleSavePermissions}>
              保存
            </Button>
          </Space>
        }
      >
        <div className="permission-assignment">
          <Tabs defaultActiveKey="tree">
            <TabPane tab="权限树" key="tree">
              <div className="permission-tree-container">
                <Tree
                  checkable
                  treeData={buildPermissionTree(permissions)}
                  checkedKeys={selectedPermissions}
                  onCheck={(checkedKeys) => {
                    setSelectedPermissions(checkedKeys as string[]);
                  }}
                  height={400}
                />
              </div>
            </TabPane>
            
            <TabPane tab="权限列表" key="list">
              <div className="permission-list-container">
                {permissions.map(permission => (
                  <div key={permission.id} className="permission-item">
                    <Checkbox
                      checked={selectedPermissions.includes(permission.id)}
                      onChange={(e) => {
                        if (e.target.checked) {
                          setSelectedPermissions([...selectedPermissions, permission.id]);
                        } else {
                          setSelectedPermissions(selectedPermissions.filter(id => id !== permission.id));
                        }
                      }}
                    >
                      <div className="permission-info">
                        <div className="permission-name">{permission.name}</div>
                        <div className="permission-description">{permission.description}</div>
                        <div className="permission-meta">
                          <Tag size="small">{permission.resource}</Tag>
                          <Tag size="small">{permission.action}</Tag>
                        </div>
                      </div>
                    </Checkbox>
                  </div>
                ))}
              </div>
            </TabPane>
          </Tabs>
        </div>
      </Modal>
    </div>
  );
};

export default AdminRoles;
