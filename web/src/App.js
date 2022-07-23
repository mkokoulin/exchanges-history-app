import React from 'react';
import {useEffect, useState} from 'react';
import { Routes, Route, Link } from "react-router-dom";
import { Layout, Menu } from 'antd';
import Calculation from "./Calculation";
import Upload from "./Upload";
import './App.css';
import 'antd/dist/antd.css';
import api from "./api";

const { Header, Content, Footer } = Layout;

const HistoryInfoContext = React.createContext({
    cryptoCurrencies: [],
    payMethods: []
});

export {HistoryInfoContext}

function App() {
    const [error, setError] = useState(null);
    const [data, setData] = useState([]);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        setLoading(true);
        api.get('/history/info')
            .then(res => {
                setData(res.data)
            })
            .catch(err => {
                setError(err)
            })
        setLoading(false);
    }, []);

    return (
      <Layout>
        <Header style={{ position: 'fixed', zIndex: 1, width: '100%' }}>
          <div className="logo" />
          <Menu
              theme="dark"
              mode="horizontal"
              defaultSelectedKeys={['1']}
              items={[{
                key: 1,
                label: (
                    <Link to="/">History</Link>
                )
              },
              {
                  key: 2,
                  label: (
                      <Link to="upload">Upload</Link>
                  )
              }]}
          />
        </Header>
        <Content className="site-layout" style={{ padding: '0 50px', marginTop: 64 }}>
          <div className="site-layout-background" style={{ padding: 24, height: 'calc(100vh - 135px)' }}>
              <HistoryInfoContext.Provider value={data}>
                  <Routes>
                      <Route path="/" element={<Calculation />} />
                      <Route path="upload" element={<Upload />} />
                  </Routes>
              </HistoryInfoContext.Provider>
          </div>
        </Content>
        <Footer style={{ textAlign: 'center' }}>Exchange History ©2022</Footer>
      </Layout>
  );
}

export default App;
