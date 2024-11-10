import React, {useState} from "react";
import {Avatar, Badge, Button, Dropdown, Input, Menu, Typography} from "antd";
import {BellOutlined, SearchOutlined, UserOutlined} from "@ant-design/icons";
import {Header} from "antd/es/layout/layout";
import {useNavigate} from "react-router-dom";

const {Title} = Typography;

const searchMenu = (
    <Menu>
        <Menu.Item key="1">Kết quả 1 cho</Menu.Item>
        <Menu.Item key="2">Kết quả 2 cho</Menu.Item>
        <Menu.Item key="3">Kết quả 3 cho</Menu.Item>
    </Menu>
);

const notificationMenu = (
    <Menu>
        <Menu.Item key="1">Thông báo 1 <Button>Add Friend</Button></Menu.Item>
        <Menu.Item key="2">Thông báo 2</Menu.Item>
        <Menu.Item key="3">Thông báo 3</Menu.Item>
    </Menu>
);

const HeaderComponent = ({reset}) => {
    const [searchVisible, setSearchVisible] = useState(false);
    const [searchValue, setSearchValue] = useState('');
    const [notificationVisible, setNotificationVisible] = useState(false);
    const nav = useNavigate()
    const handleSearchChange = (e) => {
        setSearchValue(e.target.value);
        setSearchVisible(e.target.value.length > 0); // Hiện dropdown khi có giá trị
    };
    return (
        <Header style={{
            background: 'linear-gradient(90deg, rgba(167,18,208,1) 0%, rgba(120,81,122,1) 99%)'
            , color: '#fff', padding: '0 20px'
        }}>
            <div style={{
                display: 'flex', justifyContent: 'space-between',
                alignItems: 'center', height: '100%'
            }}>
                <Title level={3} style={{color: '#fff', margin: 0, cursor: 'pointer'}} onClick={reset}>
                    Chat With WebSocket
                </Title>
                <Dropdown
                    overlay={searchMenu} visible={searchVisible} placement="bottomRight"
                    onVisibleChange={setSearchVisible}>
                    <Input
                        placeholder="Searching for friends here..."
                        prefix={<SearchOutlined/>}
                        style={{width: 400}}
                        onChange={handleSearchChange}
                        value={searchValue}
                    />
                </Dropdown>

                <div style={
                    {
                        width: '90px',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'space-between'
                    }
                }>
                    <Dropdown overlay={notificationMenu} visible={notificationVisible} placement="bottomRight"
                              onVisibleChange={setNotificationVisible}>
                        <Badge count={'!'} offset={[10, 0]}>
                            <BellOutlined style={{fontSize: 20, color: '#fff', cursor: 'pointer'}}
                                          onClick={() => setNotificationVisible(!notificationVisible)}/>
                        </Badge>
                    </Dropdown>

                    <Avatar icon={<UserOutlined/>}/>
                </div>
            </div>
        </Header>)
};

export default HeaderComponent;