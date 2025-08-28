import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Form, Input, Button, Card, Tabs, message, Row, Col, Divider, Checkbox } from 'antd';
import { UserOutlined, MailOutlined, PhoneOutlined, LockOutlined, SafetyOutlined } from '@ant-design/icons';
import { useAppDispatch } from '@/store';
import { register, sendVerificationCode } from '@/services/auth';
import './Register.css';

const { TabPane } = Tabs;

interface RegisterFormData {
  username: string;
  email: string;
  phone: string;
  password: string;
  confirmPassword: string;
  verificationCode: string;
  agreement: boolean;
}

const Register: React.FC = () => {
  const [form] = Form.useForm();
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const [loading, setLoading] = useState(false);
  const [sendingCode, setSendingCode] = useState(false);
  const [countdown, setCountdown] = useState(0);
  const [registerType, setRegisterType] = useState<'email' | 'phone'>('email');

  // 发送验证码
  const handleSendCode = async () => {
    try {
      const target = registerType === 'email' 
        ? form.getFieldValue('email')
        : form.getFieldValue('phone');
      
      if (!target) {
        message.error(registerType === 'email' ? '请输入邮箱地址' : '请输入手机号码');
        return;
      }

      setSendingCode(true);
      await sendVerificationCode({
        type: registerType,
        target,
        purpose: 'register'
      });
      
      message.success('验证码已发送');
      
      // 开始倒计时
      setCountdown(60);
      const timer = setInterval(() => {
        setCountdown(prev => {
          if (prev <= 1) {
            clearInterval(timer);
            return 0;
          }
          return prev - 1;
        });
      }, 1000);
      
    } catch (error: any) {
      message.error(error.message || '发送验证码失败');
    } finally {
      setSendingCode(false);
    }
  };

  // 注册处理
  const handleRegister = async (values: RegisterFormData) => {
    try {
      setLoading(true);
      
      const registerData = {
        username: values.username,
        email: registerType === 'email' ? values.email : '',
        phone: registerType === 'phone' ? values.phone : '',
        password: values.password,
        verificationCode: values.verificationCode,
        registerType
      };

      const response = await register(registerData);
      
      message.success('注册成功！请登录');
      navigate('/login');
      
    } catch (error: any) {
      message.error(error.message || '注册失败');
    } finally {
      setLoading(false);
    }
  };

  // 密码强度检查
  const checkPasswordStrength = (password: string) => {
    if (!password) return 0;
    
    let strength = 0;
    if (password.length >= 8) strength++;
    if (/[a-z]/.test(password)) strength++;
    if (/[A-Z]/.test(password)) strength++;
    if (/[0-9]/.test(password)) strength++;
    if (/[^A-Za-z0-9]/.test(password)) strength++;
    
    return strength;
  };

  const getPasswordStrengthText = (strength: number) => {
    switch (strength) {
      case 0:
      case 1:
        return { text: '弱', color: '#ff4d4f' };
      case 2:
      case 3:
        return { text: '中', color: '#faad14' };
      case 4:
      case 5:
        return { text: '强', color: '#52c41a' };
      default:
        return { text: '', color: '' };
    }
  };

  const renderEmailRegister = () => (
    <>
      <Form.Item
        name="email"
        rules={[
          { required: true, message: '请输入邮箱地址' },
          { type: 'email', message: '请输入有效的邮箱地址' }
        ]}
      >
        <Input
          prefix={<MailOutlined />}
          placeholder="邮箱地址"
          size="large"
        />
      </Form.Item>
      
      <Form.Item
        name="verificationCode"
        rules={[{ required: true, message: '请输入邮箱验证码' }]}
      >
        <Row gutter={8}>
          <Col span={16}>
            <Input
              prefix={<SafetyOutlined />}
              placeholder="邮箱验证码"
              size="large"
            />
          </Col>
          <Col span={8}>
            <Button
              size="large"
              onClick={handleSendCode}
              loading={sendingCode}
              disabled={countdown > 0}
              style={{ width: '100%' }}
            >
              {countdown > 0 ? `${countdown}s` : '获取验证码'}
            </Button>
          </Col>
        </Row>
      </Form.Item>
    </>
  );

  const renderPhoneRegister = () => (
    <>
      <Form.Item
        name="phone"
        rules={[
          { required: true, message: '请输入手机号码' },
          { pattern: /^1[3-9]\d{9}$/, message: '请输入有效的手机号码' }
        ]}
      >
        <Input
          prefix={<PhoneOutlined />}
          placeholder="手机号码"
          size="large"
        />
      </Form.Item>
      
      <Form.Item
        name="verificationCode"
        rules={[{ required: true, message: '请输入短信验证码' }]}
      >
        <Row gutter={8}>
          <Col span={16}>
            <Input
              prefix={<SafetyOutlined />}
              placeholder="短信验证码"
              size="large"
            />
          </Col>
          <Col span={8}>
            <Button
              size="large"
              onClick={handleSendCode}
              loading={sendingCode}
              disabled={countdown > 0}
              style={{ width: '100%' }}
            >
              {countdown > 0 ? `${countdown}s` : '获取验证码'}
            </Button>
          </Col>
        </Row>
      </Form.Item>
    </>
  );

  return (
    <div className="register-container">
      <div className="register-content">
        <Card className="register-card">
          <div className="register-header">
            <h1>用户注册</h1>
            <p>创建您的账户</p>
          </div>

          <Form
            form={form}
            onFinish={handleRegister}
            layout="vertical"
            size="large"
          >
            <Form.Item
              name="username"
              rules={[
                { required: true, message: '请输入用户名' },
                { min: 3, max: 20, message: '用户名长度为3-20位' },
                { pattern: /^[a-zA-Z0-9_-]+$/, message: '用户名只能包含字母、数字、下划线和连字符' }
              ]}
            >
              <Input
                prefix={<UserOutlined />}
                placeholder="用户名"
                size="large"
              />
            </Form.Item>

            <Tabs 
              activeKey={registerType} 
              onChange={(key) => setRegisterType(key as 'email' | 'phone')}
              centered
            >
              <TabPane tab="邮箱注册" key="email">
                {renderEmailRegister()}
              </TabPane>
              <TabPane tab="手机注册" key="phone">
                {renderPhoneRegister()}
              </TabPane>
            </Tabs>

            <Form.Item
              name="password"
              rules={[
                { required: true, message: '请输入密码' },
                { min: 8, message: '密码至少8位' },
                { 
                  validator: (_, value) => {
                    if (!value) return Promise.resolve();
                    const strength = checkPasswordStrength(value);
                    if (strength < 3) {
                      return Promise.reject(new Error('密码强度太弱，请包含大小写字母、数字和特殊字符'));
                    }
                    return Promise.resolve();
                  }
                }
              ]}
            >
              <Input.Password
                prefix={<LockOutlined />}
                placeholder="密码"
                size="large"
              />
            </Form.Item>

            <Form.Item
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
              <Input.Password
                prefix={<LockOutlined />}
                placeholder="确认密码"
                size="large"
              />
            </Form.Item>

            <Form.Item
              name="agreement"
              valuePropName="checked"
              rules={[{ required: true, message: '请阅读并同意用户协议' }]}
            >
              <Checkbox>
                我已阅读并同意 <Link to="/terms">《用户协议》</Link> 和 <Link to="/privacy">《隐私政策》</Link>
              </Checkbox>
            </Form.Item>

            <Form.Item>
              <Button
                type="primary"
                htmlType="submit"
                loading={loading}
                style={{ width: '100%' }}
                size="large"
              >
                注册
              </Button>
            </Form.Item>
          </Form>

          <Divider>其他方式注册</Divider>
          
          <div className="register-footer">
            <p>
              已有账户？ <Link to="/login">立即登录</Link>
            </p>
          </div>
        </Card>
      </div>
    </div>
  );
};

export default Register;
