import { Button, Image, Input, Spin, Typography } from "antd";
import { Content } from "antd/es/layout/layout";
import { SendOutlined } from "@ant-design/icons";
import React, { useEffect, useRef, useState } from "react";
import Logo from '../assets/logo2.png';
import request from "../utils/fetch";
import Cookies from "js-cookie";

const { Title } = Typography;

const GET_INITIAL_MESSAGES = (idRoom, page, limit) => `/api/messages/room/${idRoom}?page=${page}&limit=${limit}`;
const GET_MESSAGES_OLDER = (idRoom, page, pivotID, limit) => `/api/messages/room/${idRoom}/${pivotID}?page=${page}&limit=${limit}`;

const ChatContent = ({ idRoom, username }) => {
    const id = Cookies.get('id');
    const WEBSOCKET_URL = `ws://localhost:8080/ws?roomId=${idRoom}&userId=${id}`;
    const WEBSOCKET_URL_NOTIFICATION = `ws://localhost:8080/ws/notification?roomId=NOTIFICATION&userId=${id}`;
    const [firstAccess, setFirstAccess] = useState(true);
    const [messages, setMessages] = useState([]);
    const [page, setPage] = useState(0);
    const [isLastPage, setIsLastPage] = useState(false);
    const [loading, setLoading] = useState(false);
    const [pivotID, setPivotID] = useState('');
    const [newMessage, setNewMessage] = useState('');
    const bottomRef = useRef();
    const chatContainerRef = useRef();
    const [stopScrollBottom, setStopScrollBottom] = useState(false);
    const [socket, setSocket] = useState(null);
    const [socketNotification, setSocketNotification] = useState(null);

    useEffect(() => {
        if (!idRoom || !id) return;
        const ws = new WebSocket(WEBSOCKET_URL);

        ws.onopen = () => {
            console.log("WebSocket connection opened");
        };

        ws.onmessage = (event) => {
            const message = JSON.parse(event.data);
            console.log("message data", message);
            setMessages((prev) => [message, ...prev]);
        };

        ws.onerror = (error) => {
            console.error("WebSocket error", error);
        };

        ws.onclose = () => {
            console.log("WebSocket connection closed");
        };

        setSocket(ws);

        return () => {
            ws.close();
        };
    }, [idRoom, id]);

    useEffect(() => {
        if (!stopScrollBottom) {
            bottomRef.current?.scrollIntoView({ behavior: "smooth" });
        }
    }, [messages]);

    useEffect(() => {
        if (idRoom !== '' && username !== '') {
            setFirstAccess(false);
        } else {
            setFirstAccess(true);
        }
    }, [idRoom]);

    useEffect(() => {
        if (firstAccess) return;
        setLoading(true);
        request
            .get(GET_INITIAL_MESSAGES(idRoom, 1, 10))
            .then((res) => {
                if (res.data.data == null) {
                    return;
                }
                setMessages((prev) => [...prev, ...res.data.data]);
                setPivotID(res.data.data[res.data.data.length - 1].id);
            })
            .catch((error) => {
                console.log(error);
            })
            .finally(() => {
                setLoading(false);
            });
    }, [idRoom, firstAccess]);

    const handleSendMessage = () => {
        if (newMessage.trim() === '') return;
        setStopScrollBottom(false);
        const messageToAdd = {
            id: Date.now().toString(),
            id_sender: id,
            username: 'user2',
            id_receiver: idRoom,
            content: newMessage,
        };

        const notification = {
            id_sender: id,
            id_receiver: idRoom,
            content: newMessage,
        }

        const newArr = [messageToAdd].concat(messages);
        setMessages(newArr);
        setNewMessage('');

        if (socket && socket.readyState === WebSocket.OPEN) {
            socket.send(JSON.stringify(messageToAdd));
        } else {
            console.log("channel closed")
        }
    };

    useEffect(() => {
        if (page === 0) return;
        const chatContainer = chatContainerRef.current;
        const currentScrollHeight = chatContainer.scrollHeight;

        setLoading(true);
        request
            .get(GET_MESSAGES_OLDER(idRoom, page, pivotID, 10))
            .then((res) => {
                if (res.data.data === null) {
                    setIsLastPage(true);
                    return;
                }
                setMessages((prev) => [...prev, ...res.data.data]);
                setStopScrollBottom(true);
                chatContainer.scrollTop = chatContainer.scrollHeight - currentScrollHeight;
            })
            .catch((error) => {
                console.log(error);
            })
            .finally(() => {
                setLoading(false);
            });
    }, [page]);

    const handleScroll = (e) => {
        const top = e.target.scrollTop === 0;
        if (top && !loading && !isLastPage) {
            setPage((prev) => prev + 1);
        }
    }

    return (
        <Content style={{ padding: '16px', backgroundColor: '#b4b2b2' }}>
            <div style={{ background: '#fff', padding: '16px', height: '100%', borderRadius: '20px' }}>
                {!firstAccess ? (
                    <>
                        <Title style={{ fontFamily: 'cursive' }} level={4}>
                            Chat with {username}
                        </Title>
                        <div
                            ref={chatContainerRef}
                            style={{
                                marginTop: '16px',
                                height: 'calc(100% - 120px)',
                                overflowY: 'auto',
                                borderTop: 'solid 0.015rem gray',
                                paddingRight: '10px',
                                paddingLeft: '10px',
                            }}
                            onScroll={handleScroll}
                        >
                            {loading && (
                                <div style={{textAlign: 'center', marginTop: '16px'}}>
                                    <Spin size="large" />
                                </div>
                            )}
                            {messages.length > 0 ? (
                                messages.slice().reverse().map((message) => (
                                    <div key={message.id} style={{
                                        display: 'flex',
                                        flexDirection: 'column',
                                        alignItems: message.id_sender === id ? 'flex-end' : 'flex-start',
                                        marginBottom: '10px'
                                    }}>
                                        <span style={{
                                            fontSize: '12px',
                                            color: 'gray',
                                            marginTop: '5px'
                                        }}>{message.id_sender === id ? 'You' : message.username}</span>
                                        <div style={{
                                            padding: '10px',
                                            borderRadius: '10px',
                                            backgroundColor: message.id_sender === id ? 'rgba(166,49,197,0.27)' : '#f0f0f0',
                                            maxWidth: '60%',
                                        }}>
                                            {message.content}
                                        </div>
                                    </div>
                                ))
                            ) : (
                                <div style={{
                                    display: 'flex',
                                    justifyContent: 'center',
                                    alignItems: 'center',
                                    height: '100%',
                                    color: 'gray',
                                    fontSize: '20px',
                                    fontFamily: 'cursive',
                                    padding: '20px'
                                }}>
                                    You have no messages with {username} yet, start your conversation now!
                                </div>
                            )}
                            <div ref={bottomRef} />
                        </div>
                        <div style={{ marginTop: '10px', display: 'flex' }}>
                            <Input.TextArea
                                value={newMessage}
                                onChange={(e) => setNewMessage(e.target.value)}
                                style={{ height: '70px' }}
                                rows={1}
                                placeholder="Nhập tin nhắn..."
                                onPressEnter={(e) => {
                                    e.preventDefault();
                                    handleSendMessage();
                                }}
                            />
                            <Button
                                icon={<SendOutlined />}
                                type="primary"
                                style={{ marginLeft: '10px' }}
                                onClick={handleSendMessage}
                            >
                                Send
                            </Button>
                        </div>
                    </>
                ) : (
                    <div style={{
                        display: 'flex',
                        flexDirection: 'column',
                        justifyContent: 'center',
                        alignItems: 'center',
                        height: 'calc(100vh - 120px)'
                    }}>
                        <img src={Logo} alt="logo" style={{ width: '500px', height: '450px' }} />
                        <Title style={{
                            textAlign: 'center',
                            marginTop: '10px',
                            fontFamily: 'cursive'
                        }} level={2}>Start your messaging now!</Title>
                    </div>
                )}
            </div>
        </Content>
    );
};

export default ChatContent;