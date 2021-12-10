import { useState } from 'react';
import { Button, Form, Row } from 'react-bootstrap';

const ToggleCreateForm = (props) => {
    const [open, setOpen] = useState(false);
    const [text, setText] = useState('');

    const createPost = (e) => {
        e.preventDefault();
        fetch('http://localhost:16008/createpost', {
            method: 'POST',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json',
                cookie: `session_token=${props.token}`,
            },
            body: JSON.stringify({
                text,
            }),
        })

        props.pushPost({ author: props.username, text });
        setOpen(false);
    }

    if (open) {
        return (
            <div className="ToggleCreateForm">
                <Form onSubmit={createPost}>
                    <Form.Group>
                        <Form.Label><h5>New Post</h5></Form.Label>
                        <Form.Control
                            as="textarea"
                            placeholder="Share your idea..."
                            onChange={(e) => setText(e.target.value)}
                        />
                    </Form.Group>
                    <Row className="justify-content-center">
                        <div className="mx-3">
                            <Button variant="primary" type="submit">Publish</Button>
                        </div>
                        <div className="mx-3">
                            <Button variant="danger" onClick={() => setOpen(false)}>Cancel</Button>
                        </div>
                    </Row>
                </Form>
            </div>
        )
    } else {
        return (
            <div className="ToggleCreateForm">
                <Button onClick={() => setOpen(true)}>Create</Button>
            </div>
        );
    }
}

export default ToggleCreateForm;