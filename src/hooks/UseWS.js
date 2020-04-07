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
        const objIndex = currPlayers.findIndex(obj => obj.id === player.id);
        if (objIndex !== -1) {
          const clonePlayers = JSON.parse(JSON.stringify(currPlayers));
          clonePlayers[objIndex].shake = false;
          if (!player.show) {
            const filtered = clonePlayers.filter(clonePlayer => clonePlayer.id !== player.id);
            return filtered;
          }
          if (player.collision === true) {
            clonePlayers[objIndex].shake = true;
          }
          clonePlayers[objIndex].p = player.p;
          return clonePlayers;
        }
        //add active:true to THIS user
        let newPlayer
        if (parseData.self.id === parseData.player.id) {
          console.log('self')
          newPlayer = { ...player, active: true }
        } else {
          console.log('not')

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
