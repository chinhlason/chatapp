import React, {useEffect, useState} from "react";
import request from "../utils/fetch";
import {Avatar, Badge, Input, List, Spin} from "antd";
import {SearchOutlined, UserOutlined} from "@ant-design/icons";
import Sider from "antd/es/layout/Sider";
import LoadingIcon from "antd/lib/button/LoadingIcon";
import {useNavigate} from "react-router-dom";
import Cookies from "js-cookie";

const GET_LIST_FRIENDS_URL = (page, limit, pivot) => `/api/friend/list?page=${page}&limit=${limit}&time=${pivot}`

const SidebarComponent = ({onSelectFriend}) => {
    const id = Cookies.get('id');
    const username = Cookies.get('username');
    const WEBSOCKET_URL = `ws://localhost:8080/ws/notification?userId=${id}&username=${username}`;
    const nav = useNavigate();
    const [friends, setFriends] = useState([]);
    const [page, setPage] = useState(0);
    const [isLastPage, setIsLastPage] = useState(false);
    const [loading, setLoading] = useState(false);
    const [selectedId, setSelectedId] = useState(null);
    const [socket, setSocket] = useState(null);
    const [onlineState, setOnlineState] = useState(false);
    const [pivot, setPivot] = useState("");

    const handleSelectFriend = (id, username, idRoom) => {
        setSelectedId(id);
        onSelectFriend(idRoom, username);
        nav(`/chat?id_room=${idRoom}&username=${username}`);
    }

    useEffect(() => {
        if (!id) return;

        const ws = new WebSocket(WEBSOCKET_URL);

        ws.onopen = () => {
            console.log("WebSocket NOTIFICATION CHANNEL connection opened");
        };

        ws.onmessage = (event) => {
            const message = JSON.parse(event.data);
            console.log("WebSocket message notification received", message);
            const idSender = message.id_sender;

            //id_receiver : NOTIFICATION_1, trim only get 1
            const idReceiver = message.id_receiver.split('_')[1];

            if (message.type === 'ONLINE_NOTIFICATION') {
                setFriends((prevFriends) => {
                    return prevFriends.map((friend) => {
                        if (friend.id_friend == idSender) {
                            return {
                                ...friend,
                                is_online: message.content === 'online',
                            };
                        }
                        return friend;
                    });
                });
            }

            else {
                const newFriend = {
                    friend_username: message.username_sender,
                    id_friend: idSender,
                    id_room: idReceiver,
                    interaction_at: new Date().toISOString(),
                    is_online: true,
                    is_read: false,
                };

                setFriends((prevFriends) => {
                    const existingIndex = prevFriends.findIndex(friend => friend.id_friend == idSender);
                    if (existingIndex !== -1) {
                        const updatedFriends = [...prevFriends];
                        const [existingFriend] = updatedFriends.splice(existingIndex, 1);
                        return [existingFriend, ...updatedFriends];
                    } else if (existingIndex === 0) {
                        return [prevFriends];
                    }
                    return [newFriend, ...prevFriends];
                });
            }
        }


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
    }, [id]);

    useEffect(() => {
        setLoading(true);
        request
            .get(GET_LIST_FRIENDS_URL(1, 10, pivot))
            .then((res) => {
                if (res.data.data == null) {
                    setIsLastPage(true);
                    return;
                }
                setFriends(res.data.data);
                setPivot(res.data.data[res.data.data.length - 1].interaction_at);
            })
            .catch((error) => {
                console.log(error);
            })
            .finally(() => {
                setLoading(false);
            });
    }, []);

    const loadFriends = (page) => {
        setLoading(true);
        request
            .get(GET_LIST_FRIENDS_URL(page, 10, pivot))
            .then((res) => {
                if (res.data.data == null) {
                    setIsLastPage(true);
                    return;
                }
                setFriends((prev) => [...prev, ...res.data.data]);
                setPivot(res.data.data[res.data.data.length - 1].interaction_at);
            })
            .catch((error) => {
                console.log(error);
            })
            .finally(() => {
                setLoading(false);
            });
    }

    useEffect(() => {
        if (page >= 1) {
            loadFriends(page);
        }
    }, [page]);

    const handleScroll = (e) => {
        const bottom = e.target.scrollHeight - e.target.scrollTop === e.target.clientHeight;
        if (bottom && !loading && !isLastPage) {
            setPage((prev) => prev + 1);
        }
    }

    return (
        <Sider width={200} style={{backgroundColor: '#600080', color: 'fff'}}>
            <div style={{padding: '10px'}}>
                <Input
                    placeholder="Search"
                    prefix={<SearchOutlined />}
                    style={{marginBottom: '10px'}}
                />
                <div style={{
                    color: '#fff',
                    fontSize: '20px',
                    fontWeight: 'bold',
                    marginBottom: '10px'
                }}>Friends list</div>
                <List
                    className="custom-scrollbar"
                    style={{backgroundColor: '#600080', maxHeight: 'calc(100vh - 195px)', overflowY: 'auto'}}
                    itemLayout="horizontal"
                    dataSource={friends}
                    renderItem={(item) => (
                        <List.Item
                            className="new-noti"
                            style={{
                                cursor: 'pointer',
                                backgroundColor: selectedId === item.id ? '#8a008a' : '#600080',
                                fontWeight: selectedId === item.id ? 'bold' : 'normal',
                                borderRadius: '10px',
                                marginRight: '10px',
                        }}
                            onClick={() => handleSelectFriend(item.id, item.username, item.id_room)}
                        >
                            <List.Item.Meta
                                style={{display: 'flex', alignItems: 'center', paddingLeft: '10px'}}
                                avatar={<Badge
                                    color={item.is_online ? 'green' : 'gray'}
                                    dot
                                    offset={[-2, 25]}
                                >
                                    <Avatar icon={<UserOutlined />} />
                                </Badge>}
                                // avatar={<Avatar icon={<UserOutlined />} />}
                                title={item.friend_username}
                                description="Description"
                            />
                        </List.Item>
                    )}
                    onScroll={handleScroll}
                >
                    {loading && (
                        <div style={{textAlign: 'center', marginTop: '16px'}}>
                            <Spin size="large" style={{color: '#fff'}}/>
                        </div>
                    )}
                </List>
            </div>
        </Sider>
    );
};

export default SidebarComponent;