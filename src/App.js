
import React, { useEffect, useState, useRef } from 'react';
import './App.css';
import monkey from './monkey.png'
const socket = new WebSocket('ws://localhost:8080/ws');


function App() {

  const [players, setPlayers] = useState([]);

  function showKeyCode(e) {
    //right
    if (e.keyCode === 39) {
      socket.send(JSON.stringify({ y: "0", x: "10" }))
    }
    //left
    if (e.keyCode === 37) {
      socket.send(JSON.stringify({ y: "0", x: "-10" }))
    }


  }
  useEffect(() => {
    document.body.addEventListener('keyup', showKeyCode);
  }, [])

  useEffect(() => {
    socket.onmessage = (event) => {
      console.log('players', players)
      const parseData = JSON.parse(event.data)
      let isPlayerExist = players.find(player => player.Id === parseData.Id)
      if (isPlayerExist) {
        const cc = JSON.parse(JSON.stringify(players));
        const objIndex = cc.findIndex((obj => obj.Id == isPlayerExist.Id));
        cc[objIndex].p = parseData.p
        setPlayers(cc)
      } else {
        const newArr = [...players, parseData]
        setPlayers(newArr)
      }
    }
  }, [players]);



  useEffect(() => {

    socket.onopen = () => {
      console.log("connected successfuly")
    }

    socket.onclose = (e) => {
      console.log("socket close connection", e)
    }


    socket.onerror = (e) => {
      console.log("socket error", e)
    }
  }, []);


  return (
    <>
      {players.map(player => {
        return (
          <div style={{ position: 'absolute', right: 0, bottom: 0, left: player.p.x + 'px', color: `rgb(${player.color[0]},${player.color[1]},${player.color[2]})` }}
          >
            <span style={{position:"absolute",top:"-40px",left:"25px"}}>{player.exceptionType}</span>
            <img src={monkey} />
          </div>
        )
      })}


    </>
  );
}

export default App;


