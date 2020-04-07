
import React, { useCallback } from 'react';
import './App.css';
import monkey from './monkey.png'
import activemonkey from './activemonkey.png'
import { useWebSocket } from './hooks/useWS'
import { useBodyBounderies } from './hooks/useBodyBounderies'
import  {ws, playerMoves} from './Utills.js/Utills'
import { useEventListener } from './hooks/useEventListener'

function App() {
  const bodyBounderies = useBodyBounderies()
  const [players, sendMSG] = useWebSocket(ws, bodyBounderies)

  const showKeyCode = useCallback(
    ({ key }) => {
      sendMSG(playerMoves[key])
    },
    [sendMSG]
  );
  useEventListener('keydown', showKeyCode);
 

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
    </>
  );
}
export default App;


