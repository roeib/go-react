
import React, { useEffect, useState, useRef, useLayoutEffect } from 'react';
import './App.css';
import monkey from './monkey.png'
import activemonkey from './activemonkey.png'

const socket = new WebSocket('ws://localhost:8080/ws');
const MONKEYWIDTH = 100;
const MONKEYHEIGHT = 129;
const STEP = 10;
const NOSTEP = 0

function App() {
  const [players, setPlayers] = useState([]);
  const bodyBoundries = useRef(null);

  function showKeyCode(e) {
    if (e.key === 'ArrowRight') {
      socket.send(JSON.stringify({ y: `${NOSTEP}`, x: `${STEP}` }))
    }
    if (e.key === 'ArrowLeft') {
      socket.send(JSON.stringify({ y: `${NOSTEP}`, x: `-${STEP}` }))
    }
    if (e.key === 'ArrowUp') {
      socket.send(JSON.stringify({ y: `${STEP}`, x: `${NOSTEP}` }))
    }
    if (e.key === 'ArrowDown') {
      socket.send(JSON.stringify({ y: `-${STEP}`, x: `${NOSTEP}` }))
    }
  }
  useEffect(() => {
    document.body.addEventListener('keydown', showKeyCode);
    return () => {
      document.body.removeEventListener('keydown', showKeyCode);
    };
  }, [])

  useLayoutEffect(() => {
    const getBodyBoundries = document.body.getBoundingClientRect();
    bodyBoundries.current = {
      width: (getBodyBoundries.width - MONKEYWIDTH),
      height: (getBodyBoundries.height - MONKEYHEIGHT)
    }
  }, [])

  useEffect(() => {
    const onMessage = event => {
      const parseData = JSON.parse(event.data);
      const { player } = parseData
      setPlayers(currPlayers => {
        const objIndex = currPlayers.findIndex(obj => obj.Id === player.Id);
        if (objIndex !== -1) {
          const clonePlayers = JSON.parse(JSON.stringify(currPlayers));
          clonePlayers[objIndex].shake = false;
          if(!player.show){
            const filtered = clonePlayers.filter(clonePlayer=> clonePlayer.Id !== player.Id); 
            return filtered;
          } 
          if(player.collision === "border"){
            clonePlayers[objIndex].shake = true;
          } 
          clonePlayers[objIndex].p = player.p;
          return clonePlayers;
        }
        //add active:true to THIS user
        let newPlayer
        if(parseData.self.Id === parseData.player.Id){
          newPlayer ={...player,active:true}
        }else{
          newPlayer = {...player}

        }
        return [...currPlayers, newPlayer];
      });
    };

    socket.addEventListener('message', onMessage);

    return () => socket.removeEventListener('message', onMessage);
  },[players]);
  useEffect(() => {
    socket.onopen = () => {
      socket.send(JSON.stringify(bodyBoundries.current))
    }
    return () => {
      socket.onclose = (e) => {
        console.log("socket close connection", e)
      }
    };
  }, []);

  return (
    <>
      {players.map(player => {
        return (
          <div key={player.Id}>
              <div className={`playerImg ${player.shake === true ? "shake" : ''}`} style={{ bottom: player.p.y + 'px', left: player.p.x + 'px', color: `rgb(${player.color[0]},${player.color[1]},${player.color[2]})` }}>
                <span style={{ position: "absolute", top: "-40px", left: "-25px" }}>{player.exceptionType}</span>
                <img alt="player img" src={player.active? activemonkey : monkey } />
              </div>
          </div>
        )
      })}
    </>
  );
}

export default App;


