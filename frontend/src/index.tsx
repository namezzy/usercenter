import React from 'react';
import ReactDOM from 'react-dom/client';
import { Provider } from 'react-redux';
import { PersistGate } from 'redux-persist/integration/react';
import { BrowserRouter } from 'react-router-dom';
import { ConfigProvider, App as AntdApp, Spin } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import { store, persistor } from './store';
import App from './App';
import './index.css';

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);

root.render(
  <React.StrictMode>
    <Provider store={store}>
      <PersistGate loading={<Spin size="large" style={{ 
        position: 'fixed', 
        top: '50%', 
        left: '50%', 
        transform: 'translate(-50%, -50%)' 
      }} />} persistor={persistor}>
        <ConfigProvider locale={zhCN}>
          <BrowserRouter>
            <AntdApp>
              <App />
            </AntdApp>
          </BrowserRouter>
        </ConfigProvider>
      </PersistGate>
    </Provider>
  </React.StrictMode>
);
