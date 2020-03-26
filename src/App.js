
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
    console.log('gdfgdfg')
    socket.onmessage = (event) => {
      console.log('players',players)
      const parseData = JSON.parse(event.data)
      let isPlayerExist = players.find(player => player.Id === parseData.Id )
      if(isPlayerExist){
          const cc = JSON.parse(JSON.stringify(players));
          const objIndex = cc.findIndex((obj => obj.Id == isPlayerExist.Id));
          cc[objIndex].p =parseData.p
          setPlayers(cc)
      } else {
        console.log('new')
        const newArr = [...players,parseData]
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
        {players.map(player=>{
          return(
<>player {player.p.x}
<img className="monkey"  style={{ left: player.p.x + 'px' }} src={monkey}  />

</>
          )
        })}

       
    </>
  );
}

export default App;


