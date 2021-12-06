import React from 'react';
import { Button, Card, Col, Row } from 'react-bootstrap';
import { Navigate } from 'react-router-dom';

class Users extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            users: [],
            following: [],
        }

        this.request = this.request.bind(this);
    }

    request = async (token, uri) => {
        try {
            let response = await fetch(uri, {
                method: 'GET',
                credentials: 'include',
                headers: {
                    cookie: `session_token=${token}`,
                },
            });
            return response.json()
        } catch (err) {
            console.log(err)
        }
    }

    async componentDidMount() {
        const fp1 = this.request(this.props.token, 'http://localhost:16008/users');
        const fp2 = this.request(this.props.token, 'http://localhost:16008/following');
        const data = await Promise.all([fp1, fp2])
        const users = data[0];
        const following = data[1];
        this.setState({ users, following })
    }

    render() {
        if (this.props.token) {
            return (
                <div className="Users text-center mt-5">
                    <h1>Users</h1>
                    {
                        this.state.users.map(u => (
                            <Card>
                                <Card.Body>
                                    <Row>
                                        <Col sm={12} md={6} className="text-md-right">
                                            <h4>{ u }</h4>
                                        </Col>
                                        <Col sm={12} md={6} className="text-md-left">
                                            <Button disabled={ u === this.props.username }>
                                                { this.state.following.includes(u) ? "Follow" : "Unfollow" }
                                            </Button>
                                        </Col>
                                    </Row>
                                </Card.Body>
                            </Card>
                        ))
                    }
                </div>
            );
        } else {
            return <Navigate to="/login" />
        }
    }
}

export default Users;