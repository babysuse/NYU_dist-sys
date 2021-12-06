import React from "react";
import { Navigate } from 'react-router-dom';

const Home = (props) => {
    if (props.token) {
        return (
            <div className="Home text-center mt-5">
                <h1>Home</h1>
            </div>
        );
    } else {
        return <Navigate to="/login" />
    }
}

export default Home;