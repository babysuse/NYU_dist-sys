import { Routes, Route } from 'react-router-dom';
import Home from './components/Home';
import Login from './components/Login';
import Signup from './components/Signup';
import Users from './components/Users';
import NotFound from './components/NotFound';

const UserRoutes = (props) => {
    return (
        <Routes>
            <Route path="/" element={<Home token={props.token} />} />
            <Route path="/login" element={<Login
                setToken={props.setToken}
                setUsername={props.setUsername}
            />} />
            <Route path="/signup" element={<Signup
                setToken={props.setToken}
                setUsername={props.setUsername}
            />} />
            <Route path="/users" element={<Users
                token={props.token}
                username={props.username}
            />} />
            <Route path="*" element={<NotFound />} />
        </Routes>
    )
}

export default UserRoutes;
