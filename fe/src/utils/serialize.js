const DEFAULT_PASSWORD = "drill"

export const encodeMsg = (content) => {
    // return AES.encrypt(JSON.stringify(content), DEFAULT_PASSWORD).toString();
    // return AES.encrypt(JSON.stringify(content), DEFAULT_PASSWORD).toString();
    return JSON.stringify(content)
}

export const decodeMsg = (s) => {
    // var bytes = AES.decrypt(s, DEFAULT_PASSWORD);
    // return JSON.parse(bytes.toString(enc.Utf8));
    return JSON.parse(s)
}