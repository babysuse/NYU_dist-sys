import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Form, Button, Row, Col } from 'react-bootstrap';
import Cookies from 'js-cookie';

// TODO: signup error (username being signed up)
const Signup = (props) => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const navigate = useNavigate();

    const validateForm = () => {
        return username.length > 0 &&
                password.length > 0 &&
                confirmPassword.length > 0 &&
                password === confirmPassword;
    }

    const handleSubmit = async (e) => {
        e.preventDefault();
        await fetch('http://localhost:16008/signup', {
            method: 'POST',
            'Access-Control-Allow-Origin': 'http://localhost:3000',
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
            <Form.Group>
                <Form.Label>Confirm password</Form.Label>
                <Form.Control
                    type="password"
                    value={confirmPassword}
                    placeholder="Confirm password"
                    onChange={(e) => setConfirmPassword(e.target.value)}
                />
            </Form.Group>
            <Row className="justify-content-center">
                <Button type="submit" disabled={!validateForm()}>
                    Sign up
                </Button>
            </Row>
        </Form>
    );
}

export default Signup;