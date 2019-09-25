import React, { useState, useEffect } from 'react'
import { getLocalStorage } from "../utils/tools"
import useGlobal from "../store"
import api from "../api"
import { encodeMsg, decodeMsg } from "../utils/serialize"
import ReconnectingWebSocket from 'reconnecting-websocket';

import Quill from "quill"
import { async } from 'q';

import { AES, enc } from "crypto-js"
import "./Home.css"

var Delta = Quill.import('delta');

var editorModules = {
    syntax: true,
    toolbar: [
        // [{ 'font': [] }],
        ['bold', 'italic', 'underline', 'strike'],
        [{ 'header': '1' }, { 'header': '2' }],
        ['blockquote', 'code-block'],
        [{ 'color': [] }, { 'background': [] }],
        [{ 'script': 'super' }, { 'script': 'sub' }],
        [{ 'list': 'ordered' }, { 'list': 'bullet' }],
        [{ 'align': [] }],
        ['link', 'image', 'video'],
        // ['formula'],
        ['clean']
    ],
};

function concatTypedArrays(a, b) { // a, b TypedArray of same type
    var c = new (a.constructor)(a.length + b.length);
    c.set(a, 0);
    c.set(b, a.length);
    return c;
}

const nickname = Math.random().toString(36).substring(6, 15);

export default function Home(props) {
    const [users, setUsers] = useState([])

    const roomId = props.match.params.roomId

    var password = null;

    var quill = null;
    var ws = null;

    var initWsClient = async () => {
        var env = await api.getEnv();
        const wsUrl = env.ws_url
        ws = new ReconnectingWebSocket(wsUrl, [], {
            maxRetries: 3
        });
        ws.onopen = () => {
            const msg = "j," + encodeMsg({ room_id: roomId, nickname: nickname });
            ws.send(msg);
        }

        ws.onmessage = evt => {
            var message = decodeMsg(evt.data);
            if (!message.hasOwnProperty("action")) {
                return
            }
            if (message.action === "refresh_member") {
                setUsers(users => message.members);
            } else if (message.action === "join") {
                onReceiveJoinMsg(message);
            } else if (message.action === "editor_delta") {
                onReceiveDelta(message);
            }
        }

        ws.onclose = () => {
            console.log('disconnected')
        }
    }

    var pendingDeltaCache = [];
    var remoteVersion = 0;
    var localVersion = 0;

    var syncContent = null;

    var roomEncrypted = roomId[0] === 'S';

    var requirePassword = () => {
        if (roomEncrypted) {
            if (password === null) {
                var inputs = prompt("Enter password to enter room")
                if (typeof inputs === typeof "" && inputs.length > 0) {
                    password = inputs.trim()
                }
            }
        }
    }

    var requireCorrectPassword = () => {
        if (roomEncrypted) {
            var inputs = prompt("Password not valid, please input correct one")
            password = inputs.trim()
        }
    }


    var contentToStr = (content) => {
        if (roomEncrypted) {
            return AES.encrypt(JSON.stringify(content), password).toString();
        } else {
            return JSON.stringify(content);
        }
    }
    var strToContent = (s) => {
        if (roomEncrypted) {
            var prCount = 0
            while (true) {
                prCount += 1
                if (prCount > 3) {
                    break;
                }
                try {
                    var bytes = AES.decrypt(s, password);
                    return JSON.parse(bytes.toString(enc.Utf8));
                } catch (e) {
                    requireCorrectPassword();
                }
            }
            alert("Password error")
        } else {
            return JSON.parse(s);
        }
    }

    var handleChange = (delta, oldDelta, source) => {
        if (ws === null) {
            return;
        }
        if (source === 'user') {
            const val = contentToStr(delta)
            localVersion += 1
            pendingDeltaCache.push([localVersion, delta])

            if (localVersion > 0 && localVersion % 10 === 1) {
                const allContent = contentToStr(quill.getContents());
                const msg = 'e,' + encodeMsg({ room_id: roomId, msg: val, version: localVersion, full: true, all: allContent })
                ws.send(msg);
            } else {
                const msg = 'e,' + encodeMsg({ room_id: roomId, msg: val, version: localVersion, full: false, all: "" })
                ws.send(msg);
            }
        }
    }

    var onReceiveJoinMsg = (message) => {
        setUsers(users => message.members);
        var finalContent = new Delta();
        message.content.forEach(x => {
            const v = new Delta(strToContent(x).ops)
            finalContent = finalContent.compose(v);
        })
        const currentContent = quill.getContents();
        // const diff = currentContent.diff(finalContent);
        // quill.updateContents(diff);
        quill.setContents(finalContent)
        remoteVersion = message.current_version;
        syncContent = quill.getContents();
    }


    var onReceiveDelta = (message) => {
        var newOp = strToContent(message.delta);
        var newDelta = new Delta(newOp.ops);

        remoteVersion = message.current_version;

        syncContent = syncContent.compose(newDelta);

        if (message.sender === nickname) {
            pendingDeltaCache = pendingDeltaCache.filter(i => i[0] !== message.version)
        } else {
            var finalContent = syncContent;
            pendingDeltaCache.forEach(cacheDelta => {
                finalContent = finalContent.compose(cacheDelta[1]);
            });
            const currentContent = quill.getContents();
            const diff = currentContent.diff(finalContent);
            quill.updateContents(diff);
        }
    }

    useEffect(() => {
        requirePassword();

        quill = new Quill('#editor', {
            modules: editorModules,
            theme: 'snow'
        });
        quill.on('text-change', handleChange);

        initWsClient();

        window.addEventListener("hashchange", () => { window.location.reload(); }, false);

        return () => {
            window.removeEventListener("hashchange", () => { window.location.reload(); }, false);
        }
    }, [])

    return (
        <div className="editorContainer">
            <div className="editorMain">
                <div id="editor" className="editorInner">

                </div>
            </div>
            <div className="editerMembers">
                <div className="memberHeader">Members</div>
                <div>
                    {
                        users.map((value, index) => {
                            return <div className={value !== nickname ? "memberOther" : "memberMe"} key={value}>{'U' + value}</div>
                        })
                    }
                </div>
            </div>
        </div>
    )
}
