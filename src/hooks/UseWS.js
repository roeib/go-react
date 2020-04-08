import { useReducer, useEffect, useRef, useCallback } from 'react';
const initialState = [];
const counterReducer = (state, action) => {
  switch (action.type) {
    case "INCOMINGMSG":
      const { player } = action.by
      const objIndex = state.findIndex(obj => obj.id === player.id);
      if (objIndex !== -1) {
        const clonePlayers = JSON.parse(JSON.stringify(state));
        //delete animation
        clonePlayers[objIndex].shake = false;

        //check if player need to be seen in the screen if not remove from players
        if (!player.show) {
          const filtered = clonePlayers.filter(clonePlayer => clonePlayer.id !== player.id);
          return filtered;
        }

        //check if player hit bounderies and add animation
        if (player.collision) clonePlayers[objIndex].shake = true

        //change player cordinates on screen
        clonePlayers[objIndex].p = player.p;
        return clonePlayers;
      }
      //add active:true to the user that open connection with socket
      let newPlayer
      if (action.by.self.id === action.by.player.id) {
        newPlayer = { ...player, active: true }
      } else {
        newPlayer = { ...player }
      }
      return [...state, newPlayer];
    default:
      throw new Error();
  }
};

export const useWebSocket = (url, bounderies) => {
  const [messages, dispatch] = useReducer(counterReducer, initialState);
  const webSocket = useRef(null);

  useEffect(() => {
    webSocket.current = new WebSocket(url);
    webSocket.current.onmessage = (event) => {
      const parseData = JSON.parse(event.data);
      dispatch({ type: "INCOMINGMSG", by: parseData });
    };
  }, [url]);


  useEffect(() => {
    webSocket.current.onopen = () => {
      webSocket.current.send(JSON.stringify(bounderies.current))
    }
    return () => {
      webSocket.current.onclose = (e) => {
        console.log("socket close connection", e)
      }
    };
  }, [bounderies]);

  const sendMessage = useCallback(message => {
    webSocket.current.send(JSON.stringify(message));
  }, [webSocket]);

  return [messages, sendMessage]
};
