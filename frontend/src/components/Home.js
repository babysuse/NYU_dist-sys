import React from "react";
import { Navigate } from 'react-router-dom';
import { Card } from "react-bootstrap";

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
                        this.state.posts.map((p, index) => (
                            <Card key={index}>
                                <Card.Body>
                                    <Card.Title>{ p.author }</Card.Title>
                                    <Card.Text>{ p.text }</Card.Text>
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

export default Home;