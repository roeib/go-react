
import React, { useEffect, useState, useRef,useLayoutEffect } from 'react';
import './App.css';
import monkey from './monkey.png'
const socket = new WebSocket('ws://localhost:8080/ws');


function App() {

  const [players, setPlayers] = useState([]);
  const bodyBoundries = useRef(null);

  function showKeyCode(e) {
    //right
    if (e.keyCode === 39) {
      socket.send(JSON.stringify({ y: "0", x: "10" }))
    }
    //left
    if (e.keyCode === 37) {
      socket.send(JSON.stringify({ y: "0", x: "-10" }))
    }
    //up
    if (e.keyCode === 38) {
      socket.send(JSON.stringify({ y: "10", x: "0" }))
    }
    //down
    if (e.keyCode === 40) {
      socket.send(JSON.stringify({ y: "-10", x: "0" }))
    }

  }
  useEffect(() => {
    document.body.addEventListener('keyup', showKeyCode);
    return () => {
      document.body.remoceEventListener('keyup', showKeyCode);
    };
  }, [])

  useLayoutEffect(() => {
    const getBodyBoundries =document.body.getBoundingClientRect() 
    bodyBoundries.current = {
      x:getBodyBoundries.width,
      y:getBodyBoundries.height
    }
  }, [])

  useEffect(() => {
    socket.onmessage = (event) => {
      console.log('players', players)
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
      console.log("connected successfuly",bodyBoundries.current)
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
          <div key={player.Id} style={{ position: 'absolute', right: 0, bottom: player.p.y + 'px', left: player.p.x + 'px', color: `rgb(${player.color[0]},${player.color[1]},${player.color[2]})` }}
          >
            <span style={{ position: "absolute", top: "-40px", left: "-25px" }}>{player.exceptionType}</span>
            <img src={monkey} />
          </div>
        )
      })}


    </>
  );
}

export default App;


