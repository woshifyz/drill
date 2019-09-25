export const changeUserId = (store, userId) => {
    store.setState({ userId });
};

export const changeWsUrl = (store, url) => {
    store.setState({ wsUrl: url });
};