import React, { useEffect, useState, useRef } from 'react';
import './App.css';




const socket = new WebSocket('ws://localhost:8080/ws');




function App() {

  const [count, setCount] = useState(0);
  const [players, setPlayers] = useState([]);

  function showKeyCode(e) {
    //right
    if (e.keyCode === 39) {
      socket.send(JSON.stringify(10))
    }


    //left
    if (e.keyCode === 37) {
      socket.send(JSON.stringify(-10))
    }


  }
  useEffect(() => {
    document.body.addEventListener('keyup', showKeyCode);

  }, [])







  useEffect(() => {

    socket.onmessage = (event) => {
      // const parseData = JSON.parse(event.data)
      // setCount(x=>x+parseData.X)
      const data = {
        "bananas": [],
        "Players": [
          {
            id: 'Roei',
            x: 8,
            y: 9,
            score: 9
          },
          {
            id: 'Mey',
            x: 8,
            y: 9,
            score: 9
          }
        ]
      }
      const aaa=[data.Players[0].id,data.Players[1].id]
      setPlayers(aaa)
 }
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
      <div>
        players:
      {players.map(player => <span>{player} </span>)}
      </div>
      <div id="red" style={{ left: count + 'px' }}></div>
      <div id="green"></div>
    </>

  );
}

export default App;
