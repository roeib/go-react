import { useState, useEffect, useRef, useCallback } from 'react';

export const useWebSocket = (url,bounderies) => {
  const [messages, setMessages] = useState([]);
  const webSocket = useRef(null);
 
  useEffect(() => {
    webSocket.current = new WebSocket(url);
    webSocket.current.onmessage = (event) => {
      const parseData = JSON.parse(event.data);
      const { player } = parseData
      setMessages(currPlayers => {
        //get location from all players
        const objIndex = currPlayers.findIndex(obj => obj.id === player.id);

        if (objIndex !== -1) {
          const clonePlayers = JSON.parse(JSON.stringify(currPlayers));

          //delete anination
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
        if (parseData.self.id === parseData.player.id) {
          newPlayer = { ...player, active: true }
        } else {
          newPlayer = { ...player }
        }
        return [...currPlayers, newPlayer];
      });
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
