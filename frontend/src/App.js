import React, { useState } from 'react';
import { Navbar, Nav } from 'react-bootstrap';
import { LinkContainer } from 'react-router-bootstrap';
import UserRoutes from './UserRoutes';
import { useNavigate } from 'react-router-dom';

const Navigation = (props) => {
    const navigate = useNavigate();

    return (
        <>
            <Navbar collapseOnSelect bg="light" expand="md" className="mb=3">
                <LinkContainer to="/">
                    <Navbar.Brand className="font-weight-bold text-muted">Dist-Sys</Navbar.Brand>
                </LinkContainer>
                <Navbar.Toggle />

                <Navbar.Toggle />
                <Navbar.Collapse className="justify-content-end">
                    <Nav activeKey={window.location.pathname}>
                        {props.token ? (
                            <Nav.Link onClick={() => {
                                props.setToken('');
                                navigate('/login')
                            }}>Logout</Nav.Link>
                        ) : (
                            <>
                                <LinkContainer to="/login">
                                    <Nav.Link>Login</Nav.Link>
                                </LinkContainer>
                                <LinkContainer to="/signup">
                                    <Nav.Link>Signup</Nav.Link>
                                </LinkContainer>
                            </>
                        )}
                    </Nav>
                </Navbar.Collapse>
            </Navbar>
        </>
    );
}

const App = () => {
    const [token, setToken] = useState('testing...');

    return (
        <div className="App container py-3">
            <Navigation token={token} setToken={setToken} />
            <UserRoutes token={token} setToken={setToken} />
        </div>
    );
}

export default App;
