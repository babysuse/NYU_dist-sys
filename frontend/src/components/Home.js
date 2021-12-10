import React from "react";
import { Navigate } from 'react-router-dom';
import { Card, Row } from "react-bootstrap";
import ToggleCreateForm from "./ToggleCreateForm";

class Home extends React.Component {
    constructor (props) {
        super(props);

        this.state = {
            posts: []
        }
        this.getPosts = this.getPosts.bind(this);
    }

    getPosts = (token) => {
        fetch('http://localhost:16008/get_posts', {
            method: 'GET',
            credentials: 'include',
            headers: {
                cookie: `session_token=${token}`
            },
        })
        .then(response => response.json())
        .then(data => this.setState({posts: data}))
        .catch(err => console.log(err));
        // { author: ..., text: ... }
    }

    componentDidMount() {
        this.getPosts(this.props.token);
    }

    render() {
        if (this.props.token) {
            return (
                <div className="Home text-center mt-5">
                    <h1>Posts</h1>
                    {
                        this.state.posts && this.state.posts.map((p, index) => (
                            <Card key={index}>
                                <Card.Body>
                                    <Card.Title>{ p.author }</Card.Title>
                                    <Card.Text>{ p.text }</Card.Text>
                                </Card.Body>
                            </Card>
                        ))
                    }
                    <Row className="justify-content-center my-3">
                        <ToggleCreateForm
                            token={this.props.token}
                            username={this.props.username}
                            pushPost={(p) => {
                                if (this.state.posts) {
                                    this.setState({ posts: [...this.state.posts, p] })
                                } else {
                                    this.setState({ posts: [p] })
                                }
                            }}
                        />
                    </Row>
                </div>
            );
        } else {
            return <Navigate to="/login" />
        }
    }
}

export default Home;