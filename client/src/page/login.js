import React, { useState } from "react";
import { Form, Input, Button, notification } from "antd";
import request from "../utils/fetch";
import {useNavigate} from "react-router-dom";
import Cookies from "js-cookie";

const LOGIN_URL = "/user/login"

const Login = () => {
    const [loading, setLoading] = useState(false);
    const nav = useNavigate()
    const handleLogin = (values) => {
        setLoading(true);
        request
            .post(LOGIN_URL, values)
            .then((res) => {
                const response = res.data
                Cookies.set("token", response.data.token);
                Cookies.set("username", response.data.username);
                Cookies.set("id", response.data.id);
                setTimeout(() => {
                    setLoading(false);
                    nav('/chat')
                    notification.success({
                        message: "Login success!",
                        description: "Welcome to ChatApp.",
                    });
                }, 1000);
            })
            .catch((error) => {
                setTimeout(() => {
                    setLoading(false);
                    notification.error({
                        message: "Login failed!",
                        description: "Please try gain!",
                    });
                }, 1000);
            });
    };

    return (
        <div className="login-container" style={styles.container}>
            <div style={styles.card}>
                <h1 style={styles.title}>LOGIN</h1>
                <Form
                    name="login"
                    onFinish={handleLogin}
                    layout="vertical"
                    initialValues={{ remember: true }}
                >
                    <Form.Item
                        label="Username"
                        name="username"
                        rules={[{ required: true, message: "This field is required!" }]}
                    >
                        <Input placeholder="Username" />
                    </Form.Item>

                    <Form.Item
                        label="Password"
                        name="password"
                        rules={[{ required: true, message: "This field is required!" }]}
                        style={styles.form}
                    >
                        <Input.Password placeholder="Password" />
                    </Form.Item>

                    <Form.Item>
                        <Button
                            type="primary"
                            htmlType="submit"
                            block
                            loading={loading}
                            style={styles.button}
                        >
                            Login
                        </Button>
                    </Form.Item>
                </Form>
            </div>
        </div>
    );
};

const styles = {
    container: {
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        height: "100vh",
        background: "linear-gradient(90deg, rgba(167,18,208,1) 0%, rgba(120,81,122,1) 99%)",
    },
    card: {
        width: 400,
        padding: "20px",
        backgroundColor: "white",
        borderRadius: "8px",
        boxShadow: "0 2px 8px rgba(0, 0, 0, 0.1)",
    },
    title: {
        textAlign: "center",
        marginBottom: "20px",
        fontWeight: "bold",
        fontSize: "24px",
    },
    button: {
        marginTop: "10px",
    },
    form: {
        marginTop: "25px",
    }
};

export default Login;
