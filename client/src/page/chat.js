import React, {useEffect, useRef, useState} from 'react';
import {Layout} from 'antd';
import HeaderComponent from "../component/header";
import SidebarComponent from "../component/sidebar";
import ChatContent from "../component/chatBox";
import './style.css'
import {useNavigate} from "react-router-dom";

const Chat = () => {
    const nav = useNavigate();
    const [idRoom, setIdRoom] = useState('');
    const [username, setUsername] = useState('');
    const handleSelectFriend = (id, username) => {
        setIdRoom(id);
        setUsername(username);
    }

    const reset = () => {
        setIdRoom('');
        setUsername('');
        nav('/chat');
    }

    return (
        <Layout style={{height: '100vh'}}>
            <HeaderComponent reset={reset}/>
            <Layout>
                <SidebarComponent onSelectFriend={handleSelectFriend}/>
                <ChatContent username={username} idRoom={idRoom}/>
            </Layout>
        </Layout>
    )
};

export default Chat;
