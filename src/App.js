
import React, { useEffect, useRef, useLayoutEffect } from 'react';
import './App.css';
import monkey from './monkey.png'
import activemonkey from './activemonkey.png'
import { useWebSocket } from './hooks/UseWS'

// const socket = new WebSocket('ws://localhost:8080/ws');
const MONKEYWIDTH = 100;
const MONKEYHEIGHT = 129;
const STEP = 10;
const NOSTEP = 0

const playerMoves = {
  "ArrowRight": { y: `${NOSTEP}`, x: `${STEP}` },
  "ArrowLeft": { y: `${NOSTEP}`, x: `-${STEP}` },
  "ArrowUp": { y: `${STEP}`, x: `${NOSTEP}` },
  "ArrowDown": { y: `-${STEP}`, x: `${NOSTEP}` }
}

function App() {
  const bodyBoundries = useRef(null);
  const [players,sendMSG,ws] = useWebSocket('ws://localhost:8080/ws',bodyBoundries)


  useLayoutEffect(() => {
    const getBodyBoundries = document.body.getBoundingClientRect();
    bodyBoundries.current = {
      width: (getBodyBoundries.width - MONKEYWIDTH),
      height: (getBodyBoundries.height - MONKEYHEIGHT)
    }
  }, [])
  
  const showKeyCode = e => {
    sendMSG(playerMoves[e.key])
  };

  useEffect(() => {
    document.body.addEventListener('keydown', showKeyCode);
    return () => {
      document.body.removeEventListener('keydown', showKeyCode);
    };
  })

  
  
  // useEffect(() => {
  //   const dd = ws.current
  //   dd.onopen = () => {
  //     dd.send(JSON.stringify(bodyBoundries.current))
  //   }
  //   return () => {
  //     dd.onclose = (e) => {
  //       console.log("socket close connection", e)
  //     }
  //   };
  // }, [ws]);

  // useEffect(() => {
  //   const onMessage = event => {
  //     const parseData = JSON.parse(event.data);
  //     const { player, exception } = parseData
  //     if (parseData.exception.id !== "00000000-0000-0000-0000-000000000000") {
  //       console.log('exception')
  //       setExceptions(currPlayers => {
  //         const objIndex = currPlayers.findIndex(obj => obj.id === exception.id);
  //         if (objIndex !== -1) {
  //           const clonePlayers = JSON.parse(JSON.stringify(currPlayers));
  //           if (!player.show) {
  //             const filtered = clonePlayers.filter(clonePlayer => clonePlayer.id !== player.id);
  //             return filtered;
  //           }
  //           clonePlayers[objIndex].p = exception.p;
  //           return clonePlayers;
  //         }
  //         return [...currPlayers, exception];
  //       });
  //     } else {
  //       console.log('player')
  //       setPlayers(currPlayers => {
  //         const objIndex = currPlayers.findIndex(obj => obj.id === player.id);
  //         if (objIndex !== -1) {
  //           const clonePlayers = JSON.parse(JSON.stringify(currPlayers));
  //           clonePlayers[objIndex].shake = false;
  //           if (!player.show) {
  //             const filtered = clonePlayers.filter(clonePlayer => clonePlayer.id !== player.id);
  //             return filtered;
  //           }
  //           if (player.collision === true) {
  //             clonePlayers[objIndex].shake = true;
  //           }
  //           clonePlayers[objIndex].p = player.p;
  //           return clonePlayers;
  //         }
  //         //add active:true to THIS user
  //         let newPlayer
  //         if (parseData.self.id === parseData.player.id) {
  //           newPlayer = { ...player, active: true }
  //         } else {
  //           newPlayer = { ...player }

  //         }
  //         return [...currPlayers, newPlayer];
  //       });
  //     }

  //   };

    // socket.addEventListener('message', onMessage);

    // return () => socket.removeEventListener('message', onMessage);
  // }, [players, exceptions]);

  return (
    <>
      {players.map(player => {
        return (
          <div key={player.id}>
            <div className={`playerImg ${player.shake === true ? "shake" : ''}`} style={{ bottom: player.p.y + 'px', left: player.p.x + 'px', color: `rgb(${player.color[0]},${player.color[1]},${player.color[2]})` }}>
              <span style={{ position: "absolute", top: "-40px", left: "-25px" }}>{player.exceptionType}score{player.score}</span>
              <img alt="player img" src={player.active ? activemonkey : monkey} />
            </div>
          </div>
        )
      })}
      {/* {exceptions.map(player => {
        return (
          <div key={player.id}>
            <div className={`playerImg`} style={{ bottom: player.p.y + 'px', left: player.p.x + 'px' }}>
              im an exception
            </div>
          </div>
        )
      })} */}

    </>
  );
}
export default App;


