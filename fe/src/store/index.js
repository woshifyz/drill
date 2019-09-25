import React from 'react';
import useGlobalHook from '../utils/useGlobalHook';
import * as actions from "./actions"

const initialState = {
    userId: Math.floor(Math.random() * 100000),
    wsUrl: "",
};

const useGlobal = useGlobalHook(React, initialState, actions);

export default useGlobal;