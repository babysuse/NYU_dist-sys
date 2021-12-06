import React, { useState } from 'react';
import { Form, Button, Row, Col } from 'react-bootstrap';
import Cookies from 'js-cookie';
import { useNavigate } from 'react-router-dom';

// TODO: login error (wrong password / username not exist)
const LoginForm = (props) => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const navigate = useNavigate();

    const validateForm = () => {
        return username.length > 0 && password.length > 0;
    }

    const handleSubmit = async (e) => {
        e.preventDefault();
        // document.cookie.split(";").forEach(function(c) { document.cookie = c.replace(/^ +/, "").replace(/=.*/, "=;expires=" + new Date().toUTCString() + ";path=/"); });

        // const result = await fetch ('http://localhost:16008/query', {
        //     method: 'POST',
        //     headers: {
        //         'Content-Type': 'application/json'
        //     },
        //     body: JSON.stringify({
        //         query: `
        //             mutation: {
        //                 login(input: {username: "test", password: "test123"})
        //             }
        //         `
        //     })
        // })
        // console.log(result)
        await fetch('http://localhost:16008/login', {
            method: 'POST',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                username,
                password,
            })
        })

        props.setToken(Cookies.get('session_token'));
        props.setUsername(username);
        navigate('/');
    }

    return (
        <Form onSubmit={handleSubmit}>
            <Form.Group>
                <Form.Label>Username</Form.Label>
                <Form.Control
                    autoFocus
                    type="text"
                    value={username}
                    placeholder="Username"
                    onChange={(e) => setUsername(e.target.value)}
                />
            </Form.Group>
            <Form.Group>
                <Form.Label>Password</Form.Label>
                <Form.Control
                    type="password"
                    value={password}
                    placeholder="Password"
                    onChange={(e) => setPassword(e.target.value)}
                />
            </Form.Group>
            <Row className="justify-content-center">
                <Button type="submit" disabled={!validateForm()}>
                    Login
                </Button>
            </Row>
        </Form>
    );
}

const Login = (props) => {
    return (
        <div className="Login mt-5">
            <Col md={{ span: 6, offset: 3 }}>
                <LoginForm setToken={props.setToken} setUsername={props.setUsername} />
            </Col>
        </div>
    )
}

export default Login;