import axios from "axios"
import {TOKEN_KEY} from "./cst";

var HOST = "";

// HOST = "http://localhost:9000";

/*
{
    "code": 0,
    "error_msg": "",
    "result": {
        "ws_url": "ws://localhost:8080"
    }
}
*/
const parseResponse = (resp) => {
  const data = resp.data;
  if (data.code === 0) {
    return data.result;
  } else {
    const msg = data.error_msg;
    if (msg === 'noauth') {
      window.location.href = "/login"
    } else {
      alert(msg);
      throw msg;
    }
  }
};

// axios.interceptors.request.use(function(config) {
//   const token = localStorage.getItem(TOKEN_KEY);

//   if ( token != null ) {
//     config.headers["X-OT"] = token
//   }

//   return config;
// }, function(err) {
//   return Promise.reject(err);
// });

axios.defaults.timeout = 5000;

export default {
  getEnv: async () => {
    const resp = await axios.post(HOST + "/api/v1/util/env");
    return parseResponse(resp);
  },
  sendRegisterCode: async (email) => {
    const resp = await axios.post(HOST + "/api/v1/user_join/register_code/", JSON.stringify({email: email}));
    return parseResponse(resp);
  },
  register: async (email, token, password) => {
    const resp = await axios.post(HOST + "/api/v1/user_join/register/", JSON.stringify({email: email, token: token, password: password}));
    return parseResponse(resp);
  } ,
  login: async (email, token, password) => {
    const resp = await axios.post(HOST + "/api/v1/user_join/login/", JSON.stringify({email: email, password: password}));
    return parseResponse(resp);
  } ,
  user_info: async () => {
    const resp = await axios.get(HOST + "/api/v1/user/info/");
    return parseResponse(resp);
  } ,
  listBooks: async () => {
    const resp = await axios.get(HOST + "/api/v1/book/list/");
    return parseResponse(resp);
  },
  fetchQiniuToken: async (key) => {
    const resp = await axios.post(HOST + "/api/v1/util/qiniu_token/", JSON.stringify({key: key}));
    return parseResponse(resp);
  }
}
