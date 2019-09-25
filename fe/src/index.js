import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import * as serviceWorker from './serviceWorker';
import { HashRouter as Router, Route, Switch, Redirect } from "react-router-dom";
import Home from "./pages/Home"
import { transitions, positions, Provider as AlertProvider } from 'react-alert'
import ReactModal from "react-modal"

const options = {
  position: positions.TOP_CENTER,
  timeout: 800,
  offset: "300px",
  transition: transitions.FADE
};

const customStyles = {
  content: {
    marginLeft: "auto",
    marginRight: "auto",
    width: "200px",
    textAlign: "center",
    position: "relative",
    height: "auto",
    minHeight: "100% !important",
    borderColor: "#a7342d",
    borderBottomLeftRadius: "15px 255px",
    borderBottomRightRadius: "225px 15px",
    borderTopLeftRadius: "255px 10px",
    borderTopRightRadius: "10px 125px"
  }
};

const AlertTemplate = ({ style, options, message, close }) => (
  <ReactModal
    style={customStyles}
    closeTimeoutMS={200}
    isOpen={true}
    onRequestClose={() => { return true; }}
  >
    {message}
  </ReactModal>
);

const randomInitRoomId = '/room/R' + Math.random().toString(36).substring(2, 15);

ReactDOM.render(
  <AlertProvider template={AlertTemplate} {...options}>
    <App>
      <Router>
        <Switch>
          <Route path="/room/:roomId" component={Home} />
          <Redirect from='/' to={randomInitRoomId} />
        </Switch>
      </Router>
    </App>
  </AlertProvider>,
  document.getElementById('root'));

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
