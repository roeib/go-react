
import React, { useEffect, useState, useRef, useLayoutEffect } from 'react';
import './App.css';
import monkey from './monkey.png'
const socket = new WebSocket('ws://localhost:8080/ws');
const MONKEYWIDTH = 100;
const MONKEYHEIGHT = 129;
const STEP =10;
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
      document.body.remoceEventListener('keydown', showKeyCode);
    };
  }, [])

  useLayoutEffect(() => {
    const getBodyBoundries = document.body.getBoundingClientRect();
    bodyBoundries.current = {
      x: (getBodyBoundries.width -MONKEYWIDTH),
      y: (getBodyBoundries.height - MONKEYHEIGHT )
    }
  }, [])

  useEffect(() => {
    socket.onmessage = (event) => {
      const parseData = JSON.parse(event.data)
      const objIndex = players.findIndex((obj => obj.Id == parseData.Id));
      let clonePlayers = []
      if (objIndex !== -1) {
        clonePlayers = JSON.parse(JSON.stringify(players));
        clonePlayers[objIndex].p = parseData.p
      } else {
        clonePlayers = [...players, parseData]
      }
      setPlayers(clonePlayers)

    }
  }, [players]);



  useEffect(() => {
    socket.onopen = () => {
      console.log("connected successfuly", bodyBoundries.current)
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
          <div  className="plaeyrImg" key={player.Id} style={{ position: 'absolute', right: 0, bottom: player.p.y + 'px', left: player.p.x + 'px', color: `rgb(${player.color[0]},${player.color[1]},${player.color[2]})` }}>
            <span style={{ position: "absolute", top: "-40px", left: "-25px" }}>{player.exceptionType}</span>
            <img src={monkey} />
          </div>
        )
      })}
    </>
  );
}

export default App;


