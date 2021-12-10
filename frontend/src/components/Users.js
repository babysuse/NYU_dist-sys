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

        this.getRequest = this.getRequest.bind(this);
        this.followRequest = this.followRequest.bind(this);
    }

    getRequest = async (token, uri) => {
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

    followRequest = (event) => {
        fetch('http://localhost:16008/follow', {
            method: 'POST',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json',
                cookie: `session_token=${this.props.token}`,
            },
            body: JSON.stringify({
                'username': event.target.getAttribute('name'),
                'unfollowing': event.target.innerHTML === 'Unfollow',
            })
        })
        .then(response => response.json())
        .then(data => {
            this.setState({ following: data });
        })
        .catch((err) => {
            console.log(err);
        })
    }

    async componentDidMount() {
        const fp1 = this.getRequest(this.props.token, 'http://localhost:16008/get_users');
        const fp2 = this.getRequest(this.props.token, 'http://localhost:16008/get_following');
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
                        this.state.users && this.state.users.map(u => (
                            <Card>
                                <Card.Body>
                                    <Row>
                                        <Col sm={12} md={6} className="text-md-right">
                                            <h4>{ u }</h4>
                                        </Col>
                                        <Col sm={12} md={6} className="text-md-left">
                                            <Button
                                                key={u}
                                                name={u}
                                                onClick={ this.followRequest }
                                                disabled={ u === this.props.username }
                                            >
                                                { (this.state.following &&
                                                  this.state.following.includes(u)) ||
                                                  u === this.props.username ?
                                                    "Unfollow" : "Follow" }
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
