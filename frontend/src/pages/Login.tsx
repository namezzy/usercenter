import React, { useState, useEffect } from 'react';
import { Form, Input, Button, Checkbox, Alert, Card, Row, Col, Image } from 'antd';
import { UserOutlined, LockOutlined, SafetyOutlined } from '@ant-design/icons';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import { AppDispatch, RootState } from '@/store';
import { loginAsync, clearError } from '@/store/slices/authSlice';
import { authApi } from '@/services/auth';
import { LoginRequest } from '@/types';

const Login: React.FC = () => {
  const [form] = Form.useForm();
  const [captcha, setCaptcha] = useState<{ id: string; img: string } | null>(null);
  const [rememberMe, setRememberMe] = useState(false);
  
  const navigate = useNavigate();
  const location = useLocation();
  const dispatch = useDispatch<AppDispatch>();
  
  const { loading, error, isAuthenticated } = useSelector((state: RootState) => state.auth);

  // 获取验证码
  const getCaptcha = async () => {
    try {
      const captchaData = await authApi.getCaptcha();
      setCaptcha(captchaData);
    } catch (error) {
      console.error('获取验证码失败:', error);
    }
  };

  // 处理登录
  const handleLogin = async (values: any) => {
    try {
      const loginData: LoginRequest = {
        username: values.username,
        password: values.password,
        captcha_id: captcha?.id,
        captcha_code: values.captcha,
        device_info: {
          device_type: 'web',
          device_name: navigator.userAgent,
        },
      };

      await dispatch(loginAsync(loginData)).unwrap();
      
      // 登录成功后跳转
      const from = (location.state as any)?.from?.pathname || '/dashboard';
      navigate(from, { replace: true });
    } catch (error) {
      // 刷新验证码
      getCaptcha();
    }
  };

  // 组件挂载时获取验证码
  useEffect(() => {
    getCaptcha();
    
    // 清除之前的错误
    dispatch(clearError());
  }, [dispatch]);

  // 如果已登录，直接跳转
  useEffect(() => {
    if (isAuthenticated) {
      const from = (location.state as any)?.from?.pathname || '/dashboard';
      navigate(from, { replace: true });
    }
  }, [isAuthenticated, navigate, location]);

  return (
    <div style={{
      minHeight: '100vh',
      background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      padding: '20px'
    }}>
      <Card 
        style={{ 
          width: '100%', 
          maxWidth: 400,
          boxShadow: '0 8px 32px rgba(0,0,0,0.1)',
          borderRadius: 12
        }}
        bodyStyle={{ padding: '40px 32px' }}
      >
        {/* 标题 */}
        <div style={{ textAlign: 'center', marginBottom: 32 }}>
          <h1 style={{ 
            fontSize: 28, 
            fontWeight: 'bold', 
            margin: 0,
            background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
            WebkitBackgroundClip: 'text',
            WebkitTextFillColor: 'transparent'
          }}>
            用户中心
          </h1>
          <p style={{ color: '#666', margin: '8px 0 0 0' }}>
            欢迎回来，请登录您的账户
          </p>
        </div>

        {/* 错误提示 */}
        {error && (
          <Alert
            message={error}
            type="error"
            showIcon
            closable
            style={{ marginBottom: 24 }}
            onClose={() => dispatch(clearError())}
          />
        )}

        {/* 登录表单 */}
        <Form
          form={form}
          name="login"
          size="large"
          onFinish={handleLogin}
          autoComplete="off"
        >
          <Form.Item
            name="username"
            rules={[
              { required: true, message: '请输入用户名/邮箱/手机号' },
            ]}
          >
            <Input
              prefix={<UserOutlined />}
              placeholder="用户名/邮箱/手机号"
            />
          </Form.Item>

          <Form.Item
            name="password"
            rules={[
              { required: true, message: '请输入密码' },
            ]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="密码"
            />
          </Form.Item>

          {/* 验证码 */}
          {captcha && (
            <Form.Item>
              <Row gutter={8}>
                <Col span={14}>
                  <Form.Item
                    name="captcha"
                    noStyle
                    rules={[
                      { required: true, message: '请输入验证码' },
                    ]}
                  >
                    <Input
                      prefix={<SafetyOutlined />}
                      placeholder="验证码"
                    />
                  </Form.Item>
                </Col>
                <Col span={10}>
                  <Image
                    src={captcha.img}
                    alt="验证码"
                    style={{ 
                      width: '100%', 
                      height: 40, 
                      cursor: 'pointer',
                      borderRadius: 4
                    }}
                    preview={false}
                    onClick={getCaptcha}
                  />
                </Col>
              </Row>
            </Form.Item>
          )}

          {/* 记住我和忘记密码 */}
          <Form.Item>
            <div style={{ display: 'flex', justifyContent: 'space-between' }}>
              <Checkbox 
                checked={rememberMe} 
                onChange={(e) => setRememberMe(e.target.checked)}
              >
                记住我
              </Checkbox>
              <Link to="/forgot-password">忘记密码？</Link>
            </div>
          </Form.Item>

          {/* 登录按钮 */}
          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              style={{ width: '100%', height: 48 }}
              loading={loading}
            >
              登录
            </Button>
          </Form.Item>

          {/* 注册链接 */}
          <div style={{ textAlign: 'center' }}>
            <span style={{ color: '#666' }}>还没有账户？</span>
            <Link to="/register" style={{ marginLeft: 8 }}>
              立即注册
            </Link>
          </div>
        </Form>
      </Card>
    </div>
  );
};

export default Login;
