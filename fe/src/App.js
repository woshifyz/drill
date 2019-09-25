import React, { useEffect } from 'react';
import api from "./api"
import useGlobal from "./store"

function App(props) {
  // const [g, s] = useGlobal()
  // useEffect(() => {
  //   async function initEnv() {
  //     var env = await api.getEnv();
  //     s.changeWsUrl(env.ws_url);
  //   }
  //   initEnv();
  // }, [])
  return (
    <div>
      <div style={{ marginTop: "0px" }}>
        {props.children}
      </div>
    </div>
  )
}

export default App;
